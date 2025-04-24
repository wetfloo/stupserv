http_server: gosrc

.PHONY: run
	go run .

gosrc: *.go
	go build .

.PHONY: test
test:
	go test -v -race -buildvcs ./...

.PHONY: tidy
tidy:
	go mod tidy -v
	go fmt ./...
