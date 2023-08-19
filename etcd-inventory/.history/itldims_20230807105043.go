package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// NewGetCommand returns the cobra command for "get".
func NewGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <key> [range_end]",
		Short: "Gets the key or a range of keys",
		Run:   getCommandFunc,
	}
	return cmd
}

// getCommandFunc executes the "get" command.
func getCommandFunc(cmd *cobra.Command, args []string) {
	key, opts := getGetOp(args)
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

func getGetOp(args []string) (string, []clientv3.OpOption) {
	if len(args) == 0 {
		fmt.Println("get command needs one argument as key and an optional argument as range_end")
		os.Exit(1)
	}

	var opts []clientv3.OpOption
	key := args[0]
	if len(args) > 1 {
		opts = append(opts, clientv3.WithRange(args[1]))
	}

	return key, opts
}

func main() {
	cmd := NewGetCommand()
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
