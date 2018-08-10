
deps:
	cd src/ && dep ensure

build:
	go build -o interpreter ./src/cmd/main.go

run:
	go run ./src/cmd/main.go

test:
	go test -v ./...

docker-build:
	docker build -t interpreter .

docker-run:
	docker run -e LOG_LEVEL=debug -v $(PWD):/go/src/github.com/g-hyoga/writing-interpreter-in-go interpreter go run ./src/cmd/main.go
