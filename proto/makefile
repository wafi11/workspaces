gen-auth:
	protoc --proto_path=. \
	--go_out=../core/auth-services/gen/v1 \
	--go_opt=paths=source_relative \
	--go-grpc_out=. \
	--go-grpc_opt=paths=source_relative \
	auth.proto

gen-user:
	protoc --proto_path=. \
	--go_out=../core/user-service/gen/v1 \
	--go_out=../core/api-gateway/gen/v1 \
	--go_opt=paths=source_relative \
	--go-grpc_out=. \
	--go-grpc_opt=paths=source_relative \
	user.proto

gen-storage:
	protoc --proto_path=. \
	--go_out=../core/storage-service/gen/v1 \
	--go_out=../core/api-gateway/gen/v1 \
	--go_opt=paths=source_relative \
	--go-grpc_out=. \
	--go-grpc_opt=paths=source_relative \
	storage.proto