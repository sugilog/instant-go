package main

import (
	"fmt"
	"os"
)

func main() {
	if dir, err := os.Getwd(); err != nil {
		panic(err)
	} else {
		fmt.Println(dir)
	}
}
