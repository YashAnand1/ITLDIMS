package main

import ( //Importation of packages which help with formatting, logging and interacting with http addresses.
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra" //importation of the cobra tool for creating commands.
)

var rootITLDIMS = &cobra.Command{
	Use:   "itl",
	Short: "Interact with the etcd API",
	Long:  "A command-line tool to interact with the etcd API and tell if the connection has been made",
	Run: func(cmd *cobra.Command, args []string) {
		response, err := http.Get("http://localhost:8181/servers/")
		if err != nil {
			log.Fatalf("Failed to connect to the etcd API.")
		}
		defer response.Body.Close()

		if response.StatusCode == http.StatusOK {
			fmt.Println("Successsfully onnected with API. Interaction with etcd can be done.")
		} else {
			fmt.Println("Failed to interact with the API.")
		}
	},
}

func main() {
	if err := rootITLDIMS.Execute(); err != nil {
		log.Fatal(err)
	}
}
