package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

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
				fmt.Println("Successfully onnected with API. Interaction with etcd can be done.")
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
			serverType := "VM"
			etcdKey := fmt.Sprintf("/servers/%s/%s/%s", serverType, server, attribute)

			response, err := http.Get("http://localhost:8181" + etcdKey)
			if err != nil {
				log.Fatalf("Failed to connect to the etcd API.")
			}
			defer response.Body.Close()

			if response.StatusCode == http.StatusOK {
				body, err := io.ReadAll(response.Body)
				if err != nil {
					log.Fatalf("Failed to read response body: %v", err)
				}
				fmt.Printf("Attribute value for server IP %s and attribute %s: %s\n", server, attribute, string(body))
			}
		},
	}
)

func init() {
	rootITLDIMS.AddCommand(getCmd)
	getCmd.Flags().String("server", "", "Server IP")
	getCmd.MarkFlagRequired("server")
}

func main() {
	if err := rootITLDIMS.Execute(); err != nil {
		log.Fatal(err)
	}
}
