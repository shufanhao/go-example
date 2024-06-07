
# protobuf-example

## Practical Golang: Using Protobuffs

In this example I use Marshal and Unmarshal using protobuf instead of json.

## generate pb.go 
brew install protoc
`protoc --go_out=. ./service/clientStructure.proto`

## server

`go run server/server.go`

## client

`go run client/client.go`

### From:
https://jacobmartins.com/2016/05/24/practical-golang-using-protobuffs/




