package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	// File paths
	excelFile = "/home/user/sk/etcd-inventory/etcd.xlsx"
	csvFile   = "/home/user/sk/etcd-inventory/myetcd.csv"
	etcdHost  string
)

type ServerData map[string]string

func init() {
	flag.StringVar(&etcdHost, "etcd-host", "localhost:2379", "Etcd server host")
	flag.Parse()
}

func convertExcelToCSV(excelFile, csvFile string) {
	// ... (same as before) ...
}

func uploadToEtcd() {
	// ... (same as before) ...
}

func getServerData(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		http.Error(w, "Invalid request URL", http.StatusBadRequest)
		return
	}

	serverIP := parts[3]
	attribute := parts[4]

	// Connect to etcd
	ctx := context.TODO()
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints: []string{etcdHost},
	})
	if err != nil {
		log.Printf("Failed to connect to etcd: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer etcdClient.Close()

	// Construct the etcd key for the server data
	etcdKeyData := fmt.Sprintf("/servers/%s/data", serverIP)

	response, err := etcdClient.Get(ctx, etcdKeyData, clientv3.WithPrefix())
	if err != nil {
		log.Printf("Failed to get server data from etcd: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var serverData ServerData
	for _, kv := range response.Kvs {
		err := json.Unmarshal(kv.Value, &serverData)
		if err != nil {
			log.Printf("Failed to unmarshal server data: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		value, found := serverData[attribute]
		if found {
			fmt.Fprintf(w, "%s\n", value)
			return
		}
	}

	http.NotFound(w, r)
}

func main() {
	// Convert Excel to CSV
	convertExcelToCSV(excelFile, csvFile)
	log.Println("Excel file converted to CSV successfully.")

	// Parse command-line flags
	flag.Parse()

	// Upload CSV data to etcd
	uploadToEtcd()

	// Start API server
	log.Println("Starting API server...")
	http.HandleFunc("/servers/", getServerData)
	err := http.ListenAndServe(":8181", nil)
	if err != nil {
		log.Fatalf("Failed to start API server: %v", err)
	}
}
