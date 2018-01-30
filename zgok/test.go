package main

import (
	"fmt"
	"github.com/srtkkou/zgok"
	"os"
)

func main() {
	zfs, _ := zgok.RestoreFileSystem(os.Args[0])

	fmt.Printf("hello world paths=%#v\n", zfs.Paths())
}