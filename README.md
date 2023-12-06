# USER-SERVICE

## Installation

### WSL
#### **Database Setup**
1. Install PostgreSQL on your WSL
2. Create new database `sudo -u postgres createdb user-service`
3. Change credentials on `.env` file to match your setup regarding username and password
4. Run this sql query to add GUID datatype : `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`


#### **Install brew**

1. Install brew using this command `/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`
    
    >It will ask you to enter your sudo password  
 
2. Install gcc using `brew install gcc`

#### **Install protobuf**

`brew install buf`

#### **Run Go Mod** 

`go mod tidy`

`go mod vendor`

## How to run
run the apps using command : 

`go run main.go`

## Generate pb file from proto file
Run this command :


`buf generate`

  >It will regenerate domain/proto/user folder
