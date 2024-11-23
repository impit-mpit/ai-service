API_PROTO_FILES=$(shell find api-proto/src/proto/news -name *.proto)
.PHONY: generate
generate:
	protoc --proto_path=./api-proto/src/proto \
		   --proto_path=./api-proto/src/third_party \
 	       --go_out=paths=source_relative:./gen/go \
		   --go-grpc_out=paths=source_relative:./gen/go \
	       $(API_PROTO_FILES)