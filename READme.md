# WORDLE with gRPC

#### To generate .pb files
    protoc --go_out=./pbFiles --go_opt=paths=source_relative \
       --go-grpc_out=./pbFiles --go-grpc_opt=paths=source_relative \
       ./protos/wordle.proto