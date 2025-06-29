# 快速开始指南

## 1. 获取 Telegram API 凭据

在使用本工具之前，您需要获取 Telegram API 凭据：

1. 访问 [my.telegram.org](https://my.telegram.org)
2. 使用您的 Telegram 账号登录
3. 点击 "API development tools"
4. 创建一个新的应用程序
5. 记录下 `api_id` 和 `api_hash`

## 2. 配置应用程序

1. 复制示例配置文件：
   ```bash
   cp configs/config.example.yaml configs/config.yaml
   ```

2. 编辑配置文件，填入您的 API 凭据：
   ```yaml
   telegram:
     api_id: "你的API_ID"
     api_hash: "你的API_HASH"
     session_file: "session.json"
   ```

## 3. 安装依赖

```bash
go mod tidy
```

## 4. 登录 Telegram 账号

首次使用需要登录您的 Telegram 账号：

```bash
go run main.go login
```

系统会提示您输入：
- 手机号码（格式：+86xxxxxxxxxx）
- 验证码（从 Telegram 收到的短信）
- 两步验证密码（如果启用了两步验证）

登录成功后，会话信息会保存在 `session.json` 文件中。

## 5. 订阅 Channel

订阅您想要监控的 Channel：

```bash
go run main.go subscribe --channel @channel_name --user your_username
```

例如：
```bash
go run main.go subscribe --channel @telegram --user myuser
```

## 6. 抓取历史消息

抓取指定 Channel 的历史消息：

```bash
go run main.go fetch --channel @channel_name --limit 100
```

参数说明：
- `--channel`: Channel 用户名（必需）
- `--limit`: 抓取消息数量（可选，默认 100）

### 5. 抓取历史消息
```bash
# 使用 Channel ID 抓取（推荐）
go run main.go fetch --channel @channel_name --limit 500 --channel_id 1234567890

# 使用用户名抓取
go run main.go fetch --channel @channel_name --limit 500 --channel_name @channel_username

# 指定抓取数量
go run main.go fetch --channel @channel_name --limit 500 --channel_id 1234567890
```

**分页抓取说明**：
- 系统会自动进行分页抓取，避免被 Telegram 限流
- 每次请求 100 条消息（可在配置文件中调整）
- 请求间隔 2 秒（可在配置文件中调整）
- 支持抓取大量历史消息（如 10000 条）

## 7. 列出订阅的 Channel

查看指定用户订阅的所有 Channel：

```bash
go run main.go list --user your_username
```

例如：
```bash
go run main.go list --user myuser
```

输出示例：
```
用户 myuser 订阅的 Channel 列表:
================================================================================
用户名                标题                          成员数量        状态
--------------------------------------------------------------------------------
@telegram            Telegram                      12345678        活跃
@example_channel     示例频道                       9876543         活跃
================================================================================
总计: 2 个 Channel
```

## 8. 启动监听服务

启动服务来监听 Channel 更新：

```bash
go run main.go serve
```

服务会持续运行，定期检查订阅的 Channel 是否有新消息，并自动保存到数据库。

按 `Ctrl+C` 可以优雅地停止服务。

## 9. 查看帮助

查看所有可用命令：

```bash
go run main.go --help
```

查看特定命令的帮助：

```bash
go run main.go login --help
go run main.go subscribe --help
go run main.go fetch --help
go run main.go list --help
go run main.go serve --help
```

## 数据库结构

系统使用 SQLite 数据库存储数据，包含以下表：

- `users`: 用户信息
- `channels`: Channel 信息
- `subscriptions`: 订阅关系
- `messages`: 消息内容

## 注意事项

1. **API 限制**: Telegram API 有速率限制，请合理使用
2. **隐私**: 请遵守相关法律法规和隐私政策
3. **数据安全**: 妥善保管您的 API 凭据和会话文件
4. **备份**: 定期备份数据库文件

## 故障排除

### 登录失败
- 检查 API 凭据是否正确
- 确保手机号码格式正确
- 检查网络连接

### 抓取失败
- 确认 Channel 用户名正确
- 检查是否有权限访问该 Channel
- 确认已经成功登录

### 服务启动失败
- 检查配置文件是否正确
- 确认数据库文件路径可写
- 检查网络连接

## 使用命令

### 1. 登录认证
```bash
# 首次登录，需要输入手机号和验证码
make login
```

### 2. 订阅频道
```bash
# 订阅指定频道
make subscribe CHANNEL=@channel_username
```

### 3. 查看订阅列表
```bash
# 查看已订阅的频道（从数据库）
make list
```

### 4. 查看 Telegram 关注的 Channel
```bash
# 查看当前账号在 Telegram 中关注的所有 Channel（通过 API）
make channels
```

### 5. 抓取历史消息
```bash
# 使用 Channel ID 抓取（推荐）
make fetch CHANNEL_ID=1234567890

# 使用用户名抓取
make fetch CHANNEL_NAME=@channel_username

# 指定抓取数量
make fetch CHANNEL_ID=1234567890 LIMIT=500
```

### 6. 启动监听服务
```bash
# 启动服务监听频道更新
make serve
```