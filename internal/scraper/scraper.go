package scraper

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gotd/td/tg"
	"github.com/momaek/tgchannel/internal/database"
	"github.com/momaek/tgchannel/internal/models"
)

type Scraper struct {
	db     *database.Database
	client *tg.Client
}

// NewScraper 创建新的爬虫实例
func NewScraper(db *database.Database, client *tg.Client) *Scraper {
	return &Scraper{
		db:     db,
		client: client,
	}
}

// FetchChannelInfo 获取频道信息
func (s *Scraper) FetchChannelInfo(ctx context.Context, username string) (*models.Channel, error) {
	// 移除 @ 符号
	username = strings.TrimPrefix(username, "@")

	// 解析频道
	peer, err := s.client.ContactsResolveUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("解析频道失败: %w", err)
	}

	// 检查是否是频道
	channel, ok := peer.Chats[0].(*tg.Channel)
	if !ok {
		return nil, fmt.Errorf("不是有效的频道")
	}

	// 创建频道模型
	channelModel := &models.Channel{
		TelegramID:  channel.ID,
		Username:    username,
		Title:       channel.Title,
		Description: "", // 暂时留空
		MemberCount: int32(channel.ParticipantsCount),
		IsActive:    true,
	}

	// 保存到数据库
	if err := s.db.CreateChannel(channelModel); err != nil {
		return nil, fmt.Errorf("保存频道失败: %w", err)
	}

	log.Printf("频道信息获取成功: %s (%s)", channelModel.Title, channelModel.Username)
	return channelModel, nil
}

// FetchChannelHistory 抓取频道历史消息
func (s *Scraper) FetchChannelHistory(ctx context.Context, channelUsername string, limit int) error {
	// 获取或创建频道
	channel, err := s.db.GetChannelByUsername(channelUsername)
	if err != nil {
		// 如果频道不存在，先获取频道信息
		channel, err = s.FetchChannelInfo(ctx, channelUsername)
		if err != nil {
			return fmt.Errorf("获取频道信息失败: %w", err)
		}
	}

	// 移除 @ 符号
	username := strings.TrimPrefix(channelUsername, "@")

	// 解析频道
	peer, err := s.client.ContactsResolveUsername(ctx, username)
	if err != nil {
		return fmt.Errorf("解析频道失败: %w", err)
	}

	// 转换为 InputPeer
	peerChannel := peer.Peer.(*tg.PeerChannel)
	inputPeer := &tg.InputPeerChannel{
		ChannelID:  peerChannel.ChannelID,
		AccessHash: 0, // 暂时设为0，需要从其他地方获取
	}

	// 获取历史消息
	history, err := s.client.MessagesGetHistory(ctx, &tg.MessagesGetHistoryRequest{
		Peer:       inputPeer,
		OffsetID:   0,
		OffsetDate: 0,
		AddOffset:  0,
		Limit:      limit,
		MaxID:      0,
		MinID:      0,
		Hash:       0,
	})
	if err != nil {
		return fmt.Errorf("获取历史消息失败: %w", err)
	}

	// 处理消息
	messages, ok := history.(*tg.MessagesMessages)
	if !ok {
		return fmt.Errorf("无效的消息响应")
	}

	log.Printf("开始处理 %d 条消息...", len(messages.Messages))

	for _, msg := range messages.Messages {
		if err := s.processMessage(ctx, msg, channel.ID); err != nil {
			log.Printf("处理消息失败: %v", err)
			continue
		}
	}

	log.Printf("历史消息抓取完成，共处理 %d 条消息", len(messages.Messages))
	return nil
}

// processMessage 处理单条消息
func (s *Scraper) processMessage(ctx context.Context, msg tg.MessageClass, channelID int64) error {
	message, ok := msg.(*tg.Message)
	if !ok {
		return fmt.Errorf("无效的消息类型")
	}

	// 创建消息模型
	messageModel := &models.Message{
		TelegramID: int64(message.ID),
		ChannelID:  channelID,
		SenderID:   0,
		SenderName: "",
		Text:       message.Message,
		Views:      int32(message.Views),
		Forwards:   int32(message.Forwards),
		Replies:    0,
		Date:       time.Unix(int64(message.Date), 0),
	}

	// 处理发送者信息
	if message.FromID != nil {
		if peerUser, ok := message.FromID.(*tg.PeerUser); ok {
			messageModel.SenderID = peerUser.UserID
		}
	}

	// 处理媒体文件
	if message.Media != nil {
		if err := s.processMedia(ctx, message.Media, messageModel); err != nil {
			log.Printf("处理媒体文件失败: %v", err)
		}
	}

	// 保存到数据库
	if err := s.db.CreateMessage(messageModel); err != nil {
		return fmt.Errorf("保存消息失败: %w", err)
	}

	return nil
}

