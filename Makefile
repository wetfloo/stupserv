stupserv: gosrc

.PHONY: run
	go run .

.PHONY: gosrc
gosrc
	go build .

.PHONY: test
test:
	go test -v -race -buildvcs ./...

.PHONY: tidy
tidy:
	go mod tidy -v
	go fmt ./...
