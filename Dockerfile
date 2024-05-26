FROM golang:1.22.2-alpine

RUN apk add --no-cache postgresql-client

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .


RUN go build -o email-reminder .

CMD dockerize -wait tcp://go_db:5432 -timeout 1m ./email-reminder