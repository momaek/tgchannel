.PHONY: build clean test run help

# 默认目标
.DEFAULT_GOAL := help

# 构建可执行文件
build:
	@echo "构建 tgchannel..."
	go build -o tgchannel main.go
	@echo "构建完成!"

# 清理构建文件
clean:
	@echo "清理构建文件..."
	rm -f tgchannel
	rm -f session.json
	rm -f tgchannel.db
	rm -f tgchannel.log
	@echo "清理完成!"

# 运行测试
test:
	@echo "运行测试..."
	go test ./...

# 安装依赖
deps:
	@echo "安装依赖..."
	go mod tidy
	@echo "依赖安装完成!"

# 运行程序
run: build
	@echo "运行 tgchannel..."
	./tgchannel

# 登录
login: build
	@echo "登录 Telegram 账号..."
	./tgchannel login

# 订阅频道
subscribe: build
	@echo "订阅频道..."
	./tgchannel subscribe --channel @example --user testuser

# 列出订阅
list:
	@echo "列出已订阅的频道..."
	@go run main.go list

# 列出 Telegram 关注的 Channel
channels:
	@echo "列出 Telegram 关注的 Channel..."
	@go run main.go channels

# 抓取消息
fetch:
	@echo "抓取 Channel 历史消息..."
	@if [ -z "$(CHANNEL_ID)" ] && [ -z "$(CHANNEL_NAME)" ]; then \
		echo "请指定 Channel ID 或用户名"; \
		echo "示例: make fetch CHANNEL_ID=1234567890"; \
		echo "示例: make fetch CHANNEL_NAME=@channel_name"; \
		exit 1; \
	fi
	@if [ -n "$(CHANNEL_ID)" ]; then \
		go run main.go fetch --id $(CHANNEL_ID) --limit $(or $(LIMIT),100); \
	else \
		go run main.go fetch --name $(CHANNEL_NAME) --limit $(or $(LIMIT),100); \
	fi

# 启动服务
serve: build
	@echo "启动监听服务..."
	./tgchannel serve

# 查看抓取的消息
messages:
	@echo "查看抓取的消息..."
	@if [ -n "$(LIMIT)" ]; then \
		go run main.go messages --limit $(LIMIT); \
	else \
		go run main.go messages; \
	fi

# 显示帮助信息
help:
	@echo "tgchannel Makefile 命令:"
	@echo ""
	@echo "  build     - 构建可执行文件"
	@echo "  clean     - 清理构建文件"
	@echo "  test      - 运行测试"
	@echo "  deps      - 安装依赖"
	@echo "  run       - 构建并运行程序"
	@echo "  login     - 登录 Telegram 账号"
	@echo "  subscribe - 订阅频道"
	@echo "  list      - 列出订阅的 Channel"
	@echo "  fetch     - 抓取频道消息"
	@echo "  serve     - 启动监听服务"
	@echo "  messages  - 查看抓取的消息"
	@echo "  help      - 显示此帮助信息"
	@echo ""
	@echo "使用示例:"
	@echo "  make build     # 构建程序"
	@echo "  make login     # 登录账号"
	@echo "  make subscribe # 订阅频道"
	@echo "  make list      # 列出订阅"
	@echo "  make fetch     # 抓取消息"
	@echo "  make serve     # 启动服务"
	@echo "  make messages  # 查看抓取的消息" 