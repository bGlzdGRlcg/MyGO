package ms

import (
	"context"
	"fmt"

	"github.com/mattn/go-mastodon"
)

func PostdToot(ms_c *mastodon.Client, content, vis string) {
	toot := mastodon.Toot{
		Status:     content,
		Visibility: vis,
	}
	_, err := ms_c.PostStatus(context.Background(), &toot)
	if err != nil {
		fmt.Println(err)
	}
}

func PostdTootr(ms_c *mastodon.Client, ID mastodon.ID, content, vis string) {
	toot := mastodon.Toot{
		Status:      content,
		Visibility:  vis,
		InReplyToID: ID,
	}
	_, err := ms_c.PostStatus(context.Background(), &toot)
	if err != nil {
		fmt.Println(err)
	}
}
