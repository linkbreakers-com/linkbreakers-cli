.PHONY: build generate docs test

build:
	go build -o dist/linkbreakers ./cmd/linkbreakers

generate:
	./scripts/generate-commands.sh

docs:
	go run ./cmd/linkbreakers gendocs

test:
	go test ./...
