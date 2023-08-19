package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var (
	itldims = &cobra.Command{
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

	get = &cobra.Command{
		Use:   "get",
		Short: "Search Attributes & Values from etcd API",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			data, err := fetchDataFromEtcdAPI()
			if err != nil {
				log.Fatalf("Failed to fetch data from the etcd API: %v", err)
			}

			if len(args) == 1 {
				args = append(args, "servers")
			}

			for key, value := range data {
				if strings.Contains(key, "{") || strings.Contains(key, "}") ||
					strings.Contains(value, "{") || strings.Contains(value, "}") {
					continue
				}

				arg1 := args[0]
				arg2 := args[1]

				// Use regular expression to match exact arguments
				reArg1 := regexp.MustCompile(fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta(arg1)))
				reArg2 := regexp.MustCompile(fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta(arg2)))

				if !strings.Contains(key, "data") &&
					(reArg1.MatchString(key) || reArg1.MatchString(value)) &&
					(reArg2.MatchString(key) || reArg2.MatchString(value)) {
					fmt.Printf("key=%s\n", key)

					lines := strings.Split(value, "\n")
					for _, line := range lines {
						fmt.Println(line)
					}
					fmt.Println()
				}
			}
		},
	}
)

// ... rest of your code ...

func main() {
	if err := itldims.Execute(); err != nil {
		log.Fatal(err)
	}
}
