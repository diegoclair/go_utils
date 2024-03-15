.PHONY: mocks
mocks:
	@echo "=====> Installing mockgen"
	@go install github.com/golang/mock/mockgen@latest

	@echo "=====> Removing old mocks"
	@rm ./logger/mockgen_logger.go

	@echo "=====> Generating mocks"
	@mockgen -source=logger/logger.go -destination=logger/mockgen_logger.go -package=logger

	@echo "=====> Mocks generated"
