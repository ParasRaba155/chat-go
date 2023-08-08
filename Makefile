GO := $(HOME)/go/bin/go1.20.5

build:
	$(GO) build -o ./bin/app

run:build
	./bin/app

tidy:
	$(GO) mod tidy

vendor:
	$(GO) mod vendor

format:
	$(GO) fmt ./...

gen-data: 
	sqlc generate

run-migrate: 
	$(GO) run migrate/main.go
