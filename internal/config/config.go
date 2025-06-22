package config

import (
    "os"
)

type Config struct {
    MongoURI      string
    MongoDatabase string
    HTTPPort      string
}

func Load() *Config {
    return &Config{
        MongoURI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
        MongoDatabase: getEnv("MONGO_DATABASE", "apigw"),
        HTTPPort:      getEnv("HTTP_PORT", "8080"),
    }
}

func getEnv(key, fallback string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return fallback
}