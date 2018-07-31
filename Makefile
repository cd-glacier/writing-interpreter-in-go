
deps:
	cd src/ && dep ensure

test:
	go test -v ./...

