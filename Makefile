
deps:
	cd src/ && dep ensure

build:
	go build -o interpreter ./src/cmd/main.go

test:
	go test -v ./...

