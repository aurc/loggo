pre-req:
	brew install protobuf
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest


gen:
	protoc --proto_path=protobuf protobuf/loggo/*.proto \
		--go_out=. \
		--go-grpc_out=. \
		--grpc-gateway_out=.
	(cd web/loggo && npm run proto-gen)
