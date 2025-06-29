package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/momaek/tgchannel/internal/database"
	"github.com/momaek/tgchannel/internal/models"
	"github.com/spf13/cobra"
)

var (
	messagesChannelID   int64
	messagesChannelName string
	messagesLimit       int
	messagesOffset      int
)

// messagesCmd represents the messages command
var messagesCmd = &cobra.Command{
	Use:   "messages",
	Short: "查看抓取的消息",
	Long: `查看数据库中抓取的消息。

可以通过 Channel ID 或用户名筛选特定频道的消息。
支持分页查看，默认显示最新的 10 条消息。

示例:
  tgchannel messages --id 1234567890
  tgchannel messages --name @channel_name
  tgchannel messages --id 1234567890 --limit 20 --offset 10`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := listMessages(); err != nil {
			log.Fatalf("查看消息失败: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(messagesCmd)

	// 添加标志
	messagesCmd.Flags().Int64VarP(&messagesChannelID, "id", "i", 0, "Channel ID")
	messagesCmd.Flags().StringVarP(&messagesChannelName, "name", "n", "", "Channel 用户名")
	messagesCmd.Flags().IntVarP(&messagesLimit, "limit", "l", 10, "显示消息数量")
	messagesCmd.Flags().IntVarP(&messagesOffset, "offset", "o", 0, "偏移量")
}

func listMessages() error {
	// 初始化数据库
	db, err := database.NewDatabase(config.Database.Path)
	if err != nil {
		return fmt.Errorf("初始化数据库失败: %w", err)
	}
	defer db.Close()

	// 获取消息
	var messages []*models.Message
	var channel *models.Channel

	if messagesChannelID != 0 {
		// 通过 Channel ID 获取消息
		channel, err = db.GetChannelByTelegramID(messagesChannelID)
		if err != nil {
			return fmt.Errorf("获取频道失败: %w", err)
		}
		messages, err = db.GetChannelMessages(channel.ID, messagesLimit, messagesOffset)
	} else if messagesChannelName != "" {
		// 通过用户名获取消息
		channel, err = db.GetChannelByUsername(messagesChannelName)
		if err != nil {
			return fmt.Errorf("获取频道失败: %w", err)
		}
		messages, err = db.GetChannelMessages(channel.ID, messagesLimit, messagesOffset)
	} else {
		// 获取所有消息
		messages, err = db.GetAllMessages(messagesLimit, messagesOffset)
	}

	if err != nil {
		return fmt.Errorf("获取消息失败: %w", err)
	}

	if len(messages) == 0 {
		fmt.Println("没有找到消息")
		return nil
	}

	// 显示消息
	fmt.Printf("找到 %d 条消息:\n", len(messages))
	fmt.Println("=" + strings.Repeat("=", 100))
	fmt.Printf("%-8s %-12s %-20s %-15s %-20s\n", "ID", "Telegram ID", "频道", "发送时间", "消息长度")
	fmt.Println("-" + strings.Repeat("-", 100))

	for _, msg := range messages {
		channelTitle := "未知频道"
		if channel != nil {
			channelTitle = channel.Title
		} else {
			// 获取频道信息
			if ch, err := db.GetChannelByID(msg.ChannelID); err == nil {
				channelTitle = ch.Title
			}
		}

		textLength := len(msg.Text)
		if textLength > 50 {
			textLength = 50
		}

		fmt.Printf("%-8d %-12d %-20s %-15s %-20d\n",
			msg.ID,
			msg.TelegramID,
			truncateString(channelTitle, 18),
			msg.Date.Format("01-02 15:04"),
			textLength)
	}

	fmt.Println("=" + strings.Repeat("=", 100))

	// 显示详细消息内容
	fmt.Println("\n详细消息内容:")
	fmt.Println("=" + strings.Repeat("=", 100))

	for i, msg := range messages {
		// 获取频道标题
		channelTitle := "未知频道"
		if channel != nil {
			channelTitle = channel.Title
		} else {
			// 获取频道信息
			if ch, err := db.GetChannelByID(msg.ChannelID); err == nil {
				channelTitle = ch.Title
			}
		}

		fmt.Printf("\n[%d] 消息 ID: %d (Telegram ID: %d)\n", i+1, msg.ID, msg.TelegramID)
		fmt.Printf("频道: %s\n", channelTitle)
		fmt.Printf("时间: %s\n", msg.Date.Format("2006-01-02 15:04:05"))
		fmt.Printf("内容:\n%s\n", msg.Text)

		if msg.MediaType != "" {
			fmt.Printf("媒体类型: %s\n", msg.MediaType)
		}
		if msg.MediaURL != "" {
			fmt.Printf("媒体链接: %s\n", msg.MediaURL)
		}
		if msg.Views > 0 {
			fmt.Printf("浏览: %d\n", msg.Views)
		}
		if msg.Forwards > 0 {
			fmt.Printf("转发: %d\n", msg.Forwards)
		}
		fmt.Println("-" + strings.Repeat("-", 50))
	}

	return nil
}
