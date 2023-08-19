package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra" // Import the Cobra library for CLI functionality
)

// FOR itldims get servers - LIST ALL SERVERS
var servers = &cobra.Command{
	Use:   "servers",
	Short: "all servers displayed",
	Long:  "all servers displayed in a listed form",
	Run: func(cmd *cobra.Command, args []string) {
		data, err := fetchDataFromAPI() // Fetch data from the API
		if err != nil {
			fmt.Printf("Unable to connect due to %s", err)
		}

		// Create a map to store server IPs UNIQUELY
		IPs := make(map[string]string)

		for key := range parseKeyValuePairs(data) {
			splitKey := strings.Split(key, "/")
			serverIP := splitKey[3] // Extract the server IP from the split key
			IPs[serverIP] = serverIP
		}

		for serverIP := range IPs { // Go through stored server IPs
			fmt.Printf("%s\n----------------------------\n", serverIP)
		}
	},
}

// FOR itldims get types - LIST ALL SERVERTYPES
var Types = &cobra.Command{
	Use:   "types",
	Short: "For finding running server types",
	Long:  "For finding the running server types from all the servers",
	Run: func(cmd *cobra.Command, args []string) {
		data, err := fetchDataFromAPI()
		if err != nil {
			fmt.Printf("Could not connect to API %v", err)
		}

		Types := make(map[string]string) // for uniquely storing server types

		// Loop through keys obtained from parsing data
		for key := range parseKeyValuePairs(data) {
			splitKey := strings.Split(key, "/")
			serverType := splitKey[2]
			Types[serverType] = serverType
		}

		for serverType := range Types {
			fmt.Printf("%s\n----------------------------\n", serverType)
		}
	},
}

// FOR itldims get attributes - LISTS ATTRIBUTES
var attributes = &cobra.Command{
	Use:   "attributes",
	Short: "Displays all attributes",
	Long:  "Displays all attributes that are running",
	Run: func(cmd *cobra.Command, args []string) {
		data, err := fetchDataFromAPI()
		if err != nil {
			fmt.Printf("Could not connect to the API. Error: %v", err)
		}

		ATs := make(map[string]string) // map to uniquely store attributes

		for key := range parseKeyValuePairs(data) {
			splitKey := strings.Split(key, "/")
			attribute := splitKey[4]
			ATs[attribute] = attribute // Store the attribute in the map
		}

		for attribute := range ATs {
			fmt.Printf("%s\n----------------------------\n", attribute)
		}

	},
}
