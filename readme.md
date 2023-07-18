# My first GOLang project

Comando para compilação do protobuf

```shell
protoc --go_out=./ --go-grpc_out=./ contracts/*.proto
```

protoc -I ./contracts \
   --go_out ./contracts --go_opt paths=source_relative \
   --go-grpc_out ./contracts --go-grpc_opt paths=source_relative \
   ./contracts/*.proto
   --grpc-gateway_out ./contracts --grpc-gateway_opt paths=source_relative \

First Commit Reference: 
https://dev.to/thenicolau/introducao-ao-grpc-golang-210f
https://dev.to/bruc3mackenzi3/debugging-go-inside-docker-using-vscode-4f67


https://grpc.io/docs/protoc-installation/


Authorization and Authentication Reference
https://dev.to/techschoolguru/use-grpc-interceptor-for-authorization-with-jwt-1c5h