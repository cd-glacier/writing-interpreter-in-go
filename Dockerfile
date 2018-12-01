FROM golang:1.10

WORKDIR /go/src/github.com/g-hyoga/writing-interpreter-in-go
COPY . .

RUN go get -u github.com/golang/dep/...
RUN dep init

ENV GOOS linux
ENV GOARCH amd64

CMD ["go", "run", "src/cmd/main.go"]
