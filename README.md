# USER-SERVICE

## Install docker
If use wsl : https://docs.docker.com/engine/install/ubuntu/

## How to run
run the apps using command : 
go mod tidy
go mod vendor
sudo docker compose up --build

## Generate pb file from proto file
### Install buf
https://buf.build/docs/installation
If failed, run : brew install buf

### generate protobuf
Run : buf generate

## Reset database structures (Dont run this! Only if needed)
run : sudo docker compose down --volumes