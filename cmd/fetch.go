package cmd

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/momaek/tgchannel/internal/auth"
	"github.com/momaek/tgchannel/internal/database"
	"github.com/momaek/tgchannel/internal/scraper"
	"github.com/spf13/cobra"
)

var (
	fetchChannelUsername string
	fetchLimit           int
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "抓取 Channel 历史消息",
	Long: `抓取指定 Channel 的历史消息。

可以指定抓取的消息数量，默认抓取最新的 100 条消息。
消息会被保存到数据库中供后续分析使用。`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := fetch(); err != nil {
			log.Fatalf("抓取失败: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)

	// 添加标志
	fetchCmd.Flags().StringVarP(&fetchChannelUsername, "channel", "c", "", "Channel 用户名 (例如: @channel_name)")
	fetchCmd.Flags().IntVarP(&fetchLimit, "limit", "l", 100, "抓取消息数量")
	fetchCmd.MarkFlagRequired("channel")
}

func fetch() error {
	// 初始化数据库
	db, err := database.NewDatabase(config.Database.Path)
	if err != nil {
		return fmt.Errorf("初始化数据库失败: %w", err)
	}
	defer db.Close()

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

		// 抓取历史消息
		log.Printf("开始抓取频道 %s 的历史消息，限制 %d 条...", fetchChannelUsername, fetchLimit)

		if err := scraperClient.FetchChannelHistory(ctx, fetchChannelUsername, fetchLimit); err != nil {
			return fmt.Errorf("抓取历史消息失败: %w", err)
		}

		log.Printf("频道 %s 的历史消息抓取完成", fetchChannelUsername)
		return nil
	})

	if err != nil {
		return fmt.Errorf("抓取失败: %w", err)
	}

	return nil
}
