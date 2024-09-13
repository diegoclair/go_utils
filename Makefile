.PHONY: mocks
mocks:
	@echo "=====> Installing mockgen"
	@go install github.com/golang/mock/mockgen@latest

	@echo "=====> Removing old mocks"
	@rm ./logger/mockgen_logger.go

	@echo "=====> Generating mocks"
	@mockgen -source=logger/logger.go -destination=logger/mockgen_logger.go -package=logger

	@echo "=====> Mocks generated"

.PHONY: protoc
protoc:
	@echo "=====> Generating protobuf"
	@rm -rf ./resterrors/internal/pb/*
	@protoc -I=resterrors/internal/protodefs \
		--go_out=./resterrors/internal/pb \
		--go_opt=paths=source_relative \
		--go-grpc_out=./resterrors/internal/pb \
		--go-grpc_opt=paths=source_relative \
		resterrors/internal/protodefs/error.proto
	@echo "=====> Protobuf generated"