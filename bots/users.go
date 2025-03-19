package bots

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

var (
	User_map     = make(map[string]*User)
	userMapMutex sync.RWMutex
)

type User struct {
	Userid   string
	Subs     []string
	RSSFeeds map[string]*RSSWatcher
}

func (u *User) RmSub(index int) {
	u.RSSFeeds[u.Subs[index]].Close()
	delete(u.RSSFeeds, u.Subs[index])
	u.Subs = append(u.Subs[:index], u.Subs[index+1:]...)
}

func (u *User) AddSub(sub string) {
	u.Subs = append(u.Subs, sub)
}

func (u *User) GetSubList() string {
	var res string
	for i, sub := range u.Subs {
		res += fmt.Sprint(i) + " - " + sub + "\n"
	}
	return res
}

func (u *User) AddRSSFeed(url string) {
	if u.RSSFeeds == nil {
		u.RSSFeeds = make(map[string]*RSSWatcher)
	}
	if u.RSSFeeds[url] != nil {
		return
	}
	u.AddSub(url)
	u.RSSFeeds[url] = NewRSSWatcher(url)
	u.RSSFeeds[url].CheckNew()
}

func saveUserMap() error {
	userMapMutex.RLock()
	defer userMapMutex.RUnlock()

	data, err := json.Marshal(User_map)
	if err != nil {
		return err
	}

	return os.WriteFile("user.json", data, 0644)
}

func loadUserMap() error {
	userMapMutex.Lock()
	defer userMapMutex.Unlock()

	data, err := os.ReadFile("user.json")
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return json.Unmarshal(data, &User_map)
}
