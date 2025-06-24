package proxy

import (
	"api-gw/internal/models"
	"api-gw/internal/plugins"
	"api-gw/internal/routes"
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"
)

type contextKey string

func NewReverseProxy(routeManager *routes.RouteManager) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route, service, target, err := routeManager.ResolveWithService(r.URL.Path)
		if err != nil {
			entry := map[string]interface{}{
				"level":     "warn",
				"timestamp": time.Now().Format(time.RFC3339),
				"event":     "route_not_found",
				"method":    r.Method,
				"path":      r.URL.Path,
				"remote_ip": r.RemoteAddr,
			}
			_ = json.NewEncoder(os.Stdout).Encode(entry)
			http.Error(w, "Ruta no encontrada 3", http.StatusNotFound)
			return
		}
		// Guarda el target en el contexto para el logger
		ctx := context.WithValue(r.Context(), contextKey("target"), target)
		r = r.WithContext(ctx)
		url, _ := url.Parse(target)
		proxy := httputil.NewSingleHostReverseProxy(url)
		proxy.Transport = &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second, // Timeout de conexión
				KeepAlive: 10 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   5 * time.Second,
			ResponseHeaderTimeout: 5 * time.Second, // Timeout de respuesta
			ExpectContinueTimeout: 1 * time.Second,
		}

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
