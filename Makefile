.PHONY: build
build:
	go build cmd/edu-solution-api/main.go

.PHONY: test
test:
	go test -v ./...