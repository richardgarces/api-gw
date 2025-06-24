# Makefile for Proyecto Web Project
SHELL := /bin/bash
.DEFAULT_GOAL := push
.PHONY: run build docker push clean help

help:
	@echo "Comandos disponibles:"
	@echo "  run      - Compila y ejecuta el binario"
	@echo "  build    - Compila el binario"
	@echo "  docker   - Construye la imagen Docker (rgarces/api-gw:latest)"
	@echo "  push     - Sube la imagen Docker a Docker Hub"
	@echo "  clean    - Limpia binarios y archivos temporales"

build:
	go build -o api-gw ./cmd/main.go

run: build
	./api-gw

docker: build
	docker build --platform linux/amd64 -t rgarces/api-gw:latest .

push: docker
	docker push rgarces/api-gw:latest

clean:
	go clean
	rm -rf api-gw
