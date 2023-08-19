package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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
			} else {
				fmt.Println("Failed to interact with the API.")
			}
		},
	}

	getCmd = &cobra.Command{
		Use:   "get",
		Short: "Displays values of an attribute from a server IP",
		Long:  "Find the value of a specific attribute from a Server IP",

		Run: func(cmd *cobra.Command, args []string) {
			aip, _ := cmd.Flags().GetString("aip")
			if aip == "" || len(args) != 1 {
				log.Fatal("Please provide a server IP and attribute.")
			}

			attribute := args[0]

			// First, fetch the server type based on the provided server IP
			serverType, err := getServerType(aip)
			if err != nil {
				log.Fatalf("Failed to fetch the server type: %v", err)
			}

			etcdKey := fmt.Sprintf("/servers/%s/%s/%s", serverType, aip, attribute)
			apiURL := "http://localhost:8181" + etcdKey

			response, err := http.Get(apiURL)
			if err != nil {
				log.Fatalf("Failed to connect to the etcd API: %v", err)
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
)

type serverInfo struct {
	ServerType string `json:"server_type"`
}

func init() {
	rootITLDIMS.AddCommand(getCmd)

	getCmd.Flags().String("aip", "", "Server IP to fetch the attribute value from")
	getCmd.MarkFlagRequired("aip")
}

func main() {
	if err := rootITLDIMS.Execute(); err != nil {
		log.Fatal(err)
	}
}

// Function to fetch the server type based on the provided server IP
func getServerType(aip string) (string, error) {
	response, err := http.Get("http://localhost:8181/servers/" + aip)
	if err != nil {
		return "", fmt.Errorf("failed to connect to the etcd API: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read response body: %v", err)
		}
		fmt.Println("API Response Body:", string(body))
		var info serverInfo
		err = json.Unmarshal(body, &info)
		if err != nil {
			return "", fmt.Errorf("failed to decode server info: %v", err)
		}
		return info.ServerType, nil
	}

	return "", fmt.Errorf("failed to fetch server type for IP: %s", aip)
}
