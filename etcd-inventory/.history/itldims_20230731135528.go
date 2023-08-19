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

type ServerInfo struct {
	ServerType string            `json:"server-type"`
	Attributes map[string]string `json:"attributes"`
}

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
			}
		},
	}

	getCmd = &cobra.Command{
		Use:   "get",
		Short: "Displays values of an attribute from a server IP",
		Long:  "Find the value of a specific attribute from a Server IP",

		Run: func(cmd *cobra.Command, args []string) {
			server, _ := cmd.Flags().GetString("server")
			if server == "" || len(args) == 0 {
				log.Fatal("Enter correct server IP and attribute.")
			}

			attribute := args[0]

			// Send a request to get all keys with the given server IP
			response, err := http.Get(fmt.Sprintf("http://localhost:8181/servers/%s", server))
			if err != nil {
				log.Fatalf("Failed to connect to the etcd API.")
			}
			defer response.Body.Close()

			if response.StatusCode == http.StatusOK {
				body, err := ioutil.ReadAll(response.Body)
				if err != nil {
					log.Fatalf("Failed to read response body: %v", err)
				}

				var serverInfo ServerInfo
				if err := json.Unmarshal(body, &serverInfo); err != nil {
					log.Fatalf("Failed to unmarshal JSON response: %v", err)
				}

				if serverInfo.Attributes == nil {
					fmt.Printf("Server IP %s not found in the etcd database.\n", server)
					return
				}

				// Iterate over the attributes to find a key with a matching attribute
				for key, value := range serverInfo.Attributes {
					// Check if the key contains the given server IP and attribute
					if strings.Contains(key, server) && strings.Contains(key, attribute) {
						fmt.Printf("Attribute value for server IP %s and attribute %s: %s\n", server, attribute, value)
						return
					}
				}

				fmt.Printf("Attribute %s not found for server IP %s.\n", attribute, server)
			} else {
				fmt.Printf("Server IP %s not found in the etcd database.\n", server)
			}
		},
	}
)

func init() {
	rootITLDIMS.AddCommand(getCmd)
	getCmd.Flags().String("server", "", "Server IP to fetch the attribute value from")
	getCmd.MarkFlagRequired("server")
}

func main() {
	if err := rootITLDIMS.Execute(); err != nil {
		log.Fatal(err)
	}
}

//itldims is currently entering the server type, server ip, attribute are being filled in "localhost:8181/servers/server-type/server-ip/attribute"
//xargs cannot be used at the moment because it is
