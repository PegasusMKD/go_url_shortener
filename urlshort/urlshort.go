package urlshort

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"net/http"
)

func DBHandler(db *gorm.DB, fallback http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var pathUrls pathUrl
		db.Where(&pathUrl{Path: r.RequestURI}).First(&pathUrls)
		http.Redirect(w, r, pathUrls.URL,http.StatusFound)
		fallback.ServeHTTP(w, r)
	}

}

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if val, ok := pathsToUrls[r.RequestURI]; ok {
			http.Redirect(w, r, val, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}

}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathUrls, err := parseYaml(yml)
	if err != nil {
		return nil, err
	}
	pathsToUrls := buildMap(pathUrls)
	return MapHandler(pathsToUrls, fallback), nil
}

func JSONHandler(jsonData []byte, fallback http.Handler) (http.HandlerFunc, error) {
	data, err := parseJson(jsonData)
	if err != nil {
		panic(err)
	}
	pathsToUrls := buildMap(data)
	return MapHandler(pathsToUrls, fallback), nil
}
