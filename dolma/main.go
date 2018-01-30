package main

import (
	"fmt"
	"os"
	"github.com/kai5263499/dolma/dolma/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
