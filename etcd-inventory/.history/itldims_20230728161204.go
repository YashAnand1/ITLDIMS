package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "itldims",
		Short: "Interact with the etcd API",
		Long:  "A command-line tool to interact with the etcd API and tell if the connection has been made",
	}

	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Displays values of an attribute from a server IP",
		Long:  "Find the value of a specific attribute from a Server IP",
		Run: func(cmd *cobra.Command, args []string) {
			aip, _ := cmd.Flags().GetString("aip")
			attribute := args[0]

			serverType := "VM"
			etcdKey := fmt.Sprintf("/servers/%s/%s/%s", serverType, aip, attribute)

			response, err := http.Get("http://localhost:8181" + etcdKey)
			if err != nil {
				log.Fatalf("Failed to connect to the etcd API.")
			}
			defer response.Body.Close()

			if response.StatusCode == http.StatusOK {
				body, err := ioutil.ReadAll(response.Body)
				if err != nil {
					log.Fatalf("Failed to read response body: %v", err)
				}
				fmt.Printf("Attribute value for server IP %s and attribute %s: %s\n", aip, attribute, string(body))
			} else {
				fmt.Printf("Failed to fetch the attribute value for server IP %s and attribute %s.\n", aip, attribute)
			}
		},
	}

	getCmd.Flags().String("aip", "", "Server IP to fetch the attribute value from")
	getCmd.MarkFlagRequired("aip")

	rootCmd.AddCommand(getCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
