package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IndrajeethY/CloudFlareBackuper/config"
	"github.com/IndrajeethY/CloudFlareBackuper/notification"
	"github.com/IndrajeethY/CloudFlareBackuper/scheduler"
	"github.com/IndrajeethY/CloudFlareBackuper/storage"
)

var (
	version   = "dev"
	commit    = "none"
	buildDate = "unknown"
)

func main() {

	configPath := flag.String("config", "config.yml", "Path to configuration file")
	runOnce := flag.Bool("once", false, "Run backup once and exit")
	showVersion := flag.Bool("version", false, "Show version information")
	flag.Parse()

	if *showVersion {
		fmt.Printf("CloudFlare Backuper\n")
		fmt.Printf("Version:    %s\n", version)
		fmt.Printf("Commit:     %s\n", commit)
		fmt.Printf("Build Date: %s\n", buildDate)
		return
	}

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Println("Configuration loaded successfully")

	r2Client, err := storage.NewR2Client(
		cfg.CloudFlare.AccountID,
		cfg.CloudFlare.AccessKeyID,
		cfg.CloudFlare.SecretKey,
		cfg.CloudFlare.Bucket,
		cfg.CloudFlare.URI,
	)
	if err != nil {
		log.Fatalf("Failed to initialize R2 client: %v", err)
	}
	log.Println("CloudFlare R2 client initialized")

	// Initialize notifiers based on configuration
	var notifiers []notification.Notifier

	if cfg.Discord.WebhookURL != "" {
		discordNotifier := notification.NewDiscordNotifier(cfg.Discord.WebhookURL)
		notifiers = append(notifiers, discordNotifier)
		log.Println("Discord notifier initialized")
	}

	if cfg.Telegram.BotToken != "" && cfg.Telegram.ChatID != "" {
		telegramNotifier := notification.NewTelegramNotifier(cfg.Telegram.BotToken, cfg.Telegram.ChatID)
		notifiers = append(notifiers, telegramNotifier)
		log.Println("Telegram notifier initialized")
	}

	if len(notifiers) == 0 {
		log.Fatalf("No notification methods configured")
	}

	notifier := notification.NewMultiNotifier(notifiers...)

	backupScheduler := scheduler.NewBackupScheduler(cfg, r2Client, notifier)

	if *runOnce {
		log.Println("Running backup once...")
		if err := backupScheduler.RunOnce(); err != nil {
			log.Fatalf("Backup failed: %v", err)
		}
		log.Println("Backup completed successfully")
		return
	}

	if err := backupScheduler.Start(); err != nil {
		log.Fatalf("Failed to start scheduler: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	log.Println("CloudFlare Backuper is running. Press Ctrl+C to exit.")
	<-sigChan

	log.Println("Shutting down...")
	backupScheduler.Stop()
	log.Println("Shutdown complete")
}
