package urlshort

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	yaml "gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallbackHandler http.Handler) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		currentPath := r.URL.Path
		if destination, ok := pathsToUrls[currentPath]; ok {
			http.Redirect(rw, r, destination, http.StatusFound)
			return
		}
		fallbackHandler.ServeHTTP(rw, r)
	}
}

// JSONFileHandler will parse the provided JSON file and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the JSON file, then the
// fallback http.Handler will be called instead.
//
// JSON is expected to be in the format:
//
//	[
// 		{path:pathValue, url:urlValue},
// 		...
// 	]
//
// The  errors that can be returned are related to having
// invalid JSON data or opening json file.
func JSONFileHandler(jsonFileName string, fallbackHandler http.Handler) (http.HandlerFunc, error) {
	// 1. Read JSON file data
	jsonData, err := openFile(jsonFileName)
	if err != nil {
		return nil, err
	}
	// 2. Parse JSON data to slice of pathURLs
	pathURLs, err := parseJSONToPathURL(jsonData)
	if err != nil {
		return nil, err
	}
	// 3. Convert slice of pathURLs to map
	pathsToURLs := pathSliceToMapConversion(pathURLs)
	return MapHandler(pathsToURLs, fallbackHandler), nil
}

func parseJSONToPathURL(data []byte) ([]pathURL, error) {
	var pathsToURLs []pathURL
	err := json.Unmarshal(data, &pathsToURLs)
	if err != nil {
		return nil, err
	}
	return pathsToURLs, nil
}

// YAMLFileHandler will parse the provided YAML file and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML file, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
// 	pathsToURLs:
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The  errors that can be returned are related to having
// invalid YAML data or opening yaml file.
func YAMLFileHandler(yamlFileName string, fallbackHandler http.HandlerFunc) (http.HandlerFunc, error) {
	// 1. Open yaml file
	yamlBytes, err := openFile(yamlFileName)
	if err != nil {
		return nil, err
	}
	return YAMLHandler(yamlBytes, fallbackHandler)
}

// openFile will open the file with fileName and return the content of the file.
func openFile(fileName string) ([]byte, error) {
	fileByte, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return fileByte, nil

}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
func YAMLHandler(yml []byte, fallbackHandler http.Handler) (http.HandlerFunc, error) {
	// 1. Parse Yaml
	pathURLs, err := parseYAMLToPathURL(yml)
	if err != nil {
		return nil, err
	}
	// 2. Convert slice of pathURL to map
	pathsToURLs := pathSliceToMapConversion(pathURLs)
	return MapHandler(pathsToURLs, fallbackHandler), nil
}

func pathSliceToMapConversion(pUrls []pathURL) map[string]string {
	result := make(map[string]string)
	for _, pU := range pUrls {
		result[pU.Path] = pU.URL
	}
	return result
}

func parseYAMLToPathURL(data []byte) ([]pathURL, error) {
	var pathURLs []pathURL
	err := yaml.Unmarshal(data, &pathURLs)
	if err != nil {
		return nil, err
	}
	return pathURLs, nil
}

type pathURL struct {
	Path string `json:"path" yaml:"path"`
	URL  string `json:"url" yaml:"url"`
}
