FROM spbgit.polymetal.ru:5005/polyna/docker/images/asd-golang:1.2 as builder

WORKDIR $GOPATH/src/service/

COPY . .

RUN go build -o /go/src/service/syncTables main.go

USER asd:asd
