package ms

import (
	"context"

	"github.com/mattn/go-mastodon"
)

func GetAccList(ms_c *mastodon.Client, str string) (*mastodon.Results, error) {
	acc, err := ms_c.Search(context.Background(), str, false)
	return acc, err
}

func GetAcc(ms_c *mastodon.Client, str string) (*mastodon.Account, error) {
	acc, err := ms_c.GetAccount(context.Background(), mastodon.ID(str))
	return acc, err
}
