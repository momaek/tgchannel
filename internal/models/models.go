package models

import (
	"time"
)

// User 用户模型
type User struct {
	ID        int64     `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Phone     string    `json:"phone" db:"phone"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Channel 频道模型
type Channel struct {
	ID          int64     `json:"id" db:"id"`
	TelegramID  int64     `json:"telegram_id" db:"telegram_id"`
	Username    string    `json:"username" db:"username"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	MemberCount int32     `json:"member_count" db:"member_count"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Subscription 订阅模型
type Subscription struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	ChannelID int64     `json:"channel_id" db:"channel_id"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Message 消息模型
type Message struct {
	ID         int64     `json:"id" db:"id"`
	TelegramID int64     `json:"telegram_id" db:"telegram_id"`
	ChannelID  int64     `json:"channel_id" db:"channel_id"`
	SenderID   int64     `json:"sender_id" db:"sender_id"`
	SenderName string    `json:"sender_name" db:"sender_name"`
	Text       string    `json:"text" db:"text"`
	MediaType  string    `json:"media_type" db:"media_type"`
	MediaURL   string    `json:"media_url" db:"media_url"`
	Views      int32     `json:"views" db:"views"`
	Forwards   int32     `json:"forwards" db:"forwards"`
	Replies    int32     `json:"replies" db:"replies"`
	Date       time.Time `json:"date" db:"date"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// Config 配置模型
type Config struct {
	Telegram TelegramConfig `mapstructure:"telegram"`
	Database DatabaseConfig `mapstructure:"database"`
	Logging  LoggingConfig  `mapstructure:"logging"`
	Scraper  ScraperConfig  `mapstructure:"scraper"`
}

type TelegramConfig struct {
	APIID         string `mapstructure:"api_id"`
	APIHash       string `mapstructure:"api_hash"`
	SessionFile   string `mapstructure:"session_file"`
	DeviceModel   string `mapstructure:"device_model"`
	SystemVersion string `mapstructure:"system_version"`
	AppVersion    string `mapstructure:"app_version"`
}

type DatabaseConfig struct {
	Type string `mapstructure:"type"`
	Path string `mapstructure:"path"`
}

type LoggingConfig struct {
	Level string `mapstructure:"level"`
	File  string `mapstructure:"file"`
}

type ScraperConfig struct {
	BatchSize            int `mapstructure:"batch_size"`
	DelayBetweenRequests int `mapstructure:"delay_between_requests"`
	MaxRetries           int `mapstructure:"max_retries"`
}
