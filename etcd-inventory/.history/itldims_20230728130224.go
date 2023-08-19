package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

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
		Long:  "Find the value of a specific attribute from a Server IP",

		Run: func(cmd *cobra.Command, args []string) {
			aip, _ := cmd.Flags().GetString("aip")
			if aip == "" {
				log.Fatal("Please provide a server IP.")
			}
			if len(args) == 0 {
				log.Fatal("Please provide an attribute.")
			}
			attribute := args[0]

			// Connect to etcd
			ctx := context.TODO()
			etcdClient, err := clientv3.New(clientv3.Config{
				Endpoints: []string{"localhost:2379"},
			})
			if err != nil {
				log.Printf("Failed to connect to etcd: %v", err)
				return
			}
			defer etcdClient.Close()

			// Construct the etcd key for the attribute value
			etcdKey := fmt.Sprintf("/servers/%s/%s/%s", "servertype", aip, attribute)

			response, err := etcdClient.Get(ctx, etcdKey)
			if err != nil {
				log.Fatalf("Failed to retrieve attribute value from etcd: %v", err)
			}

			if len(response.Kvs) > 0 {
				value := string(response.Kvs[0].Value)
				fmt.Printf("Attribute value for server IP %s and attribute %s: %s\n", aip, attribute, value)
			} else {
				fmt.Printf("Attribute value not found for server IP %s and attribute %s.\n", aip, attribute)
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
