# Telegram Channel 爬虫系统

一个功能强大的 Telegram Channel 爬虫和监控工具，使用 Go 语言开发。该工具允许你订阅 Telegram 频道、抓取历史消息、监控更新并将数据存储到 SQLite 数据库中。

## ✨ 功能特性

- 🔐 **用户认证**: 支持手机号登录和验证码认证
- 📺 **频道订阅**: 订阅指定的 Telegram 频道
- 📊 **数据抓取**: 抓取频道历史消息并存储到数据库
- 🔄 **实时监听**: 监听频道更新并推送新消息
- 📋 **订阅管理**: 查看和管理已订阅的频道
- 🌐 **API 集成**: 通过 Telegram API 获取关注列表
- 💾 **数据存储**: 支持 SQLite 数据库存储
- 🎯 **消息处理**: 支持消息的进一步处理和分析
- 📱 **频道发现**: 列出你在 Telegram 上关注的所有频道

## 🚀 快速开始

### 环境要求

- Go 1.19 或更高版本
- SQLite3
- Telegram API 凭据（API ID 和 API Hash）

### 安装

1. **克隆仓库**
   ```bash
   git clone https://github.com/momaek/tgchannel.git
   cd tgchannel
   ```

2. **安装依赖**
   ```bash
   go mod download
   ```

3. **配置应用**
   ```bash
   cp configs/config.example.yaml configs/config.yaml
   ```
   
   编辑 `configs/config.yaml` 并添加你的 Telegram API 凭据：
   ```yaml
   telegram:
     api_id: "你的API_ID"
     api_hash: "你的API_HASH"
     session_file: "session.json"
   ```

### 使用方法

#### 1. 认证登录
```bash
# 首次登录，需要输入手机号和验证码
make login
```

#### 2. 订阅频道
```bash
# 订阅指定频道
make subscribe CHANNEL=@channel_username
```

#### 3. 查看订阅频道
```bash
# 查看已订阅的频道（从数据库）
make list

# 查看 Telegram 关注的 Channel（通过 API）
make channels
```

#### 4. 抓取历史消息
```bash
# 使用 Channel ID 抓取（推荐）
make fetch CHANNEL_ID=1234567890

# 使用用户名抓取
make fetch CHANNEL_NAME=@channel_name

# 指定抓取数量
make fetch CHANNEL_ID=1234567890 LIMIT=500
```

#### 5. 查看抓取的消息
```bash
# 查看所有消息（最新10条）
make messages

# 查看指定数量
make messages LIMIT=20

# 查看特定频道的消息
make messages --id 1234567890
```

#### 6. 启动监听服务
```bash
# 启动服务监听频道更新
make serve
```

## 📖 命令参考

### 认证
```bash
go run main.go login
```
使用手机号和验证码进行 Telegram 认证。

### 频道管理
```bash
# 列出订阅的频道（从数据库）
go run main.go list

# 列出关注的所有频道（通过 Telegram API）
go run main.go channels
```

### 消息操作
```bash
# 订阅频道
go run main.go subscribe --channel @channel_name

# 抓取历史消息
go run main.go fetch --id 1234567890 --limit 100
go run main.go fetch --name @channel_name --limit 100

# 查看消息
go run main.go messages --limit 10
go run main.go messages --id 1234567890
```

### 服务
```bash
# 启动监听服务
go run main.go serve
```

## ⚙️ 配置

### Telegram 配置
```yaml
telegram:
  api_id: "你的API_ID"
  api_hash: "你的API_HASH"
  session_file: "session.json"
  device_model: "Desktop"
  system_version: "Windows 10"
  app_version: "1.0.0"
```

### 数据库配置
```yaml
database:
  type: "sqlite"
  path: "tgchannel.db"
```

### 爬虫配置
```yaml
scraper:
  # 每次请求的消息数量（建议：50-200）
  batch_size: 100
  
  # 请求间隔时间（秒）（建议：1-5）
  delay_between_requests: 2
  
  # 最大重试次数
  max_retries: 3
```

## 🔧 高级功能

### 分页抓取
工具自动实现分页抓取以避免 Telegram 的限流：
- 分批抓取消息（可配置批次大小）
- 请求之间添加延迟（可配置延迟时间）
- 支持抓取大量历史消息（如 10,000 条）

### 消息处理
- 提取消息文本、媒体信息和元数据
- 存储发送者信息、浏览数、转发数和回复数
- 支持各种媒体类型（照片、文档、网页）

### 数据库架构
- **Users**: 用户认证信息
- **Channels**: 频道元数据和统计信息
- **Subscriptions**: 用户-频道订阅关系
- **Messages**: 完整的消息数据和元数据

## 📊 数据分析

存储的数据可用于：
- 内容分析和热门话题
- 参与度指标分析
- 历史数据挖掘
- 自动化内容监控

## 🤝 贡献

1. Fork 仓库
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 打开 Pull Request

## 📝 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## ⚠️ 免责声明

此工具仅供教育和研究目的使用。请尊重 Telegram 的服务条款并负责任地使用。作者不对此工具的误用负责。

## 🆘 故障排除

### 常见问题

1. **"CHANNEL_INVALID" 错误**
   - 确保你是该频道的成员
   - 检查频道 ID 是否正确

2. **限流问题**
   - 在配置中增加 `delay_between_requests`
   - 在配置中减少 `batch_size`

3. **认证问题**
   - 删除 `session.json` 并重新认证
   - 检查你的 API 凭据

### 获取帮助

- 查看[文档](docs/GETTING_STARTED.md)
- 在 GitHub 上提交 issue
- 查看配置示例

## 🔗 链接

- [Telegram API 文档](https://core.telegram.org/api)
- [gotd/td 库](https://github.com/gotd/td)
- [Go 文档](https://golang.org/doc/)

---

**为 Telegram 社区而制作 ❤️** 