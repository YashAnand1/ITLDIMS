package main

import (
	"github.com/spf13/cobra"
)

var (
	itldims = &cobra.Command{
		Use:   "itldims",
		Short: "For checking connectivity with ETCD API",
		Long:  "For checking connectivity - lets user know if connected or not",
		Run:   checkConnection,
	}
)

func checkConnection() {
	fmt.printf("yo")
}
