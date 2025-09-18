package telegram

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gotd/td/examples"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
)

// GetChannelAccessHashSimple 简化版本，尝试不同的方法
func GetChannelAccessHashSimple(channelID int64) (int64, error) {
	appID, _ := strconv.Atoi(os.Getenv("TELEGRAM_APP_ID"))
	appHash := os.Getenv("TELEGRAM_APP_HASH")

	if appID == 0 || appHash == "" {
		return 0, fmt.Errorf("请设置环境变量 TELEGRAM_APP_ID 和 TELEGRAM_APP_HASH")
	}

	fmt.Printf("使用简化方法连接 Telegram API (App ID: %d)...\n", appID)

	// 使用文件会话存储，这样可以持久化认证状态
	sessionStorage := &session.FileStorage{
		Path: "telegram_session.json",
	}

	// 使用会话存储选项
	client := telegram.NewClient(appID, appHash, telegram.Options{
		SessionStorage: sessionStorage,
	})

	// 更长的超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	var accessHash int64

	err := client.Run(ctx, func(ctx context.Context) error {
		// 首先进行认证
		fmt.Println("开始用户认证流程...")

		// 检查是否需要认证
		fmt.Println("尝试进行用户认证...")

		// 使用简单的认证方法
		authFlow := auth.NewFlow(
			examples.Terminal{
				PhoneNumber: os.Getenv("PHONE_NUMBER"),
			},
			auth.SendCodeOptions{},
		)

		if err := client.Auth().IfNecessary(ctx, authFlow); err != nil {
			fmt.Printf("认证失败: %v\n", err)
			fmt.Println("注意：Telegram 需要用户认证才能访问对话列表")
			fmt.Println("请确保：")
			fmt.Println("1. 设置了正确的 TELEGRAM_PHONE 环境变量")
			fmt.Println("2. 手机号格式正确（如：+8613800138000）")
			fmt.Println("3. 网络连接正常")
			fmt.Println("4. 如果收到验证码，请在控制台输入")
			return fmt.Errorf("认证失败: %w", err)
		}

		fmt.Println("✅ 用户认证成功！")

		api := client.API()

		fmt.Println("尝试获取配置信息...")
		config, err := api.HelpGetConfig(ctx)
		if err != nil {
			fmt.Printf("获取配置失败: %v\n", err)
			// 即使配置获取失败，也继续尝试
		} else {
			fmt.Printf("配置获取成功，DC ID: %d\n", config.ThisDC)
		}

		fmt.Println("尝试获取对话列表...")

		// 分步骤获取，减少单次请求的负载
		for limit := 10; limit <= 1000; limit += 20 {
			fmt.Printf("尝试获取 %d 个对话...\n", limit)

			dialogs, err := api.MessagesGetDialogs(ctx, &tg.MessagesGetDialogsRequest{
				OffsetDate: 0,
				OffsetID:   0,
				OffsetPeer: &tg.InputPeerEmpty{},
				Limit:      limit,
				Hash:       0,
			})

			if err != nil {
				fmt.Printf("获取 %d 个对话失败: %v\n", limit, err)
				if limit < 30 {
					continue // 尝试更少的对话数量
				}
				return fmt.Errorf("获取对话列表失败: %w", err)
			}

			// 查找目标频道
			found := false
			switch d := dialogs.(type) {
			case *tg.MessagesDialogs:
				for _, chat := range d.Chats {
					if channel, ok := chat.(*tg.Channel); ok {
						if channel.ID == channelID {
							accessHash = channel.AccessHash
							fmt.Printf("✅ 找到目标频道: %s (AccessHash: %d)\n", channel.Title, accessHash)
							found = true
							break
						}
					}
				}
			case *tg.MessagesDialogsSlice:
				for _, chat := range d.Chats {
					if channel, ok := chat.(*tg.Channel); ok {
						if channel.ID == channelID {
							accessHash = channel.AccessHash
							fmt.Printf("✅ 找到目标频道: %s (AccessHash: %d)\n", channel.Title, accessHash)
							found = true
							break
						}
					}
				}
			}

			if found {
				return nil
			}

			// 如果在少量对话中没找到，继续获取更多
			fmt.Printf("在前 %d 个对话中未找到目标频道，继续搜索...\n", limit)
		}

		return fmt.Errorf("在所有对话中未找到频道 ID: %d", channelID)
	})

	if err != nil {
		return 0, fmt.Errorf("telegram 客户端运行出错: %w", err)
	}

	return accessHash, nil
}
