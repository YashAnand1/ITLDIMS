package main

import (
	"context"               //for managing context
	"encoding/csv"          // for CSV file handling
	"encoding/json"         // helps with JSON operations
	"fmt"                   // Import the fmt package for formatted printing
	"log"                   // Import the log package for logging
	"net/http"              // Import the http package for building HTTP servers and clients
	"os"                    // Import the os package for file operations
	"strings"               // Import the strings package for string manipulation
	"github.com/tealeg/xlsx" // Import the xlsx library for Excel file handling
	clientv3 "go.etcd.io/etcd/client/v3" // Import the etcd client library from its specified location
)

var (
	// File paths
	excelFile = "/home/user/sk/etcd-inventory/etcd.xlsx" // Path to the Excel file
	csvFile   = "/home/user/sk/etcd-inventory/myetcd.csv" // Path to the CSV file
	etcdHost  = "localhost:2379" // Address of the etcd server
)

type ServerData map[string]string // Define a type for storing server data as key-value pairs

func convertExcelToCSV(excelFile, csvFile string) {
	// Open the Excel file
	xlFile, err := xlsx.OpenFile(excelFile)
	if err != nil {
		log.Fatalf("Failed to open Excel file due to the error: %v", err)
	}

	// Create the CSV file
	file, err := os.Create(csvFile)
	if err != nil {
		log.Fatalf("Failed to create CSV file: %v", err)
	}
	defer file.Close()

	// Write data to the CSV file
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Iterate over sheets and rows in the Excel file
	for _, sheet := range xlFile.Sheets {
		for _, row := range sheet.Rows {
			var rowData []string
			for _, cell := range row.Cells {
				text := cell.String()
				rowData = append(rowData, text)
			}

			// Check if the row is empty
			isEmptyRow := true
			for _, field := range rowData {
				if field != "" {
					isEmptyRow = false
					break
				}
			}

			// Skip empty rows
			if !isEmptyRow {
				writer.Write(rowData)
			}
		}
	}
}

func uploadToEtcd() {
	// Connect to etcd
	log.Println("Entered into function")
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints: []string{etcdHost},
	})
	if err != nil {
		log.Fatalf("Failed to connect to etcd: %v", err)
	}
	defer etcdClient.Close()

	// Read the CSV file
	file, err := os.Open(csvFile)
	log.Println("reading file")
	if err != nil {
		log.Fatalf("Failed to open CSV file: %v", err)
	}
	defer file.Close()

	// Parse the CSV file
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read CSV file: %v", err)
	}

	// Iterate over the records and upload to etcd
	headers := records[0]
	for _, record := range records[1:] {
		serverIP := record[0]
		serverType := record[1]
		serverData := make(ServerData)

		// Create server data dictionary
		for i := 2; i < len(headers); i++ {
			header := headers[i]
			value := record[i]
			serverData[header] = value
		}

		// Set key-value pairs in etcd for each data field
		for header, value := range serverData {
			etcdKey := fmt.Sprintf("/servers/%s/%s/%s", serverType, serverIP, header)
			etcdValue := value
			fmt.Println(etcdKey)
			fmt.Println(etcdValue)
			_, err := etcdClient.Put(context.Background(), etcdKey, etcdValue)
			if err != nil {
				log.Printf("Failed to upload key-value to etcd: %v", err)
			}
		}

		// Set key-value pair for server data
		etcdKeyData := fmt.Sprintf("/servers/%s/%s/data", serverType, serverIP)
		etcdValueData, err := json.Marshal(serverData)
		if err != nil {
			log.Printf("Failed to marshal server data: %v", err)
			continue
		}
		_, err = etcdClient.Put(context.Background(), etcdKeyData, string(etcdValueData))
		if err != nil {
			log.Printf("Failed to upload server data to etcd: %v", err)
		}
	}

	log.Println("Server details added to etcd successfully.")
}

// HTTP handler function to get server data from etcd
func getServerData(w http.ResponseWriter, r *http.Request) {
	// Extract the etcd key from the URL path
	etcdKeyData := r.URL.Path[len("/servers/"):] // Remove the "/servers/" prefix

	// Check if the key is empty or if it starts with "/"
	if etcdKeyData == "" || strings.HasPrefix(etcdKeyData, "/") {
		listAll(w, r)
	} else {
		getSpecificKey(w, r)
	}
}

func listAll(w http.ResponseWriter, r *http.Request) {
	log.Printf("response %v", r.URL.Path)
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

	// Get all keys with values using prefix
	response, err := etcdClient.Get(ctx, "/servers/", clientv3.WithPrefix())
	if err != nil {
		log.Printf("Failed to retrieve keys from etcd: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	for _, kv := range response.Kvs {
		fmt.Fprintf(w, "Key: %s\n", string(kv.Key))
		fmt.Fprintf(w, "Value: %s\n", string(kv.Value))
		fmt.Fprintf(w, "----------------------------\n")
	}
}

func getSpecificKey(w http.ResponseWriter, r *http.Request) {
	// Extract the etcd key from the URL path
	log.Printf("response %v", r.URL.Path)

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
	etcdKeyData := fmt.Sprintf(r.URL.Path)

	// Get the revision values for the key
	var revisions int
	response, err := etcdClient.Get(ctx, etcdKeyData, clientv3.WithSort(clientv3.SortByCreateRevision, clientv3.SortAscend))
	log.Printf("response %v", response)
	for _, kv := range response.Kvs {
		revisions =
