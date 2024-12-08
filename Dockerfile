FROM golang:1.23.4-alpine

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o app ./cmd/main.go

RUN apk add --no-cache curl && \
    curl -L https://github.com/jwilder/dockerize/releases/download/v0.6.1/dockerize-alpine-linux-amd64-v0.6.1.tar.gz | tar -xz -C /usr/local/bin

CMD dockerize -wait tcp://db:5432 -timeout 60s ./app