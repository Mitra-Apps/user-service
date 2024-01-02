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
RUN GOPATH= go build -o myapp

####################################################################
# This is the actual image that we will be using in production.
FROM alpine:latest

WORKDIR /app

# We need to copy the binary from the build image to the production image.
COPY --from=builder /app/myapp .

# This is the port that our application will be listening on.
EXPOSE 9000

# This is the command that will be executed when the container is started.
ENTRYPOINT ["./myapp"]