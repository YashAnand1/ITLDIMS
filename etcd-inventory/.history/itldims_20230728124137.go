package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var rootCmd = &cobra.Command{
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
			fmt.Println("interaction with etcd can be done.")
		} else {
			fmt.Println("Failed to interact with the API.")
		}
	},
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the attribute value for a specific server IP",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		aip := args[0]
		attribute := args[1]

		// Connect to etcd
		ctx := context.TODO()
		etcdClient, err := clientv3.New(clientv3.Config{
			Endpoints: []string{etcdHost},
		})
		if err != nil {
			log.Fatalf("Failed to connect to etcd: %v", err)
		}
		defer etcdClient.Close()

		// Construct the etcd key for the server IP and attribute
		etcdKeyData := fmt.Sprintf("/servers/*/%s/%s", aip, attribute)

		response, err := etcdClient.Get(ctx, etcdKeyData, clientv3.WithSort(clientv3.SortByCreateRevision, clientv3.SortAscend))
		if err != nil {
			log.Fatalf("Failed to fetch the attribute value for server IP %s and attribute %s: %v", aip, attribute, err)
		}

		if len(response.Kvs) > 0 {
			value := string(response.Kvs[0].Value)
			fmt.Println("Value:", value)
		} else {
			fmt.Printf("No values found for server IP %s and attribute %s.\n", aip, attribute)
		}
	},
}

func main() {
	rootCmd.AddCommand(getCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
