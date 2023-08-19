package main

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// Create a variable to store the command
var rootCmd = &cobra.Command{
	Use:   "itldims",
	Short: "A CLI tool to interact with the etcd API",
	Run: func(cmd *cobra.Command, args []string) {
		// By default, show the help message if no subcommand is provided.
		cmd.Help()
	},
}

// Create a function to handle the "get" subcommand
func getServerData(cmd *cobra.Command, args []string) {
	if len(args) != 3 {
		fmt.Println("Usage: itldims get [serverType] [serverIP] [key]")
		return
	}

	serverType := args[0]
	serverIP := args[1]
	key := args[2]

	// Connect to etcd
	ctx := context.TODO()
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints: []string{etcdHost},
	})
	if err != nil {
		log.Printf("Failed to connect to etcd: %v", err)
		return
	}
	defer etcdClient.Close()

	// Construct the etcd key for the server data
	etcdKeyData := fmt.Sprintf("/servers/%s/%s/%s", serverType, serverIP, key)

	response, err := etcdClient.Get(ctx, etcdKeyData, clientv3.WithSort(clientv3.SortByCreateRevision, clientv3.SortAscend))
	if err != nil {
		log.Printf("Failed to retrieve data from etcd: %v", err)
		return
	}

	if len(response.Kvs) > 0 {
		fmt.Printf("Value: %s\n", response.Kvs[0].Value)
	} else {
		fmt.Println("Value not found")
	}
}

// Add a "get" subcommand to the root command
func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "get [serverType] [serverIP] [key]",
		Short: "Retrieve data from the etcd API",
		Args:  cobra.ExactArgs(3),
		Run:   getServerData,
	})
}

// Function to execute the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
