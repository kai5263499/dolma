package main

import (
	"fmt"

	"github.com/spf13/afero"
	"io/ioutil"
)

func main() {
	var AppFs = afero.NewMemMapFs()

	var file afero.File
	var err error

	file, err = AppFs.Create("/tmp/somefile")

	file.WriteString("test data")

	file, err = AppFs.Open("/tmp/somefile")
	if err != nil {
		fmt.Printf("err=%s", err)
	}

	data, err := ioutil.ReadAll(file)


	fmt.Printf("data read %s\n", string(data)) // what the *bleep*?
}
