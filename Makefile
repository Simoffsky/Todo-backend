.PHONY: lint

run: build
	./bin/app

build:
	go build -o ./bin/app cmd/todo-app/main.go

run-auth: build-auth
	./bin/auth

build-auth:
	go build -o ./bin/auth cmd/auth/main.go

lint:
	golangci-lint run

test:
	go test ./... -race

