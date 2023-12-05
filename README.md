# USER-SERVICE

## Installation
### Create new database
install postgre
create new database using postgre and add the connection information to .env file
run this sql query to add GUID datatype : CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

### Install protobuf
Install go version 1.21.3
Install protobuf : https://github.com/protocolbuffers/protobuf/releases/tag/v25.1
Execute command line in terminal :
go get -u google.golang.org/protobuf
go get -u google.golang.org/grpc       
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go mod tidy
go mod vendor

## How to run
run the apps using command : go run main.go

## Generate pb file from proto file
Run this command :
buf generate