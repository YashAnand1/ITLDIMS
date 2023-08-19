package main

import (
	"encoding/json"
	"fmt"
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
				log.Fatalf("Failed to connect to the API.")
			}
			defer response.Body.Close()

			if response.StatusCode == http.StatusOK {
				fmt.Println("Interaction with etcd can be done.")
			} else {
				fmt.Println("Failed to interact with the API.")
			}
		},
	}

	getCmd = &cobra.Command{
		// ... (existing code)
	
		Run: func(cmd *cobra.Command, args []string) {
			// ... (existing code)
	
			response, err := http.Get(etcdURL)
			if err != nil {
				log.Fatalf("Failed to connect to the API: %v", err)
			}
			defer response.Body.Close()
	
			if response.StatusCode == http.StatusOK {
				var data map[string]interface{}
				err = json.NewDecoder(response.Body).Decode(&data)
				if err != nil {
					log.Fatalf("Failed to decode response body: %v", err)
				}
	
				// ... (existing code)
			} else {
				body, _ := ioutil.ReadAll(response.Body)
				fmt.Printf("Received status code: %d\n", response.StatusCode)
				fmt.Printf("Response body: %s\n", string(body))
				fmt.Printf("Failed to fetch the attribute value for server IP %s and attribute %s.\n", aip, attribute)
			}
		},
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
