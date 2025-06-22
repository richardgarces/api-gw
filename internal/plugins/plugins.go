package plugins

import (
	"net/http"
)

type Plugin interface {
	Wrap(http.Handler, map[string]interface{}) http.Handler
}

var registry = make(map[string]Plugin)

func Register(name string, plugin Plugin) {
	registry[name] = plugin
}

func Get(name string) Plugin {
	return registry[name]
}

func Registry() map[string]Plugin {
	return registry
}
