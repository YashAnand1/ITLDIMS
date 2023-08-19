package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

var serverIP string
var attribute string

var rootITLDIMS = &cobra.Command{
	Use:   "itldims",
	Short: "Interact with the etcd API",
	Long:  "A command-line tool to interact with the etcd API and tell if the connection has been made",
	Run: func(cmd *cobra.Command, args []string) {
		url := fmt.Sprintf("http://%s/servers/%s", serverIP, attribute)
		response, err := http.Get(url)
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

func init() {
	rootITLDIMS.PersistentFlags().StringVarP(&serverIP, "aip", "", "localhost:8181", "Server IP address and port")
	rootITLDIMS.PersistentFlags().StringVarP(&attribute, "attribute", "", "", "Attribute to interact with")

	if err := rootITLDIMS.MarkPersistentFlagRequired("attribute"); err != nil {
		log.Fatal(err)
	}
}

func main() {
	if err := rootITLDIMS.Execute(); err != nil {
		log.Fatal(err)
	}
}
