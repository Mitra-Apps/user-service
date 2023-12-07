# USER-SERVICE

## Install docker
If use wsl : https://docs.docker.com/engine/install/ubuntu/

## How to run
run the apps using command : 
go mod tidy
go mod vendor
sudo docker compose up --build

## Reset database structures
run : sudo docker-compose down --volumes

## Generate pb file from proto file
Run : buf generate