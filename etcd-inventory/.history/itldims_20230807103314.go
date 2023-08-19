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
	itldims = &cobra.Command{
		Use:   "itldims",
		Short: "Interact with the etcd API",
		Long:  "A command-line tool to interact with the etcd API and tell if the connection has been made",
		Run:   itldimsCommandFunc,
	}
)

func main() {
	if err := itldims.Execute(); err != nil {
		log.Fatal(err)
	}
}

func itldimsCommandFunc(cmd *cobra.Command, args []string) {
	// Check if we can connect to the etcd API
	response, err := http.Get("http://localhost:8181/servers/")
	if err != nil {
		log.Fatalf("Failed to connect to the etcd API.")
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		fmt.Println("Successfully connected with API. Interaction with etcd can be done.")
	}

	// Connect to etcd and run "etcdctl get" command
	if len(args) == 0 {
		log.Fatal("etcdctl get command needs at least one argument as key")
	}

	// Initialize the etcd client
	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"http://localhost:2379"},
	})
	if err != nil {
		log.Fatalf("Failed to initialize etcd client: %v", err)
	}
	defer client.Close()

	// Call the "etcdctl get" function with the provided key
	key := args[0]
	opts := getGetOp(args[1:])
	resp, err := client.Get(context.Background(), key, opts...)
	if err != nil {
		log.Fatalf("Failed to retrieve data from etcd: %v", err)
	}

	// Print the retrieved key-values
	for _, kv := range resp.Kvs {
		fmt.Printf("Key: %s, Value: %s\n", string(kv.Key), string(kv.Value))
	}
}

func getGetOp(args []string) []clientv3.OpOption {
	// Process the command arguments and extract options for the "etcdctl get" command
	// (Options have already been removed in the original "etcdctl get" code)
	var opts []clientv3.OpOption
	if len(args) > 0 {
		opts = append(opts, clientv3.WithRange(args[0]))
	}
	return opts
}
