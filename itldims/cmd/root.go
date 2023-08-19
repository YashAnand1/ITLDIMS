package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

var (
	itldims = &cobra.Command{
		Use:   "itldims",
		Short: "Interact with the etcd API",
		Long:  "A command-line tool to interact with the etcd API and check connection",
		Run: func(cmd *cobra.Command, args []string) { // Extracted function
			response, err := http.Get("http://localhost:8181/servers/")
			if err != nil {
				fmt.Printf("Failed to connect to the etcd API.")
			}
			defer response.Body.Close()

			if response.StatusCode == http.StatusOK {
				fmt.Println("Connected to API. Interaction with etcd can be done.")
			}
		},
	}
)

func init() {
	itldims.AddCommand(get)
}

func Execute() {
	if err := itldims.Execute(); err != nil {
		log.Fatal(err)
	}
}
