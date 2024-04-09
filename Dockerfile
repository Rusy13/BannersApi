# Используйте официальный образ Go как базовый
FROM golang:1.18 as builder

RUN go version
ENV GOPATH=/
COPY ./ ./
RUN go mod download
RUN go build -o avito-app ./cmd/app/main.go
RUN go get -u github.com/pressly/goose/cmd/goose
RUN go install github.com/pressly/goose/cmd/goose



CMD ["./avito-app"]
