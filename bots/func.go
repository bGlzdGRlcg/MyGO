package bots

import (
	"bGlzdGRlcg/MyGO/ms"
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/generative-ai-go/genai"
	"github.com/mattn/go-mastodon"
	"google.golang.org/api/option"
)

func Reply(client *mastodon.Client, status *mastodon.Status) {
	str := Formatstrs(status.Content)
	if status.Account.ID == Bot_ID {
		return
	} else if (str[0] == "@"+Bot_Name) && (len(str) == 1 || (str[1] == "MyGO" || str[1] == "mygo" || str[1] == "/mygo" || str[1] == "/MyGO")) {
		MyGO_rpy(client, status)
	} else if (str[0] == "@"+Bot_Name) && strings.Contains(str[1], "/sub") {
		Sub(client, status, str)
	} else if (str[0] == "@"+Bot_Name) && strings.Contains(str[1], "/unsub") {
		UnSub(client, status, str)
	} else if (str[0] == "@"+Bot_Name) && strings.Contains(str[1], "/getsublist") {
		ms.PostdTootr(client, status.ID, "@"+status.Account.Username+" \n"+User_map[string(status.Account.ID)].GetSubList(), "direct")
	} else {
		Gemini_rpy(client, status)
	}
}

func Formatstr(str string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(str))
	if err != nil {
		fmt.Printf("%v", err)
	}
	var lines []string
	doc.Find("p").Contents().Each(func(i int, s *goquery.Selection) {
		if s.Is("br") {
			lines = append(lines, "\n")
		} else {
			text := strings.TrimSpace(s.Text())
			if text != "" {
				lines = append(lines, text)
			}
		}
	})
	content := strings.Join(lines, " ")
	sContent := strings.Split(content, "\n")
	var pLines []string

	for i, line := range sContent {
		if strings.TrimSpace(line) != "" {
			if i == 0 {
				pLines = append(pLines, line)
			} else {
				pLines = append(pLines, strings.TrimLeft(line, " "))
			}
		}
	}
	content = strings.Join(pLines, "\n")
	return content
}

func Formatctx(str string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(str))
	if err != nil {
		fmt.Printf("%v", err)
	}
	var lines []string
	doc.Contents().Each(func(i int, s *goquery.Selection) {
		if s.Is("br") {
			lines = append(lines, "\n")
		} else {
			text := strings.TrimSpace(s.Text())
			if text != "" {
				lines = append(lines, text)
			}
		}
	})
	content := strings.Join(lines, " ")
	sContent := strings.Split(content, "\n")
	var pLines []string

	for i, line := range sContent {
		if strings.TrimSpace(line) != "" {
			if i == 0 {
				pLines = append(pLines, line)
			} else {
				pLines = append(pLines, strings.TrimLeft(line, " "))
			}
		}
	}
	content = strings.Join(pLines, "\n")
	return content
}

func Formatstrs(str string) []string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(str))
	if err != nil {
		fmt.Printf("%v", err)
	}
	var lines []string
	doc.Find("p").Contents().Each(func(i int, s *goquery.Selection) {
		if !s.Is("br") {
			text := strings.TrimSpace(s.Text())
			if text != "" {
				lines = append(lines, text)
			}
		}
	})
	return lines
}

func Gemini(str, pt string) string {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		fmt.Printf("%v", err)
	}
	defer client.Close()
	model := client.GenerativeModel("gemini-2.0-flash-lite")
	model.SetTemperature(0.9)
	model.SystemInstruction = genai.NewUserContent(genai.Text(pt))
	resp, err := model.GenerateContent(ctx, genai.Text(str))
	if err != nil {
		fmt.Printf("%v", err)
	}
	respStr := ""
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				respStr += fmt.Sprintln(part)
			}
		}
	}
	//fmt.Println(respStr)
	return respStr
}

