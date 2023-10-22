FROM golang:1.21.3

WORKDIR /usr/src/server

RUN go install github.com/cosmtrek/air@latest

COPY . .

RUN go mod tidy
