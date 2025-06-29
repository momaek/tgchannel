package cmd

import (
	"fmt"
	"os"

	"github.com/momaek/tgchannel/internal/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	config  *models.Config
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tgchannel",
	Short: "Telegram Channel 爬虫系统",
	Long: `tgchannel 是一个功能强大的 Telegram Channel 爬虫系统。

主要功能:
- 用户认证和登录
- Channel 订阅管理
- 历史消息抓取
- 实时更新监听
- 数据存储和管理`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// 全局标志
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "配置文件路径 (默认为 ./configs/config.yaml)")
}

// initConfig 读取配置文件
func initConfig() {
	if cfgFile != "" {
		// 使用指定的配置文件
		viper.SetConfigFile(cfgFile)
	} else {
		// 搜索配置文件
		viper.AddConfigPath("./configs")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	// 读取环境变量
	viper.AutomaticEnv()

	// 读取配置文件
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "使用配置文件:", viper.ConfigFileUsed())
	} else {
		fmt.Fprintf(os.Stderr, "读取配置文件失败: %v\n", err)
		os.Exit(1)
	}

	// 解析配置
	config = &models.Config{}
	if err := viper.Unmarshal(config); err != nil {
		fmt.Fprintf(os.Stderr, "解析配置文件失败: %v\n", err)
		os.Exit(1)
	}

}
