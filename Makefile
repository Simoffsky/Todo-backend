.PHONY: lint

run: build
	./bin/app

build:
	go build -o ./bin/app cmd/todo-app/main.go

lint:
	golangci-lint run
