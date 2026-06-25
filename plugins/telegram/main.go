package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/smtdfc/nagare/plugin-sdk/plugin"
	"github.com/smtdfc/nagare/plugin-sdk/shared"
)

func main() {

	botApiKey := os.Getenv("TELEGRAM_NAGARE_BOT_API_KEY")
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	state := &BotState{
		userChannels:  make(map[int64]string),
		isInitialized: make(map[int64]bool),
		msgQueue:      make(map[int64][]string),
	}

	plg := plugin.NewPlugin()
	plg.Connect()

	opts := []bot.Option{bot.WithDefaultHandler(NewTelegramHandler(state, plg))}
	b, err := bot.New(botApiKey, opts...)
	if err != nil {
		panic(err)
	}

	go b.Start(ctx)

	plg.Listen(func(msg shared.Message) {
		HandlePluginMessages(ctx, b, plg, state, msg)
	})

	<-ctx.Done()
}
