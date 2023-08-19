package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

var (
	rootITLDIMS = &cobra.Command{
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

	getCmd = &cobra.Command{
		Use:   "get",
		Short: "Displays values of an attribute from a server IP",
		Long:  "Find the value of a specific attribute from a Server IP",

		Run: func(cmd *cobra.Command, args []string) {
			server, _ := cmd.Flags().GetString("server")
			if server == "" || len(args) == 0 {
				log.Fatal("Enter correct server IP and attribute.")
			}

			attribute := args[0]

			// Call the function to search through all keys and find the matching ones
			err := searchKeys(server, attribute)
			if err != nil {
				log.Fatalf("Failed to search for keys in etcd: %v", err)
			}
		},
	}
)

func init() {
	rootITLDIMS.AddCommand(getCmd)
	getCmd.Flags().String("server", "", "Server IP to fetch the attribute value from")
	getCmd.MarkFlagRequired("server")
}

func searchKeys(serverIP, attribute string) error {
	// Connect to etcd
	response, err := http.Get("http://localhost:8181/servers/")
	if err != nil {
		log.Fatalf("Failed to connect to the etcd API.")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get server data from etcd")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	// Split the response body by newline to get individual lines
	lines := strings.Split(string(body), "\n")
	for _, line := range lines {
		// Extract server type, server IP, and attributes from each line
		parts := strings.Split(line, "/")
		if len(parts) >= 3 {
			srvIP := parts[2]
			attributes := parts[3:]

			// Check if the server IP and attribute match the ones provided in the command
			if srvIP == serverIP && contains(attributes, attribute) {
				// Fetch the key's value from etcd and display it
				etcdKey := strings.Join(parts[1:], "/")
				value, err := getKeyValue(etcdKey)
				if err != nil {
					return fmt.Errorf("failed to get value for key %s: %v", etcdKey, err)
				}
				fmt.Printf("Key: %s, Value: %s\n", etcdKey, value)
			}
		}
	}

	// If no match found, return an error
	return fmt.Errorf("no matching key found in etcd for server IP %s and attribute %s", serverIP, attribute)
}

// Helper function to check if a slice contains a specific element
func contains(slice []string, element string) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}

func getKeyValue(key string) (string, error) {
	// Connect to etcd to fetch the value for the given key
	response, err := http.Get("http://localhost:8181" + key)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get value for key %s from etcd", key)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body for key %s: %v", key, err)
	}

	return string(body), nil
}

func main() {
	if err := rootITLDIMS.Execute(); err != nil {
		log.Fatal(err)
	}
}
