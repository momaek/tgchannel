# Telegram Channel Scraper

A powerful Telegram Channel scraper and monitoring tool built with Go. This tool allows you to subscribe to Telegram channels, fetch historical messages, monitor updates, and store data in SQLite database.

## ✨ Features

- 🔐 **User Authentication**: Support phone number login and verification code authentication
- 📺 **Channel Subscription**: Subscribe to specified Telegram channels
- 📊 **Data Scraping**: Fetch channel historical messages and store in database
- 🔄 **Real-time Monitoring**: Monitor channel updates and push new messages
- 📋 **Subscription Management**: View and manage subscribed channels
- 🌐 **API Integration**: Get channel list through Telegram API
- 💾 **Data Storage**: Support SQLite database storage
- 🎯 **Message Processing**: Support further message processing and analysis
- 📱 **Channel Discovery**: List all channels you follow on Telegram

## 🚀 Quick Start

### Prerequisites

- Go 1.19 or higher
- SQLite3
- Telegram API credentials (API ID and API Hash)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/momaek/tgchannel.git
   cd tgchannel
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Configure the application**
   ```bash
   cp configs/config.example.yaml configs/config.yaml
   ```
   
   Edit `configs/config.yaml` and add your Telegram API credentials:
   ```yaml
   telegram:
     api_id: "your_api_id"
     api_hash: "your_api_hash"
     session_file: "session.json"
   ```

### Usage

#### 1. Authentication
```bash
# First-time login, requires phone number and verification code
make login
```

#### 2. Subscribe to Channels
```bash
# Subscribe to a specific channel
make subscribe CHANNEL=@channel_username
```

#### 3. View Subscribed Channels
```bash
# View subscribed channels (from database)
make list

# View all channels you follow on Telegram (via API)
make channels
```

#### 4. Fetch Historical Messages
```bash
# Fetch using Channel ID (recommended)
make fetch CHANNEL_ID=1234567890

# Fetch using username
make fetch CHANNEL_NAME=@channel_name

# Specify fetch limit
make fetch CHANNEL_ID=1234567890 LIMIT=500
```

#### 5. View Fetched Messages
```bash
# View all messages (latest 10)
make messages

# View specific number of messages
make messages LIMIT=20

# View messages from specific channel
make messages --id 1234567890
```

#### 6. Start Monitoring Service
```bash
# Start service to monitor channel updates
make serve
```

## 📖 Command Reference

### Authentication
```bash
go run main.go login
```
Authenticate with Telegram using phone number and verification code.

### Channel Management
```bash
# List subscribed channels (from database)
go run main.go list

# List all channels you follow (via Telegram API)
go run main.go channels
```

### Message Operations
```bash
# Subscribe to channel
go run main.go subscribe --channel @channel_name

# Fetch historical messages
go run main.go fetch --id 1234567890 --limit 100
go run main.go fetch --name @channel_name --limit 100

# View messages
go run main.go messages --limit 10
go run main.go messages --id 1234567890
```

### Service
```bash
# Start monitoring service
go run main.go serve
```

## ⚙️ Configuration

### Telegram Configuration
```yaml
telegram:
  api_id: "your_api_id"
  api_hash: "your_api_hash"
  session_file: "session.json"
  device_model: "Desktop"
  system_version: "Windows 10"
  app_version: "1.0.0"
```

### Database Configuration
```yaml
database:
  type: "sqlite"
  path: "tgchannel.db"
```

### Scraper Configuration
```yaml
scraper:
  # Number of messages per request (recommended: 50-200)
  batch_size: 100
  
  # Delay between requests in seconds (recommended: 1-5)
  delay_between_requests: 2
  
  # Maximum retry attempts
  max_retries: 3
```

## 🔧 Advanced Features

### Paginated Fetching
The tool automatically implements paginated fetching to avoid Telegram's rate limiting:
- Fetches messages in batches (configurable batch size)
- Adds delays between requests (configurable delay)
- Supports fetching large amounts of historical messages (e.g., 10,000 messages)

### Message Processing
- Extracts message text, media information, and metadata
- Stores sender information, views, forwards, and replies
- Supports various media types (photos, documents, webpages)

### Database Schema
- **Users**: User authentication information
- **Channels**: Channel metadata and statistics
- **Subscriptions**: User-channel subscription relationships
- **Messages**: Complete message data with metadata

## 📊 Data Analysis

The stored data can be used for:
- Content analysis and trending topics
- Engagement metrics analysis
- Historical data mining
- Automated content monitoring

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## ⚠️ Disclaimer

This tool is for educational and research purposes only. Please respect Telegram's Terms of Service and use responsibly. The authors are not responsible for any misuse of this tool.

## 🆘 Troubleshooting

### Common Issues

1. **"CHANNEL_INVALID" error**
   - Make sure you're a member of the channel
   - Check if the channel ID is correct

2. **Rate limiting issues**
   - Increase `delay_between_requests` in configuration
   - Reduce `batch_size` in configuration

3. **Authentication problems**
   - Delete `session.json` and re-authenticate
   - Check your API credentials

### Getting Help

- Check the [documentation](docs/GETTING_STARTED.md)
- Open an issue on GitHub
- Review the configuration examples

## 🔗 Links

- [Telegram API Documentation](https://core.telegram.org/api)
- [gotd/td Library](https://github.com/gotd/td)
- [Go Documentation](https://golang.org/doc/)

---

**Made with ❤️ for the Telegram community**
