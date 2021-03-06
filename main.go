package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func getForbesList(result *string, winners *string, losers *string) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	start := time.Now()

	url := "https://www.forbes.com/real-time-billionaires"

	var res, w, l string
	err := chromedp.Run(ctx,
		emulation.SetUserAgentOverride(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
		chromedp.Navigate(url),
		chromedp.WaitReady("table", chromedp.ByQuery),
		chromedp.WaitVisible("table tr", chromedp.ByQueryAll),
		chromedp.Text("table tbody", &res, chromedp.ByQueryAll),
		chromedp.Text("div > .winners-row", &w, chromedp.ByQueryAll),
		chromedp.Text("div > .losers-row", &l, chromedp.ByQueryAll),
	)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(l)
	fmt.Printf("\nTook: %f secs\n", time.Since(start).Seconds())
	*result = res
	*winners = w
	*losers = l

}

func runBot(result string, losers string, winners string) {
	tgKey := getEnv("TELEGRAM_KEY", "hello")
	bot, err := tgbotapi.NewBotAPI(tgKey)
	// bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		panic(err)
	}

	bot.Debug = true

	// Create a new UpdateConfig struct with an offset of 0. Offsets are used
	// to make sure Telegram knows we've handled previous values and we don't
	// need them repeated.
	updateConfig := tgbotapi.NewUpdate(0)

	// Tell Telegram we should wait up to 30 seconds on each request for an
	// update. This way we can get information just as quickly as making many
	// frequent requests without having to send nearly as many.

	updateConfig.Timeout = 30

	// Start polling Telegram for updates.
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {

		if update.Message == nil {
			continue
		}

		// Create a new MessageConfig. We don't have text yet,
		// so we leave it empty.
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// Extract the command from the Message.
		switch update.Message.Command() {
		case "help":
			msg.Text = "Display current forbes list /list\n\nDisplay today's winners /winners\n\nDisplay today's losers /losers"
		case "start":
			msg.Text = "Display current forbes list /list\n\nDisplay today's winners /winners\n\nDisplay today's losers /losers"
		case "list":
			msg.Text = result
		case "losers":
			msg.Text = losers
		case "winners":
			msg.Text = winners
		default:
			msg.Text = "Hello there ????????\nCommand not available.\nUse /help to view available commands"
		}

		// Okay, we're sending our message off! We don't care about the message
		// we just sent, so we'll discard it.
		if _, err := bot.Send(msg); err != nil {
			fmt.Println("Telegram send error")
		}
	}
}

func main() {
	ticker := time.NewTicker(3600 * time.Second)
	quit := make(chan struct{})
	var result, winners, losers string
	getForbesList(&result, &winners, &losers)
	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("Every hour fetch data")
				getForbesList(&result, &winners, &losers)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	fmt.Printf("Running - success")
	runBot(result, losers, winners)

}

// Gets default value passed if no value exist for given environment variable.
func getEnv(key, fallback string) string {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		return os.Getenv(key)
	}

	return os.Getenv(key)
}
