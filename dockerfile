FROM golang:latest

WORKDIR $GOPATH/src/service/

COPY . .

RUN go build -o /go/src/service/sync_tables main.go

USER asd:asd