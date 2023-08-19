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
		Run:   rootRun,
	}

	keyOnly bool
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
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

func rootRun(cmd *cobra.Command, args []string) {
	response, err := http.Get("http://localhost:8181/servers/")
	if err != nil {
		log.Fatalf("Failed to connect to the etcd API.")
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		fmt.Println("Successfully connected with API. Interaction with etcd can be done.")
	}
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().BoolVar(&keyOnly, "listall", false, "Displays only the keys without their values")
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Search Attributes & Values from etcd API",
	Long: `Data retrieval can be made possible 'itldims get <input1> <input2>' or 'itldims get <input1>'.

The possible command combinations that can be used are:
- itldims get <Servers> | Displays attribute values of all Servers
- itldims get <Value> | Displays all Servers containing a specific Attribute value
- itldims get <Attribute> | Displays Displays all Servers containing a specific Attribute
- itldims get <Server IP> | Displays all Attribute values of a specific Server IP
- itldims get <Server Type> <Attribute> | Displays specific Attribute values of a specific Server Type
- itldims get <Server Type> <Value> | Displays all Server Types containing a specific value
- itldims get <Value> <Server Type> | Displays all Server Types containing a specific value
- itldims get <Attribute> <Server IP> | Displays specific Attribute values of a specific Server IP
- itldims get <Server IP> <Attribute> | Displays specific Attribute values of a specific Server IP
- itldims get <Server IP> <Value> | Displays all Server IPs containing a specific value
- itldims get <Value> <Server IP> | Displays all Server IPs containing a specific value
- itldims get <Server IP> <Server Type> | Displays all Attribute values of a specific Server
	`,
	Args: cobra.RangeArgs(1, 2),
	Run:  getRun,
}

func getRun(cmd *cobra.Command, args []string) {
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

			if keyOnly {
				fmt.Println(key)
			} else {
				fmt.Println(key)
				lines := strings.Split(value, "\n")
				for _, line := range lines {
					fmt.Println(line)
				}
			}

			fmt.Println()
		}
	}
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
