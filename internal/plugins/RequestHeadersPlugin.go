package plugins

import (
    "net/http"
)

type RequestHeadersPlugin struct{}

func (p *RequestHeadersPlugin) Wrap(next http.Handler, config map[string]interface{}) http.Handler {
    headers, ok := config["set"].(map[string]interface{})
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if ok {
            for k, v := range headers {
                r.Header.Set(k, v.(string))
            }
        }
        next.ServeHTTP(w, r)
    })
}

func init() {
    Register("request_headers", &RequestHeadersPlugin{})
}