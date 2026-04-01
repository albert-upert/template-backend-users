.PHONY: build run test docker-build clean

APP_NAME := users
DOCKER_IMAGE := $(APP_NAME):latest
# Assumes build is run from repository root if using docker-build
DOCKER_CONTEXT := ../

# Build Variables
VERSION ?= $(shell git describe --tags --always --dirty || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD || echo "none")
BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD || echo "unknown")
BUILD_TIME ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS := -ldflags "-X 'yuhuu.universitaspertamina.ac.id/siak/siakup/backend/user/version.Version=$(VERSION)' \
                     -X 'yuhuu.universitaspertamina.ac.id/siak/siakup/backend/user/version.Commit=$(COMMIT)' \
                     -X 'yuhuu.universitaspertamina.ac.id/siak/siakup/backend/user/version.Branch=$(BRANCH)' \
                     -X 'yuhuu.universitaspertamina.ac.id/siak/siakup/backend/user/version.BuildTime=$(BUILD_TIME)'"

build:
	go build $(LDFLAGS) -o bin/$(APP_NAME) main.go

run:
	go run $(LDFLAGS) main.go serve

test:
	go test -v -cover ./...

clean:
	rm -rf bin/

# Docker build command to be run from the users directory, but setting context to root
docker-build:
	docker build \
		--build-arg VERSION="$(VERSION)" \
		--build-arg COMMIT="$(COMMIT)" \
		--build-arg BRANCH="$(BRANCH)" \
		--build-arg BUILD_TIME="$(BUILD_TIME)" \
		-f Dockerfile -t $(DOCKER_IMAGE) $(DOCKER_CONTEXT)

# Helper to run docker build from repo root manually:
# docker build -f users/Dockerfile -t users:latest .
