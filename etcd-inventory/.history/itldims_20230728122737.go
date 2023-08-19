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

			etcdURL := fmt.Sprintf("http://localhost:8181/servers/?prefix=%s/", aip)

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

				var data map[string]interface{}
				err = json.Unmarshal(body, &data)
				if err != nil {
					log.Fatalf("Failed to unmarshal response: %v", err)
				}

				// Find the attribute data for the given server IP
				for key, value := range data {
					parts := strings.Split(key, "/")
					if len(parts) >= 4 && parts[2] == aip {
						attr := parts[len(parts)-1]
						if attr == attribute {
							fmt.Printf("Attribute value for server IP %s and attribute %s: %v\n", aip, attribute, value)
						}
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
