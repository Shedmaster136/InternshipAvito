FROM golang:alpine AS builder
COPY . /app
WORKDIR /app
RUN go mod download
RUN go build -o main Service.go
ENTRYPOINT ./main
FROM alpine:latest
LABEL Name=WalletServer Version=1.0
EXPOSE 9090
