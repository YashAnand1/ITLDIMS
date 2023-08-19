package main

import (
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
			}
		},
	}

	getServersCmd = &cobra.Command{
		Use:   "get --servers <input1> <input2>",
		Short: "Get keys with specific inputs from etcd API",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			input1 := args[0]
			input2 := args[1]

			response, err := http.Get("http://localhost:8181/servers/")
			if err != nil {
				log.Fatalf("Failed to connect to the etcd API.")
			}
			defer response.Body.Close()

			if response.StatusCode == http.StatusOK {
				fmt.Printf("Server data with input1='%s' and input2='%s':\n", input1, input2)
				fmt.Println("---------------------")
				body, err := ioutil.ReadAll(response.Body)
				if err != nil {
					log.Fatalf("Failed to read response body: %v", err)
				}
				fmt.Println(string(body))
			}
		},
	}
)

func init() {
	rootITLDIMS.AddCommand(getServersCmd)
}

func main() {
	if err := rootITLDIMS.Execute(); err != nil {
		log.Fatal(err)
	}
}
