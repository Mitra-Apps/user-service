# Dockerfile definition for Backend application service.

# From which image we want to build. This is basically our environment.
FROM golang:1.21.3-alpine as builder

WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# This will copy all the files in our repo to the inside the container at root location.
COPY . .

# Perform go mod tidy and go mod vendor
RUN go mod tidy
RUN go mod vendor

# Build our binary at root location.
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOPATH=/go \
    go build -ldflags="-w -s" -o /usr/local/bin/user-service .

####################################################################
# This is the actual image that we will be using in production.
FROM alpine:latest

COPY --from=builder /app/docs /app/docs
COPY --from=builder /usr/local/bin/user-service /usr/local/bin/user-service

# This is the port that our application will be listening on.
EXPOSE 9100
EXPOSE 9101

WORKDIR /app

ENTRYPOINT ["/usr/local/bin/user-service"]