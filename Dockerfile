FROM golang:alpine

COPY . /src/
WORKDIR /src
RUN go build -o /go/bin/publisher cmd/publisher/main.go
RUN go build -o /go/bin/worker cmd/worker/main.go

