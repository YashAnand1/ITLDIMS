package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

var rootITLDIMS = &cobra.Command{
	Use:   "itldims",
	Short: "Interact with the ITL dimensions API",
	Long:  "A command-line tool to interact with the ITL dimensions API and display the message 'interaction with etcd can be done.'",
}

var getCommand = &cobra.Command{
	Use:   "get",
	Short: "Get attribute from the ITL dimensions API for a specific server IP",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serverIP := cmd.Flag("serverip").Value.String()
		attribute := args[0]

		// Construct the API URL using the provided server IP and attribute
		apiURL := fmt.Sprintf("http://localhost:8181/servers/%s/%s", serverIP, attribute)

		response, err := http.Get(apiURL)
		if err != nil {
			log.Fatalf("Failed to connect to the API.")
		}
		defer response.Body.Close()

		if response.StatusCode == http.StatusOK {
			fmt.Println("interaction with etcd can be done.")
			// You can parse and display the relevant data from the API response here if needed.
		} else {
			fmt.Println("Failed to interact with the API.")
		}
	},
}

func init() {
	// Add the "get" subcommand to the root command
	rootITLDIMS.AddCommand(getCommand)

	// Add a flag to the "get" command to specify the attribute
	getCommand.Flags().StringP("attribute", "a", "", "Attribute to retrieve")
	getCommand.MarkFlagRequired("attribute")

	// Add a flag to the "get" command to specify the server IP
	getCommand.Flags().StringP("serverip", "s", "", "Server IP")
	getCommand.MarkFlagRequired("serverip")
}

func main() {
	if err := rootITLDIMS.Execute(); err != nil {
		log.Fatal(err)
	}
}
