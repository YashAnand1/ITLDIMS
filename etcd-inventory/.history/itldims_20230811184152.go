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
	itldims = &cobra.Command{
		Use:   "itldims",
		Short: "Interact with the etcd API",
		Run: func(cmd *cobra.Command, args []string) {
			response, err := http.Get("http://localhost:8181/servers/")
			if err != nil {
				log.Fatal("Failed to connect to the etcd API.")
			}
			defer response.Body.Close()

			if response.StatusCode == http.StatusOK {
				fmt.Println("Connected to API. Interaction with etcd can be done.")
			}
		},
	}

	get = &cobra.Command{
		Use:   "get",
		Short: "Search Attributes & Values from etcd API",
		Long:  "Data retrieval can be done using 'itldims get <input1> <input2>' or 'itldims get <input1>'.",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			data, err := fetchDataFromAPI()
			if err != nil {
				log.Fatalf("Failed to fetch data from the etcd API: %v", err)
			}

			pairs := parseKeyValuePairs(data)
			searchValue(args, pairs, data)
		},
	}
)

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

func searchValue(args []string, pairs map[string]string, data string) {
	if len(args) == 1 && args[0] == "servers" {
		for key := range parseKeyValuePairs(data) {
			splitKey := strings.Split(key, "/")
			serverIP := splitKey[3]
			fmt.Printf("%s\n", serverIP)
		}
		return
	}

	for key, value := range pairs {
		if strings.Contains(key, "{") || strings.Contains(key, "}") ||
			strings.Contains(value, "{") || strings.Contains(value, "}") {
			continue
		}

		if !strings.Contains(key, "data") && (strings.Contains(key, args[0]) || strings.Contains(value, args[0])) {
			if len(args) > 1 && !strings.Contains(key, args[1]) && !strings.Contains(value, args[1]) {
				continue
			}

			fmt.Printf("%s:\n%s\n\n", key, value)
		}
	}
}

func init() {
	itldims.AddCommand(get)
}

func main() {
	if err := itldims.Execute(); err != nil {
		log.Fatal(err)
	}
}
