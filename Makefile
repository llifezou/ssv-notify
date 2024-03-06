SHELL=/usr/bin/env bash

.PHONY: build
build: build
	go mod tidy
	rm -rf ssv-notify
	go build -o ssv-notify main.go
