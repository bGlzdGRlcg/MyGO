package bots

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/mmcdole/gofeed"
)

const max_hashes = 500

type FeedItem struct {
	Title       string
	Link        string
	Description string
	Hash        string
}

type RSSWatcher struct {
	URL        string          `json:"url"`
	SeenHashes map[string]bool `json:"seen_hashes"`
	Parser     *gofeed.Parser  `json:"-"`
}

type rssWatcherJSON struct {
	URL        string          `json:"url"`
	SeenHashes map[string]bool `json:"seen_hashes"`
}

func (w *RSSWatcher) cleanOldHashes(currentHashes []string) {
	if len(w.SeenHashes) <= max_hashes {
		return
	}

	current := make(map[string]bool)
	for _, hash := range currentHashes {
		current[hash] = true
	}

	for hash := range w.SeenHashes {
		if !current[hash] {
			delete(w.SeenHashes, hash)
		}
		if len(w.SeenHashes) <= max_hashes {
			break
		}
	}
}

func generateHash(title, link, description string) string {
	h := sha256.New()
	h.Write([]byte(title + link + description))
	return hex.EncodeToString(h.Sum(nil))
}

func NewRSSWatcher(url string) *RSSWatcher {
	return &RSSWatcher{
		URL:        url,
		SeenHashes: make(map[string]bool),
		Parser:     gofeed.NewParser(),
	}
}

func (w *RSSWatcher) CheckNew() ([]FeedItem, error) {
	feed, err := w.Parser.ParseURL(w.URL)
	if err != nil {
		return nil, err
	}

	var newItems []FeedItem
	var currentHashes []string

	for _, item := range feed.Items {
		hash := generateHash(item.Title, item.Link, item.Description)
		currentHashes = append(currentHashes, hash)

		if !w.SeenHashes[hash] {
			newItems = append(newItems, FeedItem{
				Title:       item.Title,
				Link:        item.Link,
				Description: item.Description,
				Hash:        hash,
			})
			w.SeenHashes[hash] = true
		}
	}
	w.cleanOldHashes(currentHashes)

	return newItems, nil
}

func (w *RSSWatcher) MarshalJSON() ([]byte, error) {
	return json.Marshal(&rssWatcherJSON{
		URL:        w.URL,
		SeenHashes: w.SeenHashes,
	})
}

func (w *RSSWatcher) UnmarshalJSON(data []byte) error {
	temp := &rssWatcherJSON{}
	if err := json.Unmarshal(data, temp); err != nil {
		return err
	}
	w.URL = temp.URL
	w.SeenHashes = temp.SeenHashes
	w.Parser = gofeed.NewParser()
	return nil
}

func (w *RSSWatcher) Close() {
	w.SeenHashes = make(map[string]bool)
	w.URL = ""
	w.Parser = nil
}
