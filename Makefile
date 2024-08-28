GO_PROJECT_NAME := image-upload

go_build:
	@echo "\n....Building $(GO_PROJECT_NAME)...."
	go mod tidy && go mod vendor
	go build -o $(GO_PROJECT_NAME) ./cmd

go_run:
	@echo "\n....Running $(GO_PROJECT_NAME)...."
	./$(GO_PROJECT_NAME)

restart:
	@$(MAKE) go_build
	@$(MAKE) go_run

mocks/generate:
	mockgen -destination=mocks/repositories/upload_link_mock.go -package=mocks -source=src/repositories/upload_link.go UploadLinkRepository
	mockgen -destination=mocks/repositories/image_mock.go -package=mocks -source=src/repositories/image.go ImageRepository
	mockgen -destination=mocks/repositories/statistics_mock.go -package=mocks -source=src/repositories/statistics.go StatisticsRepository
	mockgen -destination=mocks/producers/image_uploaded_mock.go -package=mocks -source=src/producers/image_uploaded.go ImageUploadedProducer