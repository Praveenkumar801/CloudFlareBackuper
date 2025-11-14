package scheduler

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/IndrajeethY/CloudFlareBackuper/backup"
	"github.com/IndrajeethY/CloudFlareBackuper/config"
	"github.com/IndrajeethY/CloudFlareBackuper/notification"
	"github.com/IndrajeethY/CloudFlareBackuper/storage"
	"github.com/robfig/cron/v3"
)

type BackupScheduler struct {
	config   *config.Config
	r2Client *storage.R2Client
	notifier notification.Notifier
	cron     *cron.Cron
	tempDir  string
}

func NewBackupScheduler(cfg *config.Config, r2Client *storage.R2Client, notifier notification.Notifier) *BackupScheduler {
	return &BackupScheduler{
		config:   cfg,
		r2Client: r2Client,
		notifier: notifier,
		cron:     cron.New(),
		tempDir:  os.TempDir(),
	}
}

func (s *BackupScheduler) Start() error {

	_, err := s.cron.AddFunc(s.config.Backup.Schedule, func() {
		if err := s.runBackup(); err != nil {
			log.Printf("Backup failed: %v", err)
			if notifyErr := s.notifier.SendBackupFailure(err); notifyErr != nil {
				log.Printf("Failed to send failure notification: %v", notifyErr)
			}
		}
	})
	if err != nil {
		return fmt.Errorf("failed to schedule backup: %w", err)
	}

	s.cron.Start()
	log.Printf("Backup scheduler started with schedule: %s", s.config.Backup.Schedule)

	log.Println("Running initial backup...")
	if err := s.runBackup(); err != nil {
		log.Printf("Initial backup failed: %v", err)
		if notifyErr := s.notifier.SendBackupFailure(err); notifyErr != nil {
			log.Printf("Failed to send failure notification: %v", notifyErr)
		}
	}

	return nil
}

func (s *BackupScheduler) Stop() {
	s.cron.Stop()
	log.Println("Backup scheduler stopped")
}

func (s *BackupScheduler) runBackup() error {
	log.Println("Starting backup process...")

	fileName := backup.GenerateBackupFilename(s.config.Backup.NamePrefix)
	archivePath := filepath.Join(s.tempDir, fileName)

	log.Printf("Creating archive from %d folder(s)...", len(s.config.Backup.Folders))
	if err := backup.CreateArchive(s.config.Backup.Folders, archivePath); err != nil {
		return fmt.Errorf("failed to create archive: %w", err)
	}
	defer os.Remove(archivePath)

	fileInfo, err := os.Stat(archivePath)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}
	fileSize := fileInfo.Size()
	log.Printf("Archive created: %s (size: %d bytes)", fileName, fileSize)

	log.Println("Uploading to CloudFlare R2...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	fileURL, err := s.r2Client.UploadFile(ctx, archivePath)
	if err != nil {
		return fmt.Errorf("failed to upload to R2: %w", err)
	}
	log.Printf("Upload successful: %s", fileURL)

	if s.config.Backup.RetentionLimit > 0 {
		log.Printf("Checking for old backups to delete (retention limit: %d)...", s.config.Backup.RetentionLimit)
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		deletedFiles, err := s.r2Client.CleanupOldBackups(ctx, s.config.Backup.NamePrefix, s.config.Backup.RetentionLimit)
		if err != nil {
			log.Printf("Failed to cleanup old backups: %v", err)
		} else if len(deletedFiles) > 0 {
			log.Printf("Deleted %d old backup(s)", len(deletedFiles))
			for _, deletedFile := range deletedFiles {
				deletedFileURL := fmt.Sprintf("%s/%s", s.config.CloudFlare.URI, deletedFile)
				if err := s.notifier.SendBackupDeletion(deletedFile, deletedFileURL); err != nil {
					log.Printf("Failed to send deletion notification for %s: %v", deletedFile, err)
				} else {
					log.Printf("Sent deletion notification for: %s", deletedFile)
				}
			}
		} else {
			log.Println("No old backups to delete")
		}
	}

	log.Println("Sending success notification...")
	if err := s.notifier.SendBackupSuccess(fileName, fileURL, fileSize); err != nil {
		log.Printf("Failed to send notification: %v", err)
	}

	log.Println("Backup completed successfully!")
	return nil
}

func (s *BackupScheduler) RunOnce() error {
	return s.runBackup()
}
