GO_BUILD_ENV := CGO_ENABLED=0 GOOS=linux GOARCH=amd64
DOCKER_BUILD=$(shell pwd)/.docker_build
DOCKER_CMD=$(DOCKER_BUILD)/go_animapu_web

$(DOCKER_CMD): clean
	mkdir -p $(DOCKER_BUILD)
	$(GO_BUILD_ENV) go build -v -o $(DOCKER_CMD) .

clean:
	rm -rf $(DOCKER_BUILD)

heroku: $(DOCKER_CMD)
	heroku container:push web

cli_app:
	@echo "Start building binaries..."
	@go build -o build/bin/go_animapu_cli cmd/go_animapu_cli/main.go
	@chmod +x build/bin/go_animapu_cli
	@echo "Finish build"

run_cli:
	@./build/bin/go_animapu_cli

web_app:
	@echo "Start building binaries..."
	@go build -o build/bin/go_animapu_web cmd/go_animapu_web/main.go
	@chmod +x build/bin/go_animapu_web
	@echo "Finish build"

run_web:
	@./build/bin/go_animapu_web
