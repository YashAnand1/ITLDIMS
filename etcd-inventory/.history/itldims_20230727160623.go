package main

import ( //Importation of packages which help with formatting, logging and interacting with http addresses.
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra" //importation of the cobra tool for creating commands.
	// "go.etcd.io/etcd/pkg/v3/cobrautl"
)

var rootITLDIMS = &cobra.Command{
	Use:   "itldims",                                                                                    //the command
	Short: "Interact with the etcd API",                                                                 //short explanation
	Long:  "A command-line tool to interact with the etcd API and tell if the connection has been made", //longer explanation
	Run: func(cmd *cobra.Command, args []string) { //When function called, the command is represeted by cobra.Command and args is for taking possible flags
		response, err := http.Get("http://localhost:8181/servers/") //http.Get used to connect with localhost:8181
		if err != nil {
			log.Fatalf("Failed to connect to the etcd API.")
		}
		defer response.Body.Close() //to close the connection that was started

		if response.StatusCode == http.StatusOK { //If conection is set up
			fmt.Println("Successsfully onnected with API. Interaction with etcd can be done.")
		} else {
			fmt.Println("Failed to interact with the API.")
		}
	},
}

func main() { //execution of the above rootITLDIMS function is done.
	if err := rootITLDIMS.Execute(); err != nil {
		log.Fatal(err)
	}
}
