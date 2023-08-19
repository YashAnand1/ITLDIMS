package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
	clientv3 "go.etcd.io/etcd/client/v3"
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
			if aip == "" || len(args) != 1 {
				log.Fatal("Please provide a server IP and attribute.")
			}

			attribute := args[0]

			etcdClient, err := clientv3.New(clientv3.Config{
				Endpoints: []string{"http://localhost:2379"},
			})
			if err != nil {
				log.Fatalf("Failed to connect to etcd: %v", err)
			}
			defer etcdClient.Close()

			ctx := context.Background()

			// Fetch all keys that match the provided attribute prefix
			response, err := etcdClient.Get(ctx, fmt.Sprintf("/servers/*/servers/%s/%s", aip, attribute), clientv3.WithPrefix())
			if err != nil {
				log.Fatalf("Failed to fetch attribute values for server IP %s and attribute %s: %v", aip, attribute, err)
			}

			if len(response.Kvs) == 0 {
				fmt.Printf("No values found for server IP %s and attribute %s.\n", aip, attribute)
				return
			}

			for _, kv := range response.Kvs {
				// Extract server type, server IP, and attribute from the key
				keyParts := strings.Split(string(kv.Key), "/")
				serverType := keyParts[2]
				serverIP := keyParts[3]
				attribute := keyParts[4]

				fmt.Printf("Server Type: %s, Server IP: %s, Attribute: %s, Value: %s\n", serverType, serverIP, attribute, string(kv.Value))
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
