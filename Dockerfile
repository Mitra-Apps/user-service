# Dockerfile definition for Backend application service.

# From which image we want to build. This is basically our environment.
FROM golang:1.21.3-alpine as builder

WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./

# Update Goprivate env
ENV GOPRIVATE="github.com/Mitra-Apps/*"

# Set environment variable for the GitHub Personal Access Token (PAT)
ARG GH_PAT
ENV GH_PAT=${GH_PAT}

# Configure Git and Authentication
RUN apk update && apk add --no-cache git
RUN git config --global credential.helper store \
    && echo "https://${GH_PAT}@github.com" > ~/.git-credentials \
    && git config --global url."https://".insteadOf git:// \
    && echo "machine github.com login x-access-token password ${GH_PAT}" > ~/.netrc \
    && chmod 600 ~/.netrc

# Download dependencies
RUN go mod download

# This will copy all the files in our repo to the inside the container at root location.
COPY . .

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