package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	serverIP string

	itldims = &cobra.Command{
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

	get = &cobra.Command{
		Use:   "get",
		Short: "Gets the key or a range of keys",
		Run:   getCommandFunc,
	}
)

func getCommandFunc(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("get command needs one argument as key")
		os.Exit(1)
	}

	// Construct the key using the serverIP and the provided attribute
	key := fmt.Sprintf("/servers/VM/%s/%s", serverIP, args[0])

	// Remaining code remains the same
	opts := getGetOp(args)

	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"http://localhost:2379"},
	})
	if err != nil {
		fmt.Printf("Failed to initialize etcd client: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	ctx := cmd.Context()
	resp, err := client.Get(ctx, key, opts...)
	if err != nil {
		fmt.Printf("Failed to retrieve data from etcd: %v\n", err)
		os.Exit(1)
	}

	// Print the retrieved key-values
	for _, kv := range resp.Kvs {
		fmt.Printf("Key: %s, Value: %s\n", string(kv.Key), string(kv.Value))
	}
}

func getGetOp(args []string) []clientv3.OpOption {
	if len(args) == 0 {
		fmt.Println("get command needs one argument as key")
		os.Exit(1)
	}

	var opts []clientv3.OpOption
	if len(args) > 1 {
		opts = append(opts, clientv3.WithRange(args[1]))
	}

	return opts
}

func init() {
	get.Flags().StringVar(&serverIP, "server", "", "Server IP")
	get.MarkFlagRequired("server")
	itldims.AddCommand(get)
}

func main() {
	if err := itldims.Execute(); err != nil {
		log.Fatal(err)
	}
}