func Gemini_rpy(client *mastodon.Client, status *mastodon.Status) {
	o_str := status.Account.Username + ": " + Formatstr(status.Content)
	rpy := "@" + status.Account.Username + " "
	if status.InReplyToID != nil {
		t_sts, _ := client.GetStatus(context.Background(), mastodon.ID(status.InReplyToID.(string)))
		t_str := Formatstr(t_sts.Content)
		o_str += "\n" + "Reply to " + t_sts.Account.Username + ": " + t_str
	}
	//fmt.Println(o_str)
	rpy += Gemini(o_str, prompt)
	ms.PostdTootr(client, status.ID, rpy, status.Visibility)
}

func MyGO_rpy(client *mastodon.Client, status *mastodon.Status) {
	irand, _ := rand.Int(rand.Reader, big.NewInt(int64(len(MyGO_str))))
	rpy := "@" + status.Account.Username + " " + MyGO_str[irand.Int64()]
	ms.PostdTootr(client, status.ID, rpy, status.Visibility)
}

func Sub(client *mastodon.Client, status *mastodon.Status, str []string) {
	ss := strings.Split(str[1], " ")
	if len(str) > 2 {
		if strings.Contains(str[2], "http") {
			if User_map[string(status.Account.ID)] == nil {
				User_map[string(status.Account.ID)] = &User{Userid: string(status.Account.ID)}
			}
			_, err := http.Get(str[2])
			if err != nil {
				fmt.Printf("Error checking URL: %v", err)
				ms.PostdTootr(client, status.ID, "@"+status.Account.Username+" 你这链接有问题啊", "direct")
				return
			}
			User_map[string(status.Account.ID)].AddRSSFeed(str[2])
			saveUserMap()
			ms.PostdTootr(client, status.ID, "@"+status.Account.Username+" 订阅成功", "direct")
			return
		}
		ms.PostdTootr(client, status.ID, "@"+status.Account.Username+" 塞...不下的...", "direct")
		return
	}
	if len(ss) == 1 {
		ms.PostdTootr(client, status.ID, "@"+status.Account.Username+" 你的订阅链接去哪里了？", "direct")
	} else if len(ss) > 2 {
		ms.PostdTootr(client, status.ID, "@"+status.Account.Username+" 塞...不下的...", "direct")
	} else if len(ss) == 2 {
		_, err := http.Get(ss[2])
		if err != nil {
			fmt.Printf("Error checking URL: %v", err)
			ms.PostdTootr(client, status.ID, "@"+status.Account.Username+" 你这链接有问题啊", "direct")
			return
		}
		if User_map[string(status.Account.ID)] == nil {
			User_map[string(status.Account.ID)] = &User{Userid: string(status.Account.ID)}
		}
		User_map[string(status.Account.ID)].AddRSSFeed(ss[2])
		saveUserMap()
		ms.PostdTootr(client, status.ID, "@"+status.Account.Username+" 订阅成功", "direct")
		return
	}
}

func UnSub(client *mastodon.Client, status *mastodon.Status, str []string) {
	if len(User_map[string(status.Account.ID)].Subs) == 0 {
		ms.PostdTootr(client, status.ID, "@"+status.Account.Username+" 你还没订阅过任何RSS", "direct")
		return
	}
	strs := strings.Split(str[1], " ")
	if len(strs) == 1 {
		ms.PostdTootr(client, status.ID, "@"+status.Account.Username+" 请给我要退订的RSS的ID", "direct")
		return
	}
	if len(strs) > 2 {
		ms.PostdTootr(client, status.ID, "@"+status.Account.Username+" 塞...不下的...", "direct")
		return
	}
	if len(strs) == 2 {
		if strs[1] == "all" {
			for i := range User_map[string(status.Account.ID)].Subs {
				User_map[string(status.Account.ID)].RmSub(i)
			}
			saveUserMap()
			ms.PostdTootr(client, status.ID, "@"+status.Account.Username+" 退订成功", "direct")
			return
		}
		id, err := strconv.Atoi(strs[1])
		if err != nil {
			ms.PostdTootr(client, status.ID, "@"+status.Account.Username+" 请给我一个正确的ID", "direct")
			return
		}
		User_map[string(status.Account.ID)].RmSub(id)
		saveUserMap()
		ms.PostdTootr(client, status.ID, "@"+status.Account.Username+" 退订成功", "direct")
	}
}