// processMedia 处理媒体文件
func (s *Scraper) processMedia(ctx context.Context, media tg.MessageMediaClass, message *models.Message) error {
	switch m := media.(type) {
	case *tg.MessageMediaPhoto:
		message.MediaType = "photo"
	case *tg.MessageMediaDocument:
		message.MediaType = "document"
	case *tg.MessageMediaWebPage:
		message.MediaType = "webpage"
		if webpage, ok := m.Webpage.(*tg.WebPage); ok {
			message.MediaURL = webpage.URL
		}
	default:
		message.MediaType = "unknown"
	}
	return nil
}

// ListenForUpdates 监听频道更新
func (s *Scraper) ListenForUpdates(ctx context.Context) error {
	log.Println("开始监听频道更新...")

	// 获取所有订阅的频道
	// 这里需要实现获取用户订阅频道的逻辑
	// 暂时使用示例数据
	subscribedChannels := []string{"@example_channel"}

	for _, channelUsername := range subscribedChannels {
		go s.monitorChannel(ctx, channelUsername)
	}

	// 保持运行
	<-ctx.Done()
	return nil
}

// monitorChannel 监控单个频道
func (s *Scraper) monitorChannel(ctx context.Context, channelUsername string) {
	log.Printf("开始监控频道: %s", channelUsername)

	// 获取频道信息
	channel, err := s.db.GetChannelByUsername(channelUsername)
	if err != nil {
		log.Printf("获取频道信息失败: %v", err)
		return
	}

	// 获取最新消息ID
	messages, err := s.db.GetChannelMessages(channel.ID, 1, 0)
	if err != nil {
		log.Printf("获取最新消息失败: %v", err)
		return
	}

	var lastMessageID int64
	if len(messages) > 0 {
		lastMessageID = messages[0].TelegramID
	}

	// 定期检查新消息
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := s.checkNewMessages(ctx, channel, lastMessageID); err != nil {
				log.Printf("检查新消息失败: %v", err)
			}
		}
	}
}

// checkNewMessages 检查新消息
func (s *Scraper) checkNewMessages(ctx context.Context, channel *models.Channel, lastMessageID int64) error {
	// 移除 @ 符号
	username := strings.TrimPrefix(channel.Username, "@")

	// 解析频道
	peer, err := s.client.ContactsResolveUsername(ctx, username)
	if err != nil {
		return fmt.Errorf("解析频道失败: %w", err)
	}

	// 转换为 InputPeer
	peerChannel := peer.Peer.(*tg.PeerChannel)
	inputPeer := &tg.InputPeerChannel{
		ChannelID:  peerChannel.ChannelID,
		AccessHash: 0, // 暂时设为0
	}

	// 获取最新消息
	history, err := s.client.MessagesGetHistory(ctx, &tg.MessagesGetHistoryRequest{
		Peer:       inputPeer,
		OffsetID:   0,
		OffsetDate: 0,
		AddOffset:  0,
		Limit:      10,
		MaxID:      0,
		MinID:      0,
		Hash:       0,
	})
	if err != nil {
		return fmt.Errorf("获取最新消息失败: %w", err)
	}

	messages, ok := history.(*tg.MessagesMessages)
	if !ok {
		return fmt.Errorf("无效的消息响应")
	}

	// 处理新消息
	for _, msg := range messages.Messages {
		message, ok := msg.(*tg.Message)
		if !ok {
			continue
		}

		// 检查是否是新消息
		if int64(message.ID) > lastMessageID {
			if err := s.processMessage(ctx, msg, channel.ID); err != nil {
				log.Printf("处理新消息失败: %v", err)
				continue
			}
			lastMessageID = int64(message.ID)
			log.Printf("收到新消息: %s", message.Message[:min(len(message.Message), 50)])
		}
	}

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
