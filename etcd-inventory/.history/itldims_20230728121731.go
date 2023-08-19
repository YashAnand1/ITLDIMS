package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

var (
	rootITLDIMS = &cobra.Command{
		Use:   "itldims",
		Short: "Interact with the etcd API",
		Long:  "A command-line tool to interact with the etcd API and tell if the connection has been made",
		Run: func(cmd *cobra.Command, args []string) {
			response, err := http.Get("http://localhost:8181/servers/")
			if err != nil {
				log.Fatalf("Failed to connect to the etcd API.")
			}
			defer response.Body.Close()

			if response.StatusCode == http.StatusOK {
				fmt.Println("Successfully connected with API. Interaction with etcd can be done.")
			} else {
				fmt.Println("Failed to interact with the API.")
			}
		},
	}

	getCmd = &cobra.Command{
		Use:   "get",
		Short: "Displays values of an attribute from a server IP",
		Long:  "Find the value of a specific attribute from a Server IP alone without having to enter the key",

		Run: func(cmd *cobra.Command, args []string) {
			aip, _ := cmd.Flags().GetString("aip")
			if aip == "" || len(args) == 0 {
				log.Fatal("Please provide a server IP and attribute.")
			}

			attribute := args[0]

			// Fetch all keys from the etcd database with the provided server IP as the prefix
			// and then search for the key that contains the provided attribute
			response, err := http.Get(fmt.Sprintf("http://localhost:8181/servers/%s", aip))
			if err != nil {
				log.Fatalf("Failed to connect to the etcd API.")
			}
			defer response.Body.Close()

			if response.StatusCode != http.StatusOK {
				fmt.Printf("Failed to fetch data for server IP %s.\n", aip)
				return
			}

			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Fatalf("Failed to read response body: %v", err)
			}

			// Parse the JSON response containing the key-value pairs
			var data map[string]interface{}
			err = json.Unmarshal(body, &data)
			if err != nil {
				log.Fatalf("Failed to parse JSON response: %v", err)
			}

			// Search for the key that contains the desired attribute
			var value interface{}
			for key, v := range data {
				if strings.Contains(key, attribute) {
					value = v
					break
				}
			}

			if value == nil {
				fmt.Printf("No values found for server IP %s and attribute %s.\n", aip, attribute)
			} else {
				fmt.Printf("Attribute value for server IP %s and attribute %s: %v\n", aip, attribute, value)
			}
		},
	}
)

func init() {
	rootITLDIMS.AddCommand(getCmd)

	getCmd.Flags().String("aip", "", "Server IP to fetch the attribute value from")
	getCmd.MarkFlagRequired("aip")
}

func main() {
	if err := rootITLDIMS.Execute(); err != nil {
		log.Fatal(err)
	}
}
