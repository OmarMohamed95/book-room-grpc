grpc_gen:
	mkdir -p pb && protoc --go_out=./pb --go_opt=paths=source_relative \
	--go-grpc_out=require_unimplemented_servers=false:./pb --go-grpc_opt=paths=source_relative \
    --proto_path=proto proto/*.proto

grpc_clean:
	rm -R pb

server:
	go run cmd/server/main.go -port 8080

client_create:
	go run cmd/client/main.go -address 0.0.0.0:8080 -operation create

client_find:
	go run cmd/client/main.go -address 0.0.0.0:8080 -operation find

client_update:
	go run cmd/client/main.go -address 0.0.0.0:8080 -operation update

client_delete:
	go run cmd/client/main.go -address 0.0.0.0:8080 -operation delete

client_book:
	go run cmd/client/main.go -address 0.0.0.0:8080 -operation book

client_upload:
	go run cmd/client/main.go -address 0.0.0.0:8080 -operation upload

.PHONY: grpc_gen grpc_clean server client