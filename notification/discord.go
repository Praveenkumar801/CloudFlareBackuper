package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type DiscordNotifier struct {
	webhookURL string
}

func NewDiscordNotifier(webhookURL string) *DiscordNotifier {
	return &DiscordNotifier{
		webhookURL: webhookURL,
	}
}

type DiscordMessage struct {
	Content string         `json:"content,omitempty"`
	Embeds  []DiscordEmbed `json:"embeds,omitempty"`
}

type DiscordEmbed struct {
	Title       string              `json:"title,omitempty"`
	Description string              `json:"description,omitempty"`
	Color       int                 `json:"color,omitempty"`
	Fields      []DiscordEmbedField `json:"fields,omitempty"`
	Timestamp   string              `json:"timestamp,omitempty"`
}

type DiscordEmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

func (d *DiscordNotifier) SendBackupSuccess(fileName, fileURL string, fileSize int64) error {
	embed := DiscordEmbed{
		Title:       "✅ Backup Successful",
		Description: "A new backup has been created and uploaded successfully!",
		Color:       3066993,
		Fields: []DiscordEmbedField{
			{
				Name:   "File Name",
				Value:  fileName,
				Inline: false,
			},
			{
				Name:   "File Size",
				Value:  formatFileSize(fileSize),
				Inline: true,
			},
			{
				Name:   "Download Link",
				Value:  fmt.Sprintf("[Click here to download](%s)", fileURL),
				Inline: false,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	message := DiscordMessage{
		Embeds: []DiscordEmbed{embed},
	}

	return d.sendMessage(message)
}

func (d *DiscordNotifier) SendBackupFailure(err error) error {
	embed := DiscordEmbed{
		Title:       "❌ Backup Failed",
		Description: "The backup process encountered an error.",
		Color:       15158332,
		Fields: []DiscordEmbedField{
			{
				Name:   "Error",
				Value:  err.Error(),
				Inline: false,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	message := DiscordMessage{
		Embeds: []DiscordEmbed{embed},
	}

	return d.sendMessage(message)
}

func (d *DiscordNotifier) SendBackupDeletion(fileName, fileURL string) error {
	embed := DiscordEmbed{
		Title:       "🗑️ Old Backup Deleted",
		Description: "An old backup has been automatically deleted due to retention limit.",
		Color:       16776960,
		Fields: []DiscordEmbedField{
			{
				Name:   "Deleted File",
				Value:  fileName,
				Inline: false,
			},
			{
				Name:   "Previous Download Link",
				Value:  fmt.Sprintf("`%s`", fileURL),
				Inline: false,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	message := DiscordMessage{
		Embeds: []DiscordEmbed{embed},
	}

	return d.sendMessage(message)
}

func (d *DiscordNotifier) sendMessage(message DiscordMessage) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal Discord message: %w", err)
	}

	resp, err := http.Post(d.webhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to send Discord webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Discord webhook returned status code: %d", resp.StatusCode)
	}

	return nil
}

func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
