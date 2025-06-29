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
	fetchChannelID   int64
	fetchChannelName string
	fetchLimit       int
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "抓取 Channel 历史消息",
	Long: `抓取指定 Channel 的历史消息。

可以通过 Channel ID 或用户名指定频道，推荐使用 Channel ID。
可以指定抓取的消息数量，默认抓取最新的 100 条消息。
消息会被保存到数据库中供后续分析使用。

示例:
  tgchannel fetch --id 1234567890
  tgchannel fetch --name @channel_name
  tgchannel fetch --id 1234567890 --limit 500`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := fetch(); err != nil {
			log.Fatalf("抓取失败: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)

	// 添加标志
	fetchCmd.Flags().Int64VarP(&fetchChannelID, "id", "i", 0, "Channel ID (推荐使用)")
	fetchCmd.Flags().StringVarP(&fetchChannelName, "name", "n", "", "Channel 用户名 (例如: @channel_name)")
	fetchCmd.Flags().IntVarP(&fetchLimit, "limit", "l", 100, "抓取消息数量")

	// 至少需要指定 ID 或用户名之一
	fetchCmd.MarkFlagsMutuallyExclusive("id", "name")
}

func fetch() error {
	// 检查参数
	if fetchChannelID == 0 && fetchChannelName == "" {
		return fmt.Errorf("请指定 Channel ID (--id) 或用户名 (--name)")
	}

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
		scraperClient := scraper.NewScraper(db, client.API(), &config.Scraper)

		// 抓取历史消息
		if fetchChannelID != 0 {
			log.Printf("开始抓取频道 ID %d 的历史消息，限制 %d 条...", fetchChannelID, fetchLimit)
			if err := scraperClient.FetchChannelHistoryByID(ctx, fetchChannelID, fetchLimit); err != nil {
				return fmt.Errorf("抓取历史消息失败: %w", err)
			}
			log.Printf("频道 ID %d 的历史消息抓取完成", fetchChannelID)
		} else {
			log.Printf("开始抓取频道 %s 的历史消息，限制 %d 条...", fetchChannelName, fetchLimit)
			if err := scraperClient.FetchChannelHistory(ctx, fetchChannelName, fetchLimit); err != nil {
				return fmt.Errorf("抓取历史消息失败: %w", err)
			}
			log.Printf("频道 %s 的历史消息抓取完成", fetchChannelName)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("抓取失败: %w", err)
	}

	return nil
}
