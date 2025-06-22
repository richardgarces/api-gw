build:
    go build -o api-gw ./cmd/main.go

run: build
    ./api-gw

docker-build:
    docker build -t api-gw .

docker-run:
    docker run -p 8080:8080 api-gw

test:
    go test ./internal/...