.PHONY: all test regen

regen:
	@echo "protoc regen..."
	@protoc -I proto ./proto/machrpc.proto \
		--experimental_allow_proto3_optional \
		--go_out=./machrpc --go_opt=paths=source_relative \
		--go-grpc_out=./machrpc --go-grpc_opt=paths=source_relative