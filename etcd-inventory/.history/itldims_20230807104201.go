package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "itldims",
		Short: "Interact with the etcd API",
		Long:  "A command-line tool to interact with the etcd API and tell if the connection has been made",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Successfully connected with API. Interaction with etcd can be done.")
		},
	}

	getCmd = &cobra.Command{
		Use:   "get <key>",
		Short: "Gets the value of a key from the etcd API",
		Run:   getCommandFunc,
	}
)

func main() {
	rootCmd.AddCommand(getCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func getCommandFunc(cmd *cobra.Command, args []string) {
	// Check if there are enough arguments
	if len(args) < 1 {
		log.Fatal("etcdctl get command needs one argument as key")
	}

	// Construct the API URL
	apiURL := fmt.Sprintf("http://localhost:8181%s", args[0])

	// Check if we can connect to the etcd API
	response, err := http.Get(apiURL)
	if err != nil {
		log.Fatalf("Failed to connect to the etcd API: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Fatalf("Failed to retrieve data from etcd API. Status code: %d", response.StatusCode)
	}

	// Print the retrieved key-value
	value := response.Header.Get("Value")
	if value == "" {
		log.Fatalf("Value not found for the specified key")
	}

	fmt.Printf("Key: %s, Value: %s\n", args[0], value)
}
