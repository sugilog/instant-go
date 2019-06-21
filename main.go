package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	flag.Parse()
	dir := workdir(flag.Args())
	server := http.FileServer(http.Dir(dir))
	http.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		log.Printf("%s %s", request.Method, request.URL.String())
		server.ServeHTTP(response, request)
	})
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
