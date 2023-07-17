# Wordle-with-gRPC-and-CobraCLI

A Simple wordle game implementation in Golang with gRPC server and CLI client
## Technologies used
- Golang
- Protocol Buffers
- gRPC
- Cobra-CLI
- Git/Github
## Features
- Real time LeaderBoard
    - Hostorical LeaderBoard
    - Today's Leader
- Own InMemoryDB
- gRPC server
- Efficient data communication powered by HTTP/2 protocol
- Simple Client Interface with Cobra-CLI
- Data persistance in Mongo Atlas (Todo)
- JWT authentication (Todo)

## Implentation Details

- Implemented own In Memory Database from the scratch with proper Mutex Locking to avoid data race conditions
- Used Protocol buffers and gRPC architecture for efficient data communication
- Implemented real time leader board with accurate timestamps
- Written Cron Jobs to update word, store and reset leader board daily
- Implemented custom random function to get a random word out of 3000+ words based on date
- Handled all cases in Wordle game

### Notes
#### To Change root Command
- change the name of client file (client -> wordle)
- delete client (old binary) from home/go/bin
- run these commands
  
          go build

          go install

          export PATH=$PATH:$GOPATH/bin


#### To generated stubs
    protoc --go_out=./pbFiles --go_opt=paths=source_relative --go-grpc_out=./pbFiles --go-grpc_opt=paths=source_relative ./protos/wordle.proto
