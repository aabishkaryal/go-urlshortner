package urlshort

import (
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
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallbackHandler http.Handler) (http.HandlerFunc, error) {
	// 1. Parse Yaml
	pathURLs, err := parseYamlToPathURL(yml)
	if err != nil {
		return nil, err
	}
	// 2. Convert slice of pathURL to map
	pathsTOUrls := pathSliceToMapConversion(pathURLs)
	return MapHandler(pathsTOUrls, fallbackHandler), nil
}

func pathSliceToMapConversion(pUrls []pathURL) map[string]string {
	result := make(map[string]string)
	for _, pU := range pUrls {
		result[pU.Path] = pU.URL
	}
	return result
}

func parseYamlToPathURL(data []byte) ([]pathURL, error) {
	var pathURLs []pathURL
	err := yaml.Unmarshal(data, &pathURLs)
	if err != nil {
		return nil, err
	}
	return pathURLs, nil
}

type pathURL struct {
	Path string
	URL  string
}
