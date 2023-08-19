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
		Short: "Get the value of a specific attribute from a server IP",
		Run: func(cmd *cobra.Command, args []string) {
			serverIP, _ := cmd.Flags().GetString("aip")
			attribute, _ := cmd.Flags().GetString("attribute")

			fmt.Printf("Fetching attribute %s from server IP %s\n", attribute, serverIP)
			fmt.Println("Attribute value: DummyValue") // Replace this with your actual etcd interaction logic
		},
	}
)

func init() {
	rootCmd.Flags().StringP("aip", "a", "", "Server IP for attribute retrieval")
	rootCmd.Flags().StringP("attribute", "t", "", "Attribute name to retrieve")
}

func main() {
	rootCmd.AddCommand(getCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
