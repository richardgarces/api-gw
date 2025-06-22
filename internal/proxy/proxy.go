package proxy

import (
	"api-gw/internal/models"
	"api-gw/internal/plugins"
	"api-gw/internal/routes"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewReverseProxy(routeManager *routes.RouteManager) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route, service, target, err := routeManager.ResolveWithService(r.URL.Path)
		if err != nil {
			http.Error(w, "Ruta no encontrada", http.StatusNotFound)
			return
		}
		url, _ := url.Parse(target)
		proxy := httputil.NewSingleHostReverseProxy(url)

		// Aplica plugins de ruta y servicio (ruta tiene prioridad)
		var pluginsToApply []models.PluginConfig
		if service != nil {
			pluginsToApply = append(pluginsToApply, service.Plugins...)
		}
		if route != nil {
			pluginsToApply = append(pluginsToApply, route.Plugins...)
		}
		handler := ApplyPlugins(proxy, pluginsToApply)
		handler.ServeHTTP(w, r)
	})
}

func ApplyPlugins(handler http.Handler, pluginConfigs []models.PluginConfig) http.Handler {
	for i := len(pluginConfigs) - 1; i >= 0; i-- {
		plugin := plugins.Get(pluginConfigs[i].Type)
		if plugin != nil {
			handler = plugin.Wrap(handler, pluginConfigs[i].Config)
		}
	}
	return handler
}
