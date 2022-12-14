run:
	docker-compose --env-file ./config/.env.dev up --build

proto:
	protoc --go_out=internal --go_opt=paths=source_relative \
	  --go-grpc_out=internal --go-grpc_opt=paths=source_relative \
	  protos/**/*.proto