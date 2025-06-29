package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/momaek/tgchannel/internal/database"
	"github.com/spf13/cobra"
)

var (
	listUserUsername string
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出订阅的 Channel",
	Long: `列出指定用户订阅的所有 Channel。

显示每个 Channel 的基本信息，包括用户名、标题、成员数量等。`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := list(); err != nil {
			log.Fatalf("列出失败: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// 添加标志
	listCmd.Flags().StringVarP(&listUserUsername, "user", "u", "", "用户用户名")
	listCmd.MarkFlagRequired("user")
}

func list() error {
	// 初始化数据库
	db, err := database.NewDatabase(config.Database.Path)
	if err != nil {
		return fmt.Errorf("初始化数据库失败: %w", err)
	}
	defer db.Close()

	// 获取用户信息
	user, err := db.GetUserByUsername(listUserUsername)
	if err != nil {
		return fmt.Errorf("用户不存在: %w", err)
	}

	// 获取用户订阅的频道
	channels, err := db.GetUserSubscriptions(user.ID)
	if err != nil {
		return fmt.Errorf("获取订阅失败: %w", err)
	}

	if len(channels) == 0 {
		fmt.Printf("用户 %s 还没有订阅任何 Channel\n", listUserUsername)
		return nil
	}

	// 显示订阅的频道列表
	fmt.Printf("用户 %s 订阅的 Channel 列表:\n", listUserUsername)
	fmt.Println("=" + strings.Repeat("=", 80))
	fmt.Printf("%-20s %-30s %-15s %-10s\n", "用户名", "标题", "成员数量", "状态")
	fmt.Println("-" + strings.Repeat("-", 80))

	for _, channel := range channels {
		status := "活跃"
		if !channel.IsActive {
			status = "非活跃"
		}
		fmt.Printf("%-20s %-30s %-15d %-10s\n",
			"@"+channel.Username,
			truncateString(channel.Title, 28),
			channel.MemberCount,
			status)
	}

	fmt.Println("=" + strings.Repeat("=", 80))
	fmt.Printf("总计: %d 个 Channel\n", len(channels))

	return nil
}

// truncateString 截断字符串到指定长度
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
