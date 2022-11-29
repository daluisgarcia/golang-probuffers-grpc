# Protobuf and gRPC project in Golang

Protobuffers is a binary serialization format for structured data. It is language and platform neutral. It is used to serialize structured data for use in communications protocols, data storage, and more.

The protoc compiler is used to generate code from .proto files. The generated code can be used to populate, serialize, and retrieve structured data.

The gRPC project is a modern, open source, high-performance remote procedure call (RPC) framework that can run anywhere. It enables client and server applications to communicate transparently, and makes it easier to build connected systems.


## Installation 
First you need to install the protoc compiler. You can download it from [here](https://github.com/protocolbuffers/protobuf/releases). Be sure to download the correct version for your OS.

Then you need to install the protoc-gen-go plugin. You can run the command ```go install google.golang.org/protobuf/cmd/protoc-gen-go@latest``` to install the latest version of the compiler.

Also you need to install the protoc-gen-go-grpc plugin. You can run the command ```go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest``` to install the latest version of the plugin.

## Usage
To generate the gRPC code from the proto file to a .go file, you can run the command ```protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative <relative-route-to-.proto-file>```.

To run the **Student service** run ```go run server-student/main.go```. This one will run in port 5060.

To run the **Test service** run ```go run server-test/main.go```. This one will run in port 5061.

To run the **Client**, uncomment/comment the functions depending on what type of connection you want and then run ```go run client/main.go```.