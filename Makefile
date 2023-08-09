build:
	go build -o ./bin/app

run-dev:build
	./bin/app --environment dev

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
