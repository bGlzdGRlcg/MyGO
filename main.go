package main

import (
	"bGlzdGRlcg/MyGO/bots"
	"bGlzdGRlcg/MyGO/ms"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}
	ms.HOST = os.Getenv("MS_HOST")
	ms.Cid = os.Getenv("MS_CID")
	ms.Secret = os.Getenv("MS_SECRET")
	ms.Token = os.Getenv("MS_TOKEN")
	bots.Start()
}
