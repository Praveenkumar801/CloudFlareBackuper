package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type TelegramNotifier struct {
	botToken string
	chatID   string
}

func NewTelegramNotifier(botToken, chatID string) *TelegramNotifier {
	return &TelegramNotifier{
		botToken: botToken,
		chatID:   chatID,
	}
}

type TelegramMessage struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

func (t *TelegramNotifier) SendBackupSuccess(fileName, fileURL string, fileSize int64) error {
	message := fmt.Sprintf(
		"✅ *Backup Successful*\n\n"+
			"A new backup has been created and uploaded successfully!\n\n"+
			"*File Name:* `%s`\n"+
			"*File Size:* %s\n"+
			"*Download Link:* [Click here](%s)",
		fileName,
		formatFileSize(fileSize),
		fileURL,
	)

	return t.sendMessage(message)
}

func (t *TelegramNotifier) SendBackupFailure(err error) error {
	message := fmt.Sprintf(
		"❌ *Backup Failed*\n\n"+
			"The backup process encountered an error.\n\n"+
			"*Error:* `%s`",
		err.Error(),
	)

	return t.sendMessage(message)
}

func (t *TelegramNotifier) SendBackupDeletion(fileName, fileURL string) error {
	message := fmt.Sprintf(
		"🗑️ *Old Backup Deleted*\n\n"+
			"An old backup has been automatically deleted due to retention limit.\n\n"+
			"*Deleted File:* `%s`\n"+
			"*Previous Download Link:* `%s`",
		fileName,
		fileURL,
	)

	return t.sendMessage(message)
}

func (t *TelegramNotifier) sendMessage(text string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.botToken)

	message := TelegramMessage{
		ChatID:    t.chatID,
		Text:      text,
		ParseMode: "Markdown",
	}

	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal Telegram message: %w", err)
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to send Telegram message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Telegram API returned status code: %d", resp.StatusCode)
	}

	return nil
}
