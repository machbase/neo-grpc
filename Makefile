.PHONY: all test regen regen-machrpc regen-mgmt

regen: regen-machrpc regen-mgmt

regen-machrpc:
	@echo "protoc regen machrpc..."
	@protoc -I proto ./proto/machrpc.proto \
		--experimental_allow_proto3_optional \
		--go_out=./machrpc --go_opt=paths=source_relative \
		--go-grpc_out=./machrpc --go-grpc_opt=paths=source_relative

regen-mgmt:
	@echo "protoc regen mgmt..."
	@protoc -I proto ./proto/mgmt.proto \
		--go_out=./mgmt --go_opt=paths=source_relative \
		--go-grpc_out=./mgmt --go-grpc_opt=paths=source_relative