FROM golang:alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN apk add --no-cache bash

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

RUN go build -o main simple-http-server/main.go

WORKDIR /dist

RUN cp /build/entrypoint.sh .
RUN cp /build/main .

ENTRYPOINT ["/dist/entrypoint.sh"]