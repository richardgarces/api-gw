package plugins

import (
    "net/http"
)

type ResponseHeadersPlugin struct{}

func (p *ResponseHeadersPlugin) Wrap(next http.Handler, config map[string]interface{}) http.Handler {
    headers, ok := config["set"].(map[string]interface{})
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        rw := &headerWriter{ResponseWriter: w, headers: headers, ok: ok}
        next.ServeHTTP(rw, r)
    })
}

type headerWriter struct {
    http.ResponseWriter
    headers map[string]interface{}
    ok      bool
}

func (w *headerWriter) WriteHeader(statusCode int) {
    if w.ok {
        for k, v := range w.headers {
            w.Header().Set(k, v.(string))
        }
    }
    w.ResponseWriter.WriteHeader(statusCode)
}

func init() {
    Register("response_headers", &ResponseHeadersPlugin{})
}