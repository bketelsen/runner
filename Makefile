.PHONY: build-cli
build-cli:
	go build -o runcli ./cli

build-server:
	go build -o runner.

run-server:
	go run .
