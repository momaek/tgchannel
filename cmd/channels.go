package cmd

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gotd/td/tg"
	"github.com/momaek/tgchannel/internal/auth"
	"github.com/spf13/cobra"
)

// channelsCmd represents the channels command
var channelsCmd = &cobra.Command{
	Use:   "channels",
	Short: "列出 Telegram 关注的 Channel",
	Long: `通过 Telegram API 获取当前账号关注的所有 Channel。

显示每个 Channel 的基本信息，包括用户名、标题、成员数量等。
这个命令会调用 Telegram 的 API 来获取实时的关注列表。`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := listChannels(); err != nil {
			log.Fatalf("获取 Channel 列表失败: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(channelsCmd)
}

func listChannels() error {
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
		// 获取当前用户信息
		me, err := client.Self(ctx)
		if err != nil {
			return fmt.Errorf("获取用户信息失败: %w", err)
		}

		fmt.Printf("当前用户: %s (@%s)\n", me.FirstName, me.Username)
		fmt.Println()

		// 获取对话列表（包括 Channel）
		dialogs, err := client.API().MessagesGetDialogs(ctx, &tg.MessagesGetDialogsRequest{
			Limit:      100,                  // 获取前100个对话
			OffsetPeer: &tg.InputPeerEmpty{}, // 添加必需的 offset_peer 字段
		})
		if err != nil {
			return fmt.Errorf("获取对话列表失败: %w", err)
		}

		// 处理对话列表 - 支持多种返回类型
		var chats []tg.ChatClass
		switch d := dialogs.(type) {
		case *tg.MessagesDialogs:
			chats = d.Chats
		case *tg.MessagesDialogsSlice:
			chats = d.Chats
		default:
			return fmt.Errorf("未知的对话响应类型: %T", dialogs)
		}

		// 统计 Channel 数量
		var channels []*tg.Channel
		for _, chat := range chats {
			if channel, ok := chat.(*tg.Channel); ok {
				channels = append(channels, channel)
			}
		}

		if len(channels) == 0 {
			fmt.Println("你还没有关注任何 Channel")
			return nil
		}

		// 显示 Channel 列表
		fmt.Printf("你关注的 Channel 列表 (共 %d 个):\n", len(channels))
		fmt.Println("=" + strings.Repeat("=", 80))
		fmt.Printf("%-12s %-40s %-15s %-10s\n", "ID", "标题", "成员数量", "类型")
		fmt.Println("-" + strings.Repeat("-", 80))

		for _, channel := range channels {
			channelType := "Channel"
			if channel.Broadcast {
				channelType = "Broadcast"
			}

			fmt.Printf("%-12d %-40s %-15d %-10s\n",
				channel.ID,
				truncateString(channel.Title, 38),
				channel.ParticipantsCount,
				channelType)
		}

		fmt.Println("=" + strings.Repeat("=", 80))

		return nil
	})

	if err != nil {
		return fmt.Errorf("获取 Channel 列表失败: %w", err)
	}

	return nil
}
