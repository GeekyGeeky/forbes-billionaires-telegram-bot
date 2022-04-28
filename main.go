package main

// "os"

//
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

func getForbesList() string {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	start := time.Now()

	url := "https://www.forbes.com/real-time-billionaires"

	var res string
	// var nodes []*cdp.Node
	err := chromedp.Run(ctx,
		emulation.SetUserAgentOverride(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
		chromedp.Navigate(url),
		// chromedp.ScrollIntoView(`table`),
		chromedp.WaitReady("table", chromedp.ByQuery),
		chromedp.WaitVisible("table tr", chromedp.ByQueryAll),
		chromedp.Text("table tbody", &res, chromedp.ByQueryAll),
	)

	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(strings.TrimSpace(res))
	// fmt.Printf("%s___________\n", res)
	fmt.Printf("\nTook: %f secs\n", time.Since(start).Seconds())
	return res
	// for _, n := range nodes {
	// 	u := n.Children
	// 	fmt.Println(u)
	// }

}

func runBot(result string) {
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
	// get forbes list

	// Let's go through each update that we're getting from Telegram.
	for update := range updates {
		// Telegram can send many types of updates depending on what your Bot
		// is up to. We only want to look at messages for now, so we can
		// discard any other updates.
		if update.Message == nil {
			continue
		}

		// Now that we know we've gotten a new message, we can construct a
		// reply! We'll take the Chat ID and Text from the incoming message
		// and use it to create a new message.
		///	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		// We'll also say that this message is a reply to the previous message.
		// For any other specifications than Chat ID or Text, you'll need to
		// set fields on the `MessageConfig`.
		///	msg.ReplyToMessageID = update.Message.MessageID

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
			// tgbotapi
		default:
			msg.Text = "Hello there 👋🏻\nCommand not available.\nUse /help to view available commands"
		}

		// Okay, we're sending our message off! We don't care about the message
		// we just sent, so we'll discard it.
		if _, err := bot.Send(msg); err != nil {
			fmt.Println("Telegram send error")
			// Note that panics are a bad way to handle errors. Telegram can
			// have service outages or network errors, you should retry sending
			// messages or more gracefully handle failures.
			// panic(err)
		}
	}
}

func main() {
	ticker := time.NewTicker(3600 * time.Second)
	quit := make(chan struct{})
	result := getForbesList()
	go func() {
		for {
			select {
			case <-ticker.C:
				// do stuff
				fmt.Println("Every hour fetch data")
				result = getForbesList()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	fmt.Printf("Running - success")
	runBot(result)
}

// Gets default value passed if no value exist for given environment variable.
func getEnv(key, fallback string) string {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		return os.Getenv(key)
		// log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
	// web: bin/forbes-billionaires-telegram-bot

}
