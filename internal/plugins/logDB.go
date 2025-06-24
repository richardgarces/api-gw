package plugins

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoLogPlugin struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func (p *MongoLogPlugin) init(config map[string]interface{}) error {
	uri := os.Getenv("MONGO_URI")
	db := os.Getenv("MONGO_DATABASE")
	if config != nil {
		if v, ok := config["uri"].(string); ok && v != "" {
			uri = v
		}
		if v, ok := config["database"].(string); ok && v != "" {
			db = v
		}
	}
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	p.client = client
	p.collection = client.Database(db).Collection("logs")
	return nil
}

func (p *MongoLogPlugin) Wrap(next http.Handler, config map[string]interface{}) http.Handler {
	if p.collection == nil {
		_ = p.init(config)
	}
	saveReqHeaders := false
	saveReqBody := false
	saveRespHeaders := false
	saveRespBody := false
	var headerFilter []string
	if config != nil {
		if v, ok := config["request_headers"].(bool); ok {
			saveReqHeaders = v
		}
		if v, ok := config["request_body"].(bool); ok {
			saveReqBody = v
		}
		if v, ok := config["response_headers"].(bool); ok {
			saveRespHeaders = v
		}
		if v, ok := config["response_body"].(bool); ok {
			saveRespBody = v
		}
		if v, ok := config["header_filter"].([]interface{}); ok {
			for _, h := range v {
				if hs, ok := h.(string); ok {
					headerFilter = append(headerFilter, strings.ToLower(hs))
				}
			}
		}
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &bodyStatusWriter{ResponseWriter: w, status: 200}
		var reqBody []byte
		if saveReqBody && r.Body != nil {
			bodyBytes, _ := io.ReadAll(r.Body)
			reqBody = bodyBytes
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
		next.ServeHTTP(sw, r)
		entry := map[string]interface{}{
			"timestamp":  time.Now(),
			"remote_ip":  r.RemoteAddr,
			"method":     r.Method,
			"path":       r.URL.Path,
			"status":     sw.status,
			"user_agent": r.UserAgent(),
			"referer":    r.Referer(),
			"duration":   time.Since(start).Milliseconds(),
		}
		if saveReqHeaders {
			entry["request_headers"] = filterHeaders(r.Header, headerFilter)
		}
		if saveReqBody {
			entry["request_body"] = string(reqBody)
		}
		if saveRespHeaders {
			entry["response_headers"] = filterHeaders(sw.Header(), headerFilter)
		}
		if saveRespBody {
			entry["response_body"] = string(sw.body)
		}
		go func(e map[string]interface{}) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			_, _ = p.collection.InsertOne(ctx, e)
		}(entry)
	})
}

func filterHeaders(headers http.Header, filter []string) map[string][]string {
	if len(filter) == 0 {
		return headers
	}
	result := make(map[string][]string)
	for _, key := range filter {
		for h, v := range headers {
			if strings.ToLower(h) == key {
				result[h] = v
			}
		}
	}
	return result
}

type bodyStatusWriter struct {
	http.ResponseWriter
	status int
	body   []byte
}

func (w *bodyStatusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *bodyStatusWriter) Write(b []byte) (int, error) {
	w.body = append(w.body, b...)
	return w.ResponseWriter.Write(b)
}

func init() {
	Register("mongo_log", &MongoLogPlugin{})
}
