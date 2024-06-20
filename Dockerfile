# Start from golang base image
FROM golang:1.20-alpine

# Update Alpine
RUN apk update

# Install utils packages.
RUN apk add --no-cache git

RUN apk add --no-cache build-base

RUN git --version
# Set the current working directory inside the container
WORKDIR /go/src/github.com/jordanlanch/stori-test

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed
RUN go mod download

RUN go get -v github.com/bokwoon95/wgo@latest && go install github.com/bokwoon95/wgo@latest

ENV GIN_MODE="debug"

EXPOSE 8080

ENTRYPOINT ["wgo", "run", "-buildvcs=false", "main.go"]
