build:
	go build -o ./bin/app

run-dev:
	air ## air is hot reloader for golang

tidy:
	go mod tidy

vendor:
	go mod vendor

format:
	go fmt ./...

gen-data: 
	sqlc generate

run-migrate: 
	go run migrate/main.go

install-dev:
	scripts/install.sh
