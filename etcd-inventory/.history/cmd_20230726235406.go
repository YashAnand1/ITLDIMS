package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "itldims",
	Short: "Interact with the ITL dimensions API",
	Long:  "A command-line tool to interact with the ITL dimensions API and display the message 'interaction with etcd can be done.'",
	Run: func(cmd *cobra.Command, args []string) {
		response, err := http.Get("http://localhost:8181/servers/")
		if err != nil {
			log.Fatalf("Failed to connect to the API.")
		}
		defer response.Body.Close()

		if response.StatusCode == http.StatusOK {
			fmt.Println("interaction with etcd can be done.")
		} else {
			fmt.Println("Failed to interact with the API.")
		}
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
