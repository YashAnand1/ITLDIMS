package main

import (
	"context"       //package
	"encoding/csv"  //package
	"encoding/json" //package
	"flag"          //package
	"fmt"           //package
	"log"           //package
	"net/http"      //package
	"os"            //package
	"strings"

	"github.com/tealeg/xlsx"             //package
	clientv3 "go.etcd.io/etcd/client/v3" //package

	"github.com/spf13/cobra"
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
			//fmt.Println(value)
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

func getServerData(w http.ResponseWriter, r *http.Request) {
	// Extract the server type and IP from the URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}

	serverType := parts[2]
	serverIP := parts[3]

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
	etcdKeyData := fmt.Sprintf("/servers/%s/%s/data", serverType, serverIP)
	response, err := etcdClient.Get(ctx, etcdKeyData)
	if err != nil {
		log.Printf("Failed to retrieve data from etcd: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if len(response.Kvs) == 0 {
		http.Error(w, "No data found for the server", http.StatusNotFound)
		return
	}

	// Get the latest revision
	latestRevision := response.Kvs[0].ModRevision

	// Fetch the last 5 revisions for the key
	var revisions []int64
	for i := latestRevision; i > latestRevision-5; i-- {
		revisions = append(revisions, i)
	}

	// Get values for the specified revisions
	values, err := getRevisionValues(etcdClient, etcdKeyData, revisions)
	if err != nil {
		log.Printf("Failed to retrieve data from etcd: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Prepare the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(values)
}

func getRevisionValues(client *clientv3.Client, key string, revisions []int64) ([]string, error) {
	ctx := context.TODO()

	var values []string

	for _, rev := range revisions {
		response, err := client.Get(ctx, key, clientv3.WithRev(rev))
		if err != nil {
			return nil, err
		}

		if len(response.Kvs) > 0 {
			value := string(response.Kvs[0].Value)
			values = append(values, value)
		}
	}

	return values, nil
}

func main() {
	// Convert Excel to CSV
	convertExcelToCSV(excelFile, csvFile)
	log.Println("Excel file converted to CSV successfully.")

	// Parse command-line flags
	flag.Parse()

	// Upload CSV data to etcd
	uploadToEtcd()

	// Create the 'itldims' command using Cobra
	var rootCmd = &cobra.Command{Use: "itldims"}

	var getCmd = &cobra.Command{
		Use:   "get [serverType] [serverIP]",
		Short: "Retrieve data from the etcd API for a specific server",
		Args:  cobra.ExactArgs(2),
		Run:   getServerDataCommand,
	}

	rootCmd.AddCommand(getCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
