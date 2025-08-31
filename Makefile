# 设置变量
APP_NAME = price-query
SRC = etc/main.go
SCRIPT_APP_NAME = price-query
SCRIPT = etc/script/script.go

# 固定目标 OS 为 linux（写死）
GOOS = linux

# 架构：默认 amd64（可按需修改或写死）
GOARCH ?= amd64

# 禁用 cgo 以生成静态二进制
CGO_ENABLED = 0

# go build 命令封装（静态构建）
GOBUILD := CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags="-s -w"

# 默认目标
all: build

# 构建项目（无 cgo，静态链接，目标为 linux）
build:
	$(GOBUILD) -o $(APP_NAME) $(SRC)

# 构建 script（无 cgo，静态链接，目标为 linux）
build-history:
	$(GOBUILD) -o $(SCRIPT_APP_NAME) $(SCRIPT)

# 运行（本地调试）
run: build
	./$(APP_NAME)

# 清理构建产物
clean:
	rm -f $(APP_NAME) $(SCRIPT_APP_NAME)

.PHONY: all build build-history run clean