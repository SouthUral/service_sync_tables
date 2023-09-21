FROM golang:latest

WORKDIR /usr/src/service/sync_tables

COPY . .

RUN go build -o /usr/src/service/sync_tables main.go

# USER asd:asd