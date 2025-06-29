package cmd

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/momaek/tgchannel/internal/auth"
	"github.com/momaek/tgchannel/internal/database"
	"github.com/momaek/tgchannel/internal/models"
	"github.com/momaek/tgchannel/internal/scraper"
	"github.com/spf13/cobra"
)

var (
	channelUsername string
	userUsername    string
)

// subscribeCmd represents the subscribe command
var subscribeCmd = &cobra.Command{
	Use:   "subscribe",
	Short: "订阅 Telegram Channel",
	Long: `订阅指定的 Telegram Channel。

订阅后，系统会自动抓取该 Channel 的历史消息，
并在有新消息时自动保存到数据库。`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := subscribe(); err != nil {
			log.Fatalf("订阅失败: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(subscribeCmd)

	// 添加标志
	subscribeCmd.Flags().StringVarP(&channelUsername, "channel", "c", "", "Channel 用户名 (例如: @channel_name)")
	subscribeCmd.Flags().StringVarP(&userUsername, "user", "u", "", "用户用户名")
	subscribeCmd.MarkFlagRequired("channel")
	subscribeCmd.MarkFlagRequired("user")
}

func subscribe() error {
	// 初始化数据库
	db, err := database.NewDatabase(config.Database.Path)
	if err != nil {
		return fmt.Errorf("初始化数据库失败: %w", err)
	}
	defer db.Close()

	// 获取或创建用户
	user, err := db.GetUserByUsername(userUsername)
	if err != nil {
		// 用户不存在，创建新用户
		user = &models.User{
			Username: userUsername,
		}
		if err := db.CreateUser(user); err != nil {
			return fmt.Errorf("创建用户失败: %w", err)
		}
		log.Printf("创建新用户: %s", userUsername)
	}

	// 解析 API ID
	apiID, err := strconv.Atoi(config.Telegram.APIID)
	if err != nil {
		return fmt.Errorf("无效的 API ID: %w", err)
	}

	// 创建认证客户端
	authClient := auth.NewAuth(apiID, config.Telegram.APIHash, config.Telegram.SessionFile)
	client := authClient.GetClient()

	// 连接到 Telegram
	ctx := context.Background()
	err = client.Run(ctx, func(ctx context.Context) error {
		// 创建爬虫实例
		scraperClient := scraper.NewScraper(db, client.API())

		// 获取频道信息
		channel, err := scraperClient.FetchChannelInfo(ctx, channelUsername)
		if err != nil {
			return fmt.Errorf("获取频道信息失败: %w", err)
		}

		// 创建订阅
		subscription := &models.Subscription{
			UserID:    user.ID,
			ChannelID: channel.ID,
			IsActive:  true,
		}

		if err := db.CreateSubscription(subscription); err != nil {
			return fmt.Errorf("创建订阅失败: %w", err)
		}

		log.Printf("成功订阅频道: %s (@%s)", channel.Title, channel.Username)
		return nil
	})

	if err != nil {
		return fmt.Errorf("订阅失败: %w", err)
	}

	return nil
}
