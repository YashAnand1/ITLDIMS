package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "itldims",
		Short: "Interact with the etcd API",
		Long:  "A command-line tool to interact with the etcd API and tell if the connection has been made",
	}

	getServersCmd = &cobra.Command{
		Use:   "servers",
		Short: "Get keys with specific inputs from etcd API",
		Args:  cobra.ExactArgs(2),
		Run:   getServersData,
	}
)

func init() {
	rootCmd.AddCommand(getServersCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func getServersData(cmd *cobra.Command, args []string) {
	data, err := fetchDataFromEtcdAPI()
	if err != nil {
		log.Fatalf("Failed to fetch data from the etcd API: %v", err)
	}

	fmt.Printf("Server data with input1='%s' and input2='%s':\n", args[0], args[1])
	fmt.Println("---------------------")

	for key, value := range data {
		if strings.Contains(key, args[0]) && strings.Contains(key, args[1]) {
			fmt.Printf("\n %s\n %s\n", key, value)
		}
	}
}

func fetchDataFromEtcdAPI() (map[string]string, error) {
	response, err := http.Get("http://localhost:8181/servers/")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the etcd API: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data from the etcd API. Status code: %d", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return parseKeyValuePairs(string(body)), nil
}

func parseKeyValuePairs(data string) map[string]string {
	result := make(map[string]string)
	lines := strings.Split(data, "\n")
	for i := 0; i < len(lines)-1; i += 2 {
		result[strings.TrimSpace(lines[i])] = strings.TrimSpace(lines[i+1])
	}
	return result
}
