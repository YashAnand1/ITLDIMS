// This is the main package of the Go program.

// Import required packages.
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/cobra" // Importing the Cobra library for command-line interaction.
)

// Define a root command named "itldims".
var (
	itldims = &cobra.Command{
		Use:   "itldims",
		Short: "Interact with the etcd API", // Short description of the command.
		Long:  "A command-line tool to interact with the etcd API and check connection", // Longer description of the command.
		Run: func(cmd *cobra.Command, args []string) { // Function to execute when the "itldims" command is run.
			response, err := http.Get("http://localhost:8181/servers/") // Send a GET request to the etcd API.
			if err != nil { // If there's an error, log a message and exit.
				log.Fatalf("Failed to connect to the etcd API.")
			}
			defer response.Body.Close() // Close the response body when done.

			if response.StatusCode == http.StatusOK {
				fmt.Println("Connected to API. Interaction with etcd can be done.") // Print a message if the connection is successful.
			}
		},
	}

	keyOnly bool // Flag to determine whether to display only keys without values.
)

// Define a subcommand named "get".
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Search Attributes & Values from etcd API", // Short description of the subcommand.
	Long: `Data retrieval can be done using 'itldims get <input1> <input2>' or 'itldims get <input1>'.
	// Long description of the subcommand, explaining its usage and available options.
	Args: cobra.RangeArgs(1, 2), // Define the range of allowed arguments (1 to 2).

	Run: func(cmd *cobra.Command, args []string) { // Function to execute when the "get" subcommand is run.
		data, err := fetchDataFromAPI() // Fetch data from the etcd API.
		if err != nil {
			log.Fatalf("Failed to fetch data from the etcd API: %v", err) // Log an error if data fetching fails.
		}

		if len(args) == 1 {
			args = append(args, "servers") // If only one argument is provided, default to "servers".
		}

		for key, value := range parseKeyValuePairs(data) { // Loop through the parsed key-value pairs.
			if strings.Contains(key, "{") || strings.Contains(key, "}") ||
				strings.Contains(value, "{") || strings.Contains(value, "}") {
				continue // Skip key-value pairs with problematic characters.
			}

			if !strings.Contains(key, "data") && strings.Contains(key, args[0]) || strings.Contains(value, args[0]) {
				if len(args) > 1 && !strings.Contains(key, args[1]) && !strings.Contains(value, args[1]) {
					continue // Skip if arguments don't match the key or value.
				}

				fmt.Println(key) // Print the key.
				if !keyOnly {
					fmt.Println(value) // Print the value if the "keyOnly" flag is not set.
				}

				fmt.Println() // Print an empty line for separation.
			}
		}
	},
}

// Function to fetch data from the etcd API.
func fetchDataFromAPI() (string, error) {
	response, err := http.Get("http://localhost:8181/servers/") // Send a GET request to the etcd API.
	if err != nil {
		return "", err // Return an error if the request fails.
	}
	defer response.Body.Close() // Close the response body when done.

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch data. Status code: %d", response.StatusCode) // Return an error if the status code is not OK.
	}

	data, err := io.ReadAll(response.Body) // Read the response body.
	if err != nil {
		return "", err // Return an error if reading the body fails.
	}

	return string(data), nil // Return the fetched data as a string.
}

// Function to parse the fetched data into key-value pairs.
func parseKeyValuePairs(data string) map[string]string {
	result := make(map[string]string) // Initialize a map to store key-value pairs.

	keyValuePairs := strings.Split(data, "Key:") // Split data into key-value pairs using "Key:" as the delimiter.

	for _, kv := range keyValuePairs {
		kv = strings.TrimSpace(kv)
		if len(kv) == 0 {
			continue // Skip empty key-value pairs.
		}

		lines := strings.Split(kv, "Value:") // Split each key-value pair into lines using "Value:" as the delimiter.
		if len(lines) == 2 {
			key := strings.TrimSpace(lines[0])   // Extract and trim the key.
			value := strings.TrimSpace(lines[1]) // Extract and trim the value.
			result[key] = value                  // Store the key-value pair in the result map.
		}
	}
	return result // Return the map of parsed key-value pairs.
}

// Initialization function, where flags and subcommands are added.
func init() {
	itldims.AddCommand(getCmd) // Add the "get" subcommand to the "itldims" command.
	getCmd.Flags().BoolVar(&keyOnly, "listall", false, "Display only keys without values") // Define the "listall" flag for the "get" subcommand.
}

// The main function, where the program execution begins.

func main() {
	if err := itldims.Execute(); err != nil {
		log.Fatal(err) // Log an error and exit if command execution fails.
	}
}
