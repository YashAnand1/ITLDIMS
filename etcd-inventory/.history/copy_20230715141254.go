package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/tealeg/xlsx"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	// File paths
	excelFile = "/home/user/sk/etcd-inventory/etcd.xlsx"
	csvFile   = "/home/user/sk/etcd-inventory/myetcd.csv"
	etcdHost  = "localhost:2379"
)

type ServerData map[string]string

func convertExcelToCSV(excelFile, csvFile string) {
	// Open the Excel file
	xlFile, err := xlsx.OpenFile(excelFile)
	if err != nil {
		log.Fatalf("Failed to open Excel file: %v", err)
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
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints: []string{etcdHost},
	})
	if err != nil {
		log.Fatalf("Failed to connect to etcd: %v", err)
	}
	defer etcdClient.Close()

	// Read the CSV file
	file, err := os.Open(csvFile)
	if err != nil {
		log.Fatalf("Failed to open CSV file: %v", err)a, string(etcdValueData))
		if err != nil {
			log.Printf("Failed to upload server data to etcd: %v", err)
		}
	}

	log.Println("Server details added to etcd successfully.")
}

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

			// Store the creation time along with the value
			timestamp := time.Now().Format(time.RFC3339)
			etcdValueWithTimestamp := fmt.Sprintf("%s - %s", etcdValue, timestamp)

			_, err := etcdClient.Put(context.Background(), etcdKey, etcdValueWithTimestamp)
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

func getServerData(w http.ResponseWriter, r *http.Request) {
	// Extract the etcd key from the URL path
	etcdKeyData := r.URL.Path

	// Connect to etcd
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints: []string{etcdHost},
	})
	if err != nil {
		log.Printf("Failed to connect to etcd: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer etcdClient.Close()

	// Retrieve the key-value pairs with timestamps
	response, err := etcdClient.Get(context.Background(), etcdKeyData, clientv3.WithPrefix())
	if err != nil {
		log.Printf("Failed to get key-value pairs from etcd: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")

	for _, kv := range response.Kvs {
		value := string(kv.Value)
		parts := strings.Split(value, " - ")
		if len(parts) > 1 {
			creationTime := parts[1]
			fmt.Fprintf(w, "Key: %s\n", string(kv.Key))
			fmt.Fprintf(w, "Value: %s\n", parts[0])
			fmt.Fprintf(w, "Creation Time: %s\n", creationTime)
			fmt.Fprintln(w, "-------------------")
		}
	}
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
