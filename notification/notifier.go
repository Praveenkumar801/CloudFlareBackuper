package notification

// Notifier defines the interface for sending backup notifications
type Notifier interface {
	SendBackupSuccess(fileName, fileURL string, fileSize int64) error
	SendBackupFailure(err error) error
	SendBackupDeletion(fileName, fileURL string) error
}
