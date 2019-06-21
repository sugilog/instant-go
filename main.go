package main

import (
	"flag"
	"net/http"
	"os"
)

func main() {
	flag.Parse()
	dir := workdir(flag.Args())
	server := http.FileServer(http.Dir(dir))
	http.Handle("/", server)
	http.ListenAndServe(":3000", nil)
}

func workdir(args []string) string {
	var dir string

	if len(args) == 1 {
		if _, err := os.Stat(args[0]); !os.IsNotExist(err) {
			dir = args[0]
		}
	}

	if dir == "" {
		if d, err := os.Getwd(); err != nil {
			panic(err)
		} else {
			dir = d
		}
	}

	return dir
}
