package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/momaek/tgchannel/internal/auth"
	"github.com/momaek/tgchannel/internal/database"
	"github.com/momaek/tgchannel/internal/scraper"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "启动监听服务",
	Long: `启动监听服务，持续监控订阅的 Channel 更新。

服务会定期检查订阅的 Channel 是否有新消息，
并自动保存到数据库中。`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := serve(); err != nil {
			log.Fatalf("服务启动失败: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func serve() error {
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

	// 创建上下文，支持优雅关闭
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 处理信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("收到关闭信号，正在优雅关闭...")
		cancel()
	}()

	// 连接到 Telegram
	err = client.Run(ctx, func(ctx context.Context) error {
		// 创建爬虫实例
		scraperClient := scraper.NewScraper(db, client.API(), &config.Scraper)

		log.Println("启动监听服务...")
		log.Println("按 Ctrl+C 停止服务")

		// 开始监听更新
		return scraperClient.ListenForUpdates(ctx)
	})

	if err != nil {
		return fmt.Errorf("服务运行失败: %w", err)
	}

	log.Println("服务已停止")
	return nil
}
