// package main

// import (
// 	"fmt"
// 	"log"

// 	"go.etcd.io/bbolt"
// )

// func main() {
// 	filePath := "/home/user/my.db/reader/default.etcd/member/snap/db"

// 	db, err := bbolt.Open(filePath, 0666, nil)
// 	if err != nil {
// 		log.Fatalf("Error opening the database: %v", err)
// 	}
// 	defer db.Close()

// 	err = db.View(func(tx *bbolt.Tx) error {

// 		bucketName := "key"
// 		bucket := tx.Bucket([]byte(bucketName))

// 		fmt.Printf("Current Bucket: %s\n\nAll existing buckets:\n", bucketName)

// 		err = tx.ForEach(func(name []byte, _ *bbolt.Bucket) error {
// 			fmt.Println(string(name))
// 			return nil
// 		})

// 		fmt.Printf("\nKey-Value Data From The Bucket Is As Folows\n")

// 		err = bucket.ForEach(func(k, v []byte) error {
// 			fmt.Printf("Key: %s, Value: %s\n", k, v)
// 			return nil
// 		})

// 		return nil
// 	})

// 	if err != nil {
// 		log.Fatalf("Error reading the database: %v", err)
// 	}
// }
