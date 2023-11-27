.PHONY: all test regen regen-machrpc regen-mgmt regen-bridge

regen: regen-machrpc regen-mgmt regen-bridge regen-schedule

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

regen-bridge:
	@echo "protoc regen bridge..."
	@protoc -I proto ./proto/bridge.proto \
		--go_out=./bridge --go_opt=paths=source_relative \
		--go-grpc_out=./bridge --go-grpc_opt=paths=source_relative

regen-schedule:
	@echo "protoc regen schedule..."
	@protoc -I proto ./proto/schedule.proto \
		--go_out=./schedule --go_opt=paths=source_relative \
		--go-grpc_out=./schedule --go-grpc_opt=paths=source_relative

test:
	go test -cover ./machrpc ./driver