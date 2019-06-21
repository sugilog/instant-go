package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path"
)

func main() {
	flag.Parse()
	dir := workdir(flag.Args())
	server := http.FileServer(http.Dir(dir))
	http.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		log.Printf("%s %s\n", request.Method, request.URL.String())
		response.Header().Set("Accept", detectMimeType(request.URL.Path))
		server.ServeHTTP(response, request)
	})
	http.ListenAndServe(":3000", nil)
}

func detectMimeType(pathstring string) string {
	if base := path.Ext(pathstring); base == "" {
		return "plain/text"
	} else {
		switch base {
		case ".js":
			return "text/javascript"
		case ".mjs":
			return "text/javascript"
		case ".css":
			return "text/css"
		default:
			return "plain/html"
		}
	}
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
