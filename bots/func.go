package bots

import (
	"bGlzdGRlcg/MyGO/ms"
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"os"
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
	} else if (str[0] == "@"+Bot_Name) && (len(str) == 1 || (str[1] == "mygo" || str[1] == "/mygo")) {
		irand, _ := rand.Int(rand.Reader, big.NewInt(int64(len(MyGO_str))))
		rpy := "@" + status.Account.Username + " " + MyGO_str[irand.Int64()]
		ms.PostdTootr(client, status.ID, rpy, status.Visibility)
	} else {
		o_str := status.Account.Username + ": " + Formatstr(status.Content)
		rpy := "@" + status.Account.Username + " "
		if status.InReplyToID != nil {
			t_sts, _ := client.GetStatus(context.Background(), mastodon.ID(status.InReplyToID.(string)))
			t_str := Formatstr(t_sts.Content)
			o_str += "\n" + "Reply to " + t_sts.Account.Username + ": " + t_str
		}
		//fmt.Println(o_str)
		rpy += Gemini_rpy(o_str)
		ms.PostdTootr(client, status.ID, rpy, status.Visibility)
	}
}

func Formatstr(str string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(str))
	if err != nil {
		log.Fatal(err)
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

func Formatstrs(str string) []string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(str))
	if err != nil {
		log.Fatal(err)
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

func Gemini_rpy(str string) string {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	model := client.GenerativeModel("gemini-2.0-flash-lite")
	model.SetTemperature(0.9)
	model.SystemInstruction = genai.NewUserContent(genai.Text(prompt))
	resp, err := model.GenerateContent(ctx, genai.Text(str))
	if err != nil {
		log.Fatal(err)
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
