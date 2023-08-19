package main

	import (
		"fmt"
		"io"
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
	
	func main() {
		if err := itldims.Execute(); err != nil {
			log.Fatal(err)
		}
	}
	
)

func fetchDataFromEtcdAPI() (map[string]string, error) {
	response, err := http.Get("http://localhost:8181/servers/")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the etcd API: %v", err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data from the etcd API. Status code: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return parseKeyValuePairs(string(body)), nil
}

func parseKeyValuePairs(data string) map[string]string {
	result := make(map[string]string)

	keyValuePairs := strings.Split(data, "Key:")

	for _, kv := range keyValuePairs {
		kv = strings.TrimSpace(kv)
		if len(kv) == 0 {
			continue
		}

		lines := strings.Split(kv, "Value:")
		if len(lines) == 2 {
			key := strings.TrimSpace(lines[0])
			value := strings.TrimSpace(lines[1])
			result[key] = value
		}
	}

	return result
}

func init() {
	itldims.AddCommand(get)
	get.PersistentFlags().String("all", "", "For displaying all the values and components")
}

func main() {
	if err := itldims.Execute(); err != nil {
		log.Fatal(err)
	}
}
