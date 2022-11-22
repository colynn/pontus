NAME=pontus
VERSION=1.0.0

.PHONY: build
## build: Compile the packages.
build:
	@go build -o $(NAME) cmd/pontus/main.go

.PHONY: run
## run: Build and Run in local mode.
run: build
	@ENV=local ./$(NAME)

.PHONY: run-dev
## run-dev: Build and Run in development mode.
run-dev: build
	@ENV=dev GIN_MODE=test ./$(NAME)

.PHONY: run-prod
## run-prod: Build and Run in production mode.
run-prod: build
	@GIN_MODE=release ENV=prod ./$(NAME)

.PHONY: clean
## clean: Clean project and previous builds.
clean:
	@rm -f $(NAME)


.PHONY: help
all: help
# help: show this help message
help: Makefile
	@echo
	@echo " Choose a command to run in "$(APP_NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
