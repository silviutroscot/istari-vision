# add all your cmd/<things> in here
TARGETS = executable
NAME := auth-service
TIMESTAMP := $$(date +%s)
IMG    := ${NAME}:${TIMESTAMP}
LATEST := ${NAME}:latest

# install linter
golangci-lint = ./bin/golangci-lint
$(golangci-lint):
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.27.0

# install gosec
gosec = ./bin/gosec
$(gosec):
	curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s v2.3.0

CMD_DIR := ./cmd
PKG_DIR := ./pkg
OUT_DIR := ./out

COV_FILE := cover.out

GO111MODULE := on

GO_TEST_FLAGS := -v -count=1 -coverprofile=$(OUT_DIR)/$(COV_FILE) -covermode=atomic

.PHONY: $(OUT_DIR) clean build test mod cover test-deps fmt vet purge bench lint sec docker-build

all: clean mod fmt vet test build lint sec

$(OUT_DIR):
	@mkdir -p $(OUT_DIR)

clean:
	@rm -rf $(OUT_DIR)

purge: clean
	go mod tidy
	go clean -cache
	go clean -testcache
	go clean -modcache

build: $(OUT_DIR)
	$(foreach target,$(TARGETS),go build -o $(OUT_DIR)/$(target) $(CMD_DIR)/$(target)/*.go;)

test: $(OUT_DIR)
	go test $(GO_TEST_FLAGS) ./...

mod:
	go mod tidy
	go mod verify

cover:
	go tool cover -html=$(OUT_DIR)/$(COV_FILE) -o ./frontend/cover.html

test-deps:
	go test all

fmt:
	go fmt ./...

vet:
	go vet ./...

bench:
	go test -bench=. -benchmem -benchtime=10s ./...

lint:
	golangci-lint run

sec: ## Security scan
	gosec -exclude G601,G404 ./...

docker-build:
	@docker build -t ${IMG} -t ${LATEST} -f prod-dockerfile .

# ROOT NODE
export NODE_PATH := $(shell npm root -g)

build-frontend:
	@mkdir -p ./frontend/assets/css
	tailwind -c tailwind.config.cjs -m --no-autoprefixer -o ./frontend/assets/css/main.min.css


