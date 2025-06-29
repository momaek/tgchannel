# tgchannel

Telegram Channel 爬虫系统

## 功能特性

- 🔐 **用户认证**: 支持手机号登录和验证码认证
- 📺 **频道订阅**: 订阅指定的 Telegram 频道
- 📊 **数据抓取**: 抓取频道历史消息并存储到数据库
- 🔄 **实时监听**: 监听频道更新并推送新消息
- 📋 **订阅管理**: 查看和管理已订阅的频道
- 🌐 **API 集成**: 通过 Telegram API 获取关注列表
- 💾 **数据存储**: 支持 SQLite 数据库存储
- 🎯 **消息处理**: 支持消息的进一步处理和分析

## 快速开始

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 配置环境

创建配置文件 `config.yaml`:

```yaml
telegram:
  api_id: "你的API_ID"
  api_hash: "你的API_HASH"
  session_file: "session.json"

database:
  type: "sqlite"
  path: "tgchannel.db"

logging:
  level: "info"
  file: "tgchannel.log"
```

### 3. 运行程序

```bash
# 登录 Telegram 账号
go run main.go login

# 订阅 Channel
go run main.go subscribe --channel @channel_name --user your_username

# 列出订阅的 Channel
go run main.go list --user your_username

# 抓取历史数据
go run main.go fetch --channel @channel_name

# 启动监听服务
go run main.go serve
```

## 项目结构

```
tgchannel/
├── cmd/                    # 命令行工具
├── internal/              # 内部包
│   ├── auth/             # 认证相关
│   ├── client/           # Telegram 客户端
│   ├── database/         # 数据库操作
│   ├── models/           # 数据模型
│   └── scraper/          # 爬虫逻辑
├── configs/              # 配置文件
├── main.go               # 主程序入口
└── README.md
```

## 许可证

MIT License

## 使用示例

```bash
# 登录认证
make login

# 订阅频道
make subscribe CHANNEL=@example_channel

# 查看已订阅的频道（从数据库）
make list

# 查看 Telegram 关注的 Channel（通过 API）
make channels

# 抓取历史消息
make fetch CHANNEL=@example_channel

# 启动监听服务
make serve
```
