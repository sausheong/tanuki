package cmd

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Executable struct {
	Method string `yaml:"method"`
	Route  string `yaml:"route"`
	File   string `yaml:"file"`
}

type Executables []Executable

type Listener struct {
	Method string `yaml:"method"`
	Route  string `yaml:"route"`
	Local  bool   `yaml:"local"`
	Path   string `yaml:"path"`
	Port   string
}

type Listeners []Listener

func getBins() (bins Executables, err error) {
	data, err := ioutil.ReadFile("bins.yaml")
	if err != nil {
		danger("Cannot read bins.yaml file", err)
		return
	}
	err = yaml.Unmarshal(data, &bins)
	if err != nil {
		danger("Cannot unmarshal bins.yaml file", err)
		return
	}
	return
}

func (bins Executables) getBin(method, route string) (file string, ok bool) {
	for _, bin := range bins {
		if bin.Method == method && bin.Route == route {
			ok = true
			file = bin.File
			return
		}
	}
	ok = false
	return
}

func getListeners() (list Listeners, err error) {
	data, err := ioutil.ReadFile("listeners.yaml")
	if err != nil {
		danger("Cannot read listeners.yaml file", err)
		return
	}
	err = yaml.Unmarshal(data, &list)
	if err != nil {
		danger("Cannot unmarshal listeners.yaml file", err)
		return
	}
	return
}

func (list Listeners) getListener(method, route string) (l Listener, ok bool) {
	for _, listener := range list {
		if listener.Method == method && listener.Route == route {
			ok = true
			l = listener
			return
		}
	}
	ok = false
	return
}
