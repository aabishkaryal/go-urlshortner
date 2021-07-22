package main

import (
	"flag"
	"fmt"
	"net/http"

	urlshort "github.com/aabishkaryal/go-urlshortner"
)

func main() {
	jsonFileName := flag.String("json",
		"paths.json",
		"JSON file mapping paths to urls in format:\n [{path:pathValue, url:urlValue}, {path:pathValue, url:urlValue}, ...]\n")
	yamlFileName := flag.String("yaml",
		"paths.yaml",
		"Yaml file mapping paths to urls in format:\n-path:path\n url:url\n-path:path\n url:url\n ....\n")
	flag.Parse()

	mux := defaultMux()
	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}

	jsonFileHandler, err := urlshort.JSONFileHandler(*jsonFileName, yamlHandler)
	if err != nil {
		panic(err)
	}

	yamlFileHandler, err := urlshort.YAMLFileHandler(*yamlFileName, jsonFileHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlFileHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
