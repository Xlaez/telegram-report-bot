package main

import (
	"log"
)

func main() {
	// Create a new bot instance using your Telegram Bot Token.
	bot, err := tgbotapi.NewBotAPI("YOUR_TELEGRAM_BOT_TOKEN")
	if err != nil {
		log.Fatal(err)
	}

	// Uncomment the following line to enable debugging.
	// bot.Debug = true

	// Configure a new update config with the update type and time interval.
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	// Start polling for new updates using the update config.
	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Loop through all incoming updates.
	for update := range updates {
		// Check if the incoming update is a message.
		if update.Message != nil {
			// Check if the incoming message contains the "/report" command.
			if update.Message.IsCommand() && update.Message.Command() == "report" {
				// Send a response message to the user.
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please enter your report:")
				bot.Send(msg)

				// Wait for the user's report.
				report := <-waitForMessage(bot, update.Message.Chat.ID)

				// Send the report to a designated report chat.
				// Inplace of ), your report chat id should be here
				reportMsg := tgbotapi.NewMessage(0, report.Text)
				bot.Send(reportMsg)

				// Send a confirmation message to the user.
				confirmationMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Thank you for your report!")
				bot.Send(confirmationMsg)
			}
		}
	}
}

// This function waits for the next incoming message from the specified chat ID and returns it.
func waitForMessage(bot *tgbotapi.BotAPI, chatID int64) chan tgbotapi.Message {
	messageChan := make(chan tgbotapi.Message)

	go func() {
		for update := range bot.ListenForWebhook("/") {
			if update.Message != nil && update.Message.Chat.ID == chatID {
				messageChan <- *update.Message
				break
			}
		}
	}()

	return messageChan
}
