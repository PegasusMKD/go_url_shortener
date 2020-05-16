package main

import (
	"HandlerProject/main/config"
	"HandlerProject/urlshort"
	"flag"
	"fmt"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
	"strings"
)

var databaseConfig string
var db *gorm.DB

func main() {
	dataLocation := flag.String("file", "default.yml", "A file containing all of the shortened paths and urls in YAML or JSON format")
	database := flag.Bool("db", false, "Boolean for saying whether you want to use the local database instead of a file")
	createDatabase := flag.Bool("create", false, "Boolean whether you want to create the data inside of the database")
	flag.Parse()
	mux := defaultMux()

	configuration, err := config.GetConfigFile()
	if err == nil {
		databaseConfig = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
			configuration.DatabaseHost, configuration.DatabasePort, configuration.DatabaseUser, configuration.DatabaseName,
			configuration.DatabasePassword, configuration.SSLMode)
		db, _ = gorm.Open("postgres", databaseConfig)
	}

	if *createDatabase == true {
		urlshort.CreateDB(*dataLocation, db)
	}

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	var mainHandler http.HandlerFunc
	if *database == true {
		mainHandler = urlshort.DBHandler(db, mapHandler)
	} else {
		data, err := ioutil.ReadFile(*dataLocation)
		if err != nil {
			panic(err)
		}
		ext := strings.Split(*dataLocation, ".")[1]
		if ext == "json" {
			mainHandler, err = urlshort.JSONHandler(data, mapHandler)
		} else if ext == "yml" || ext == "yaml" {
			mainHandler, err = urlshort.YAMLHandler(data, mapHandler)
		}
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", mainHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}