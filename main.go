package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path"
)

// https://hackernoon.com/simple-http-middleware-with-go-79a4ad62889b
type middleware func(next http.HandlerFunc) http.HandlerFunc

func chainMiddleware(mw ...middleware) middleware {
	return func(final http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			last := final

			for i := len(mw) - 1; i >= 0; i-- {
				last = mw[i](last)
			}

			last(w, r)
		}
	}
}

func main() {
	flag.Parse()
	dir := workdir(flag.Args())
	chain := chainMiddleware(withContentType, withNoCache, withLogging)
	http.HandleFunc("/", chain(http.FileServer(http.Dir(dir)).ServeHTTP))
	log.Print("Invoke server @ http://localhost:3000/")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func withContentType(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var contentType string

		if base := path.Ext(r.URL.Path); base == "" {
			contentType = "text/html"
		} else {
			switch base {
			case ".js":
				contentType = "text/javascript"
			case ".mjs":
				contentType = "text/javascript"
			case ".css":
				contentType = "text/css"
			default:
				contentType = "text/html"
			}
		}

		w.Header().Set("Content-Type", contentType)
		next.ServeHTTP(w, r)
	}
}

func withNoCache(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=0; private; no-cache")
		w.Header().Set("Pragma", "no-cache")
		next.ServeHTTP(w, r)
	}
}

func withLogging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s for %s\n", r.Method, r.URL.String(), w.Header().Get("Content-Type"))
		next.ServeHTTP(w, r)
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
