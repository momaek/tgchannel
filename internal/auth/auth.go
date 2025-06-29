package auth

import (
	"context"
	"fmt"
	"log"

	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
	"golang.org/x/term"
)

type Auth struct {
	apiID       int
	apiHash     string
	sessionFile string
}

// NewAuth 创建新的认证实例
func NewAuth(apiID int, apiHash, sessionFile string) *Auth {
	return &Auth{
		apiID:       apiID,
		apiHash:     apiHash,
		sessionFile: sessionFile,
	}
}

// Login 登录 Telegram 账号
func (a *Auth) Login(ctx context.Context) error {
	client := telegram.NewClient(a.apiID, a.apiHash, telegram.Options{
		SessionStorage: &session.FileStorage{Path: a.sessionFile},
	})

	handler := &authHandler{}
	flow := auth.NewFlow(
		handler,
		auth.SendCodeOptions{},
	)

	return client.Run(ctx, func(ctx context.Context) error {
		if err := client.Auth().IfNecessary(ctx, flow); err != nil {
			return fmt.Errorf("认证失败: %w", err)
		}

		me, err := client.Self(ctx)
		if err != nil {
			return fmt.Errorf("获取用户信息失败: %w", err)
		}
		log.Printf("登录成功! 用户: %s (@%s)", me.FirstName, me.Username)
		return nil
	})
}

// authHandler 处理认证回调
type authHandler struct{}

func (h *authHandler) Phone(ctx context.Context) (string, error) {
	fmt.Print("请输入手机号码 (格式: +86xxxxxxxxxxx): ")
	var phone string
	fmt.Scanln(&phone)
	return phone, nil
}

func (h *authHandler) Password(ctx context.Context) (string, error) {
	fmt.Print("请输入密码: ")
	password, err := term.ReadPassword(0)
	if err != nil {
		return "", err
	}
	fmt.Println()
	return string(password), nil
}

func (h *authHandler) Code(ctx context.Context, _ *tg.AuthSentCode) (string, error) {
	fmt.Print("请输入验证码: ")
	var code string
	fmt.Scanln(&code)
	return code, nil
}

func (h *authHandler) SignUp(ctx context.Context) (auth.UserInfo, error) {
	return auth.UserInfo{}, fmt.Errorf("不支持注册新账号")
}

func (h *authHandler) AcceptTermsOfService(ctx context.Context, tos tg.HelpTermsOfService) error {
	fmt.Println("需要接受服务条款:", tos.Text)
	fmt.Print("是否接受? (y/n): ")
	var input string
	fmt.Scanln(&input)
	if input == "y" || input == "Y" {
		return nil
	}
	return fmt.Errorf("未接受服务条款")
}

// GetClient 获取 Telegram 客户端
func (a *Auth) GetClient() *telegram.Client {
	return telegram.NewClient(a.apiID, a.apiHash, telegram.Options{
		SessionStorage: &session.FileStorage{Path: a.sessionFile},
	})
}
