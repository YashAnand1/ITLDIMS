package main

import (
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
		Long:  "Find the value of a specific attribute from a Server IP",

		Run: func(cmd *cobra.Command, args []string) {
			aip, _ := cmd.Flags().GetString("aip")
			if aip == "" || len(args) == 0 {
				log.Fatal("Please provide a server IP and attribute.")
			}

			attribute := args[0]

			// Send GET request to retrieve the attribute value
			response, err := http.Get(fmt.Sprintf("http://localhost:8181/servers/?prefix=%s", aip))
			if err != nil {
				log.Fatalf("Failed to connect to the etcd API.")
			}
			defer response.Body.Close()

			if response.StatusCode == http.StatusOK {
				body, err := ioutil.ReadAll(response.Body)
				if err != nil {
					log.Fatalf("Failed to read response body: %v", err)
				}

				// Parse the response and find the attribute value
				attributeValue := findAttributeValue(string(body), aip, attribute)
				if attributeValue != "" {
					fmt.Printf("Attribute value for server IP %s and attribute %s: %s\n", aip, attribute, attributeValue)
				} else {
					fmt.Printf("Attribute value not found for server IP %s and attribute %s.\n", aip, attribute)
				}
			} else {
				fmt.Printf("Failed to fetch the attribute value for server IP %s and attribute %s.\n", aip, attribute)
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

func findAttributeValue(response, aip, attribute string) string {
	lines := strings.Split(response, "\n")
	for _, line := range lines {
		if strings.Contains(line, aip+"/"+attribute) {
			parts := strings.Split(line, "/")
			if len(parts) > 0 {
				return parts[len(parts)-1]
			}
		}
	}
	return ""
}
