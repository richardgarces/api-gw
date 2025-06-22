package plugins

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestAPIKeyPlugin(t *testing.T) {
    plugin := &APIKeyPlugin{}
    handler := plugin.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(200)
    }), map[string]interface{}{"key": "testkey"})

    req := httptest.NewRequest("GET", "/", nil)
    req.Header.Set("X-API-Key", "testkey")
    rr := httptest.NewRecorder()
    handler.ServeHTTP(rr, req)
    if rr.Code != 200 {
        t.Errorf("esperado 200, obtuve %d", rr.Code)
    }

    req = httptest.NewRequest("GET", "/", nil)
    req.Header.Set("X-API-Key", "wrong")
    rr = httptest.NewRecorder()
    handler.ServeHTTP(rr, req)
    if rr.Code != 401 {
        t.Errorf("esperado 401, obtuve %d", rr.Code)
    }
}