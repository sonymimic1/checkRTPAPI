package main

import (
	"context"

	"sonymimic1/Golang_server/checkRTP/config"
	"sonymimic1/Golang_server/checkRTP/internal/app"
)

// @title Check RTP API and Sechdule Clear RTP task
// @version 1.0
// @description This is a check RTP API for gamecode and Sechdule to Clear bet/win value in Redis.
func main() {
	cfg, err := config.LoadConfig("config")
	if err != nil {
		panic(err)
	}

	app.NewApp(context.Background(), cfg).Run()
}
