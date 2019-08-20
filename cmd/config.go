package cmd

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Handler struct {
	Method string `yaml:"method"`
	Route  string `yaml:"route"`
	Type   string `yaml:"type"`
	Local  bool   `yaml:"local,omitempty"` // only for listeners
	Path   string `yaml:"path"`            // file path if bin or local listener, host:port to remote listener
	Port   string
}

type Handlers []Handler

func getHandlers(path string) (handlers Handlers, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		danger("Cannot read handlers.yaml file", err)
		return
	}
	err = yaml.Unmarshal(data, &handlers)
	if err != nil {
		danger("Cannot unmarshal handlers.yaml file", err)
		return
	}
	return
}

func (handlers Handlers) getHandler(method, route string) (handler Handler, ok bool) {
	for _, handler = range handlers {
		if handler.Method == method && handler.Route == route {
			ok = true
			return
		}
	}
	ok = false
	return
}
