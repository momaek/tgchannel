package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/momaek/tgchannel/internal/models"
)

type Database struct {
	db *sql.DB
}

// NewDatabase 创建新的数据库连接
func NewDatabase(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	database := &Database{db: db}
	if err := database.InitTables(); err != nil {
		return nil, fmt.Errorf("failed to init tables: %w", err)
	}

	return database, nil
}

// InitTables 初始化数据库表
func (d *Database) InitTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE,
			phone TEXT UNIQUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS channels (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			telegram_id INTEGER UNIQUE,
			username TEXT UNIQUE,
			title TEXT,
			description TEXT,
			member_count INTEGER DEFAULT 0,
			is_active BOOLEAN DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS subscriptions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			channel_id INTEGER,
			is_active BOOLEAN DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users (id),
			FOREIGN KEY (channel_id) REFERENCES channels (id),
			UNIQUE(user_id, channel_id)
		)`,
		`CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			telegram_id INTEGER,
			channel_id INTEGER,
			sender_id INTEGER,
			sender_name TEXT,
			text TEXT,
			media_type TEXT,
			media_url TEXT,
			views INTEGER DEFAULT 0,
			forwards INTEGER DEFAULT 0,
			replies INTEGER DEFAULT 0,
			date DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (channel_id) REFERENCES channels (id),
			UNIQUE(telegram_id, channel_id)
		)`,
	}

	for _, query := range queries {
		if _, err := d.db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}

	log.Println("数据库表初始化完成")
	return nil
}

// CreateUser 创建用户
func (d *Database) CreateUser(user *models.User) error {
	query := `INSERT INTO users (username, phone) VALUES (?, ?)`
	result, err := d.db.Exec(query, user.Username, user.Phone)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	user.ID = id
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return nil
}

// GetUserByUsername 根据用户名获取用户
func (d *Database) GetUserByUsername(username string) (*models.User, error) {
	query := `SELECT id, username, phone, created_at, updated_at FROM users WHERE username = ?`
	user := &models.User{}
	err := d.db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Phone, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// CreateChannel 创建频道
func (d *Database) CreateChannel(channel *models.Channel) error {
	query := `INSERT INTO channels (telegram_id, username, title, description, member_count) 
			  VALUES (?, ?, ?, ?, ?)`
	result, err := d.db.Exec(query, channel.TelegramID, channel.Username,
		channel.Title, channel.Description, channel.MemberCount)
	if err != nil {
		return fmt.Errorf("failed to create channel: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	channel.ID = id
	channel.CreatedAt = time.Now()
	channel.UpdatedAt = time.Now()

	return nil
}

// GetChannelByUsername 根据用户名获取频道
func (d *Database) GetChannelByUsername(username string) (*models.Channel, error) {
	query := `SELECT id, telegram_id, username, title, description, member_count, 
			  is_active, created_at, updated_at FROM channels WHERE username = ?`
	channel := &models.Channel{}
	err := d.db.QueryRow(query, username).Scan(
		&channel.ID, &channel.TelegramID, &channel.Username, &channel.Title,
		&channel.Description, &channel.MemberCount, &channel.IsActive,
		&channel.CreatedAt, &channel.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get channel: %w", err)
	}
	return channel, nil
}

// GetChannelByTelegramID 根据 Telegram ID 获取频道
func (d *Database) GetChannelByTelegramID(telegramID int64) (*models.Channel, error) {
	query := `SELECT id, telegram_id, username, title, description, member_count, 
			  is_active, created_at, updated_at FROM channels WHERE telegram_id = ?`
	channel := &models.Channel{}
	err := d.db.QueryRow(query, telegramID).Scan(
		&channel.ID, &channel.TelegramID, &channel.Username, &channel.Title,
		&channel.Description, &channel.MemberCount, &channel.IsActive,
		&channel.CreatedAt, &channel.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get channel: %w", err)
	}
	return channel, nil
}

// CreateSubscription 创建订阅
func (d *Database) CreateSubscription(subscription *models.Subscription) error {
	query := `INSERT INTO subscriptions (user_id, channel_id) VALUES (?, ?)`
	result, err := d.db.Exec(query, subscription.UserID, subscription.ChannelID)
	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	subscription.ID = id
	subscription.CreatedAt = time.Now()
	subscription.UpdatedAt = time.Now()

	return nil
}

// GetUserSubscriptions 获取用户订阅的频道
func (d *Database) GetUserSubscriptions(userID int64) ([]*models.Channel, error) {
	query := `SELECT c.id, c.telegram_id, c.username, c.title, c.description, 
			  c.member_count, c.is_active, c.created_at, c.updated_at 
			  FROM channels c 
			  JOIN subscriptions s ON c.id = s.channel_id 
			  WHERE s.user_id = ? AND s.is_active = 1`

	rows, err := d.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query subscriptions: %w", err)
	}
	defer rows.Close()

	var channels []*models.Channel
	for rows.Next() {
		channel := &models.Channel{}
		err := rows.Scan(
			&channel.ID, &channel.TelegramID, &channel.Username, &channel.Title,
			&channel.Description, &channel.MemberCount, &channel.IsActive,
			&channel.CreatedAt, &channel.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan channel: %w", err)
		}
		channels = append(channels, channel)
	}

	return channels, nil
}

// CreateMessage 创建消息
func (d *Database) CreateMessage(message *models.Message) error {
	query := `INSERT INTO messages (telegram_id, channel_id, sender_id, sender_name, 
			  text, media_type, media_url, views, forwards, replies, date) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := d.db.Exec(query, message.TelegramID, message.ChannelID,
		message.SenderID, message.SenderName, message.Text, message.MediaType,
		message.MediaURL, message.Views, message.Forwards, message.Replies, message.Date)
	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	message.ID = id
	message.CreatedAt = time.Now()
	message.UpdatedAt = time.Now()

	return nil
}

