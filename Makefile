# 设置变量
APP_NAME = price-query
SRC = etc/main.go

# 默认目标
all: build

# 构建项目
build:
	go build -o $(APP_NAME) $(SRC)

run: build
	./$(APP_NAME)