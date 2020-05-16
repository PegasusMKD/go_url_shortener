package config

import (
	"github.com/tkanos/gonfig"
	"path"
	filepath2 "path/filepath"
	"runtime"
)

func GetConfigFile() (Config, error){
	_, dirname ,_, _ := runtime.Caller(0)
	filepath := path.Join(filepath2.Dir(dirname), "config.json")
	configuration := Config{}
	err := gonfig.GetConf(filepath, &configuration)
	return configuration, err
}