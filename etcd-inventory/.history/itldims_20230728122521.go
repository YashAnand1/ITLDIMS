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
				log.Fatalf("Failed to connect to the API.")
			}
			defer response.Body.Close()

			if response.StatusCode == http.StatusOK {
				fmt.Println("Interaction with etcd can be done.")
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

			etcdURL := fmt.Sprintf("http://localhost:8181/servers/?prefix=true&keysOnly=true&separator=%s", "/")

			response, err := http.Get(etcdURL)
			if err != nil {
				log.Fatalf("Failed to connect to the API.")
			}
			defer response.Body.Close()

			if response.StatusCode == http.StatusOK {
				body, err := ioutil.ReadAll(response.Body)
				if err != nil {
					log.Fatalf("Failed to read response body: %v", err)
				}

				// Split the response into individual lines (each line represents a key path)
				lines := strings.Split(string(body), "\n")

				// Filter the lines to find the matching server IP and attribute
				var foundKeys []string
				for _, line := range lines {
					if strings.Contains(line, aip) && strings.Contains(line, attribute) {
						foundKeys = append(foundKeys, line)
					}
				}

				if len(foundKeys) == 0 {
					fmt.Printf("No values found for server IP %s and attribute %s.\n", aip, attribute)
					return
				}

				// Fetch the values for each found key
				for _, key := range foundKeys {
					etcdKeyURL := fmt.Sprintf("http://localhost:8181%s", key)
					valueResponse, err := http.Get(etcdKeyURL)
					if err != nil {
						log.Fatalf("Failed to fetch attribute value for server IP %s and attribute %s.", aip, attribute)
					}
					defer valueResponse.Body.Close()

					if valueResponse.StatusCode == http.StatusOK {
						valueBody, err := ioutil.ReadAll(valueResponse.Body)
						if err != nil {
							log.Fatalf("Failed to read response body: %v", err)
						}

						// Extract the value from the response body
						value := string(valueBody)
						fmt.Printf("Attribute value for server IP %s and attribute %s: %s\n", aip, attribute, value)
					} else {
						fmt.Printf("Failed to fetch attribute value for server IP %s and attribute %s.\n", aip, attribute)
					}
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
