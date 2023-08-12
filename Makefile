grpc_gen:
	mkdir -p pb && protoc --go_out=./pb --go_opt=paths=source_relative \
	--go-grpc_out=require_unimplemented_servers=false:./pb --go-grpc_opt=paths=source_relative \
    --proto_path=proto proto/*.proto

grpc_clean:
	rm -R pb

server:
	go run cmd/server/main.go -port 8080

client:
	go run cmd/client/main.go -address 0.0.0.0:8080

.PHONY: grpc_gen grpc_clean server client