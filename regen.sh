set -e

MOD=$1

echo "protoc regen proto/$MOD.proto..."

protoc -I proto $MOD.proto \
	--experimental_allow_proto3_optional \
	--go_out=./$MOD --go_opt=paths=source_relative \
	--go-grpc_out=./$MOD --go-grpc_opt=paths=source_relative

# unused grpc-gwateway    
#    --grpc-gateway_out=./$MOD --grpc-gateway_opt=logtostderr=true --grpc-gateway_opt=paths=source_relative
