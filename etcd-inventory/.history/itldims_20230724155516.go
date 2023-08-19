package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	// Etcd host
	etcdHost = "localhost:2379"
)

func getDataByServerIPAttribute(etcdClient *clientv3.Client, serverIP, attribute string) {
	// Construct the etcd key based on the server IP and attribute
	etcdKey := fmt.Sprintf("/servers/%s/%s", serverIP, attribute)

	// Use the etcd client to retrieve the value associated with the etcd key
	response, err := etcdClient.Get(context.Background(), etcdKey)
	if err != nil {
		log.Fatalf("Failed to retrieve data from etcd: %v", err)
	}

	// Check if the key exists and has a value
	if len(response.Kvs) > 0 {
		value := string(response.Kvs[0].Value)
		fmt.Printf("Value for attribute '%s' from server IP '%s': %s\n", attribute, serverIP, value)
	} else {
		fmt.Printf("No data found for attribute '%s' from server IP '%s'\n", attribute, serverIP)
	}
}

func main() {
	// Parse command-line flags
	flag.String("aip", "", "The server IP")
	flag.String("attribute", "", "The attribute to retrieve")
	flag.Parse()

	// Get the values of the flags
	serverIP := flag.Lookup("aip").Value.String()
	attribute := flag.Lookup("attribute").Value.String()

	// Check if the required flags are provided
	if serverIP == "" || attribute == "" {
		log.Println("Usage: itldims --aip <server IP> --attribute <attribute>")
		return
	}

	// Connect to etcd
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints: []string{etcdHost},
	})
	if err != nil {
		log.Fatalf("Failed to connect to etcd: %v", err)
	}
	defer etcdClient.Close()

	// Check if the command is for 'get' and 'aip'
	if serverIP != "" && attribute != "" {
		getDataByServerIPAttribute(etcdClient, serverIP, attribute)
		return
	}

	log.Println("Invalid command. Use 'itldims --help' for usage instructions.")
}
