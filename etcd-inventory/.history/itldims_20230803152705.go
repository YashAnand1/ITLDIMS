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
	itldims = &cobra.Command{
		Use:   "itldims",
		Short: "Interact with the etcd API",
		Long:  "A command-line tool to interact with the etcd API and tell if the connection has been made",
	}

	getServersCmd = &cobra.Command{
		Use:   "servers",
		Short: "Get keys with specific inputs from etcd API",
		Args:  cobra.ExactArgs(2),
		Run:   getServersData}
)

func init() {
	itldims.AddCommand(getServersCmd)
}

func getServersData(cmd *cobra.Command, args []string) {
	data, err := fetchetcddata()
	if err != nil {
		log.Fatalf("Failed to fetch data from the etcd API: %v", err)
	}

	for key, value := range data {
		if containsAll(key, args[0], args[1]) {
			fmt.Printf("\n %s\n %s\n", key, value)
		}
	}
}

func fetchetcddata() (map[string]string, error) {
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

	return parsepairs(string(body)), nil
}

func parsepairs(data string) map[string]string {
	result := make(map[string]string)
	lines := strings.Split(data, "\n")
	for i := 0; i < len(lines)-1; i += 2 {
		result[strings.TrimSpace(lines[i])] = strings.TrimSpace(lines[i+1])
	}
	return result
}

func containsAll(s, substr1, substr2 string) bool {
	return strings.Contains(s, substr1) && strings.Contains(s, substr2)
}

func main() {
	if err := itldims.Execute(); err != nil {
		log.Fatal(err)
	}
}
