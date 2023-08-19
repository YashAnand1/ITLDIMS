package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	aip := flag.String("aip", "", "Server IP for attribute retrieval")
	attribute := flag.String("attribute", "", "Attribute name to retrieve")
	flag.Parse()

	if *aip == "" || *attribute == "" {
		log.Fatal("Server IP and attribute flags are required.")
	}

	fmt.Printf("Fetching attribute %s from server IP %s\n", *attribute, *aip)
	fmt.Println("Attribute value: DummyValue") // Replace this with your actual etcd interaction logic
}
