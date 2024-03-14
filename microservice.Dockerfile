FROM golang:1.22.1

WORKDIR /usr/src/app

RUN git config --global --add safe.directory /usr/src/app
RUN go install github.com/cosmtrek/air@latest

COPY . .
RUN go mod tidy
