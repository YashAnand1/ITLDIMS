//Github version (older) 10 august

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "itldims",
		Short: "Interact with the etcd API",
		Long:  "A command-line tool to interact with the etcd API and check connection",
		Run: func(cmd *cobra.Command, args []string) {
			response, err := http.Get("http://localhost:8181/servers/")
			if err != nil {
				log.Fatalf("Failed to connect to the etcd API.")
			}
			defer response.Body.Close()

			if response.StatusCode == http.StatusOK {
				fmt.Println("Connected to API. Interaction with etcd can be done.")
			}
		},
	}

	keyOnly bool
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Search Attributes & Values from etcd API",
	Long: `Data retrieval: 'itldims get <input1> <input2>' or 'itldims get <input1>'.

Commands:
- itldims get <Servers>          		| Display attributes of all Servers
- itldims get <Value>            		| Display Servers with a specific Attribute value
- itldims get <Attribute>        		| Display Servers with a specific Attribute
- itldims get <Server IP>        		| Display Attribute values of a specific Server IP
- itldims get <Server Type> <Attribute> | Display specific Attribute values of a Server Type
- itldims get <Server Type> <Value>     | Display Server Types containing a specific value
- itldims get <Value> <Server Type>     | Display Server Types containing a specific value
- itldims get <Attribute> <Server IP>   | Display specific Attribute values of a Server IP
- itldims get <Server IP> <Attribute>   | Display specific Attribute values of a Server IP
- itldims get <Server IP> <Value>       | Display Server IPs containing a specific value
- itldims get <Value> <Server IP>       | Display Server IPs containing a specific value
- itldims get <Server IP> <Server Type> | Display Attribute values of a specific Server
	`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		data, err := fetchDataFromAPI()
		if err != nil {
			log.Fatalf("Failed to fetch data from the etcd API: %v", err)
		}

		if len(args) == 1 {
			args = append(args, "servers")
		}

		for key, value := range parseKeyValuePairs(data) {
			if strings.Contains(key, args[0]) || strings.Contains(value, args[0]) {
				if len(args) > 1 && !strings.Contains(key, args[1]) && !strings.Contains(value, args[1]) {
					continue
				}

				if shouldIgnore(key, value) {
					continue
				}

				fmt.Println(key)
				if !keyOnly {
					fmt.Println(value)
				}

				fmt.Println()
			}
		}
	},
}

func fetchDataFromAPI() (string, error) {
	response, err := http.Get("http://localhost:8181/servers/")
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch data. Status code: %d", response.StatusCode)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func shouldIgnore(key, value string) bool {
	ignoreKeywords := []string{"data", "{", "}"}
	for _, keyword := range ignoreKeywords {
		if strings.Contains(key, keyword) || strings.Contains(value, keyword) {
			return true
		}
	}
	return false
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
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().BoolVar(&keyOnly, "listall", false, "Display only keys without values")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