// GetChannelMessages 获取频道的消息
func (d *Database) GetChannelMessages(channelID int64, limit, offset int) ([]*models.Message, error) {
	query := `SELECT id, telegram_id, channel_id, sender_id, sender_name, 
			  text, media_type, media_url, views, forwards, replies, date, 
			  created_at, updated_at 
			  FROM messages 
			  WHERE channel_id = ? 
			  ORDER BY date DESC 
			  LIMIT ? OFFSET ?`

	rows, err := d.db.Query(query, channelID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %w", err)
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		message := &models.Message{}
		err := rows.Scan(
			&message.ID, &message.TelegramID, &message.ChannelID, &message.SenderID,
			&message.SenderName, &message.Text, &message.MediaType, &message.MediaURL,
			&message.Views, &message.Forwards, &message.Replies, &message.Date,
			&message.CreatedAt, &message.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, message)
	}

	return messages, nil
}

// GetAllMessages 获取所有消息
func (d *Database) GetAllMessages(limit, offset int) ([]*models.Message, error) {
	query := `SELECT id, telegram_id, channel_id, sender_id, sender_name, text, media_type, media_url, 
			  views, forwards, replies, date, created_at, updated_at 
			  FROM messages 
			  ORDER BY date DESC 
			  LIMIT ? OFFSET ?`

	rows, err := d.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %w", err)
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		msg := &models.Message{}
		err := rows.Scan(
			&msg.ID, &msg.TelegramID, &msg.ChannelID, &msg.SenderID, &msg.SenderName,
			&msg.Text, &msg.MediaType, &msg.MediaURL, &msg.Views, &msg.Forwards,
			&msg.Replies, &msg.Date, &msg.CreatedAt, &msg.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

// GetChannelByID 根据 ID 获取频道
func (d *Database) GetChannelByID(channelID int64) (*models.Channel, error) {
	query := `SELECT id, telegram_id, username, title, description, member_count, 
			  is_active, created_at, updated_at FROM channels WHERE id = ?`
	channel := &models.Channel{}
	err := d.db.QueryRow(query, channelID).Scan(
		&channel.ID, &channel.TelegramID, &channel.Username, &channel.Title,
		&channel.Description, &channel.MemberCount, &channel.IsActive,
		&channel.CreatedAt, &channel.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get channel: %w", err)
	}
	return channel, nil
}

// Close 关闭数据库连接
func (d *Database) Close() error {
	return d.db.Close()
}
