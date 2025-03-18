package bots

import (
	"bGlzdGRlcg/MyGO/ms"
	"context"
	"fmt"
	"log"

	"github.com/mattn/go-mastodon"
)

var (
	Bot_ID mastodon.ID
)

func Start() {
	c := mastodon.NewClient(&mastodon.Config{
		Server:       ms.HOST,
		ClientID:     ms.Cid,
		ClientSecret: ms.Secret,
		AccessToken:  ms.Token,
	})

	bot, err := c.GetAccountCurrentUser(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	Bot_ID = bot.ID

	ws_client := c.NewWSClient()

	ctx := context.Background()

	events, err := ws_client.StreamingWSUser(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for event := range events {
		switch e := event.(type) {
		case *mastodon.NotificationEvent:
			fmt.Printf("收到新通知: %s\n", e.Notification.Type)
			if e.Notification.Type == "follow" {
				c.AccountFollow(ctx, e.Notification.Account.ID)
			} else if e.Notification.Type == "mention" {
				fmt.Printf("收到新嘟文: %s\n", e.Notification.Status.Content)
				Reply(c, e.Notification.Status)
			}
		case *mastodon.DeleteEvent:
			fmt.Printf("嘟文被删除: %s\n", e.ID)
		case *mastodon.ErrorEvent:
			fmt.Printf("发生错误: %v\n", e.Err)
		}
	}
}
