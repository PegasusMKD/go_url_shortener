package urlshort

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func buildMap(pathUrls []pathUrl) map[string]string {
	pathsToUrls := make(map[string]string)
	for _,v := range pathUrls {
		pathsToUrls[v.Path] = v.URL
	}
	return pathsToUrls
}

func parseYaml(yml []byte) ([]pathUrl, error) {
	var pathUrls []pathUrl
	err := yaml.Unmarshal(yml, &pathUrls)
	if err != nil {
		return nil, err
	}
	return pathUrls, nil
}

func parseJson(jsonData []byte) ([]pathUrl, error){
	var pathUrls []pathUrl
	err := json.Unmarshal(jsonData, &pathUrls)
	if err != nil {
		return nil, err
	}
	return pathUrls, nil
}

func CreateDB(dataLocation string, db *gorm.DB) {
	file, err := ioutil.ReadFile(dataLocation)
	pathUrls, err := parseJson(file)
	if err != nil {
		panic(err)
	}
	for _, v := range pathUrls {
		db.Create(&v)
	}
}