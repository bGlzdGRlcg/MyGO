package bots

import (
	"bGlzdGRlcg/MyGO/ms"
	"context"
	"fmt"
	"time"

	"github.com/mattn/go-mastodon"
)

var (
	Bot_ID   mastodon.ID
	Bot_Name string
)

func Start() {
	c := mastodon.NewClient(&mastodon.Config{
		Server:       ms.HOST,
		ClientID:     ms.Cid,
		ClientSecret: ms.Secret,
		AccessToken:  ms.Token,
	})

	if err := loadUserMap(); err != nil {
		fmt.Println(err)
	}

	bot, err := c.GetAccountCurrentUser(context.Background())
	if err != nil {
		fmt.Printf("%v", err)
	}
	Bot_ID = bot.ID
	Bot_Name = bot.Username

	ws_client := c.NewWSClient()

	ctx := context.Background()

	events, err := ws_client.StreamingWSUser(ctx)
	if err != nil {
		fmt.Printf("%v", err)
	}

	tag_events, err := ws_client.StreamingWSHashtag(ctx, "MyGO", true)
	if err != nil {
		fmt.Printf("%v", err)
	}

	go func() {
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
	}()

	go func() {
		for tag_event := range tag_events {
			switch e := tag_event.(type) {
			case *mastodon.UpdateEvent:
				fmt.Printf("收到新嘟文: %s\n", e.Status.Content)
				MyGO_rpy(c, e.Status)
			case *mastodon.ErrorEvent:
				fmt.Printf("发生错误: %v\n", e.Err)
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for range ticker.C {
			for _, user := range User_map {
				if user.RSSFeeds != nil {
					for _, watcher := range user.RSSFeeds {
						items, err := watcher.CheckNew()
						if err != nil {
							fmt.Printf("Error checking RSS feed: %v\n", err)
							continue
						}
						m_user, err := ms.GetAcc(c, user.Userid)
						if err != nil {
							fmt.Printf("Error getting Mastodon user: %v\n", err)
							continue
						}
						for _, item := range items {
							content := fmt.Sprintf("@%s \n📰 %s\n\n %s\n\nLink: %s",
								m_user.Username,
								item.Title,
								Formatctx(item.Description),
								item.Link)
							ms.PostdToot(c, content, "direct")
						}
					}
				}
			}
			saveUserMap()
		}
	}()

	select {}
}
