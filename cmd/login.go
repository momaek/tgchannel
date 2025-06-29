package cmd

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/momaek/tgchannel/internal/auth"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "登录 Telegram 账号",
	Long: `登录到 Telegram 账号。

首次登录需要提供手机号码和验证码。
如果启用了两步验证，还需要输入密码。`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := login(); err != nil {
			log.Fatalf("登录失败: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

func login() error {
	// 添加调试信息
	fmt.Printf("调试信息 - API ID: '%s', API Hash: '%s'\n", config.Telegram.APIID, config.Telegram.APIHash)

	// 解析 API ID
	apiID, err := strconv.Atoi(config.Telegram.APIID)
	if err != nil {
		return fmt.Errorf("无效的 API ID: %w", err)
	}

	// 创建认证实例
	authClient := auth.NewAuth(apiID, config.Telegram.APIHash, config.Telegram.SessionFile)

	// 执行登录
	ctx := context.Background()
	if err := authClient.Login(ctx); err != nil {
		return fmt.Errorf("登录失败: %w", err)
	}

	log.Println("登录成功!")
	return nil
}
