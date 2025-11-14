package notification

import (
	"log"
)

// MultiNotifier sends notifications to multiple notifiers
type MultiNotifier struct {
	notifiers []Notifier
}

// NewMultiNotifier creates a new multi-notifier
func NewMultiNotifier(notifiers ...Notifier) *MultiNotifier {
	return &MultiNotifier{
		notifiers: notifiers,
	}
}

func (m *MultiNotifier) SendBackupSuccess(fileName, fileURL string, fileSize int64) error {
	var lastErr error
	for _, notifier := range m.notifiers {
		if err := notifier.SendBackupSuccess(fileName, fileURL, fileSize); err != nil {
			log.Printf("Failed to send success notification: %v", err)
			lastErr = err
		}
	}
	return lastErr
}

func (m *MultiNotifier) SendBackupFailure(err error) error {
	var lastErr error
	for _, notifier := range m.notifiers {
		if notifyErr := notifier.SendBackupFailure(err); notifyErr != nil {
			log.Printf("Failed to send failure notification: %v", notifyErr)
			lastErr = notifyErr
		}
	}
	return lastErr
}

func (m *MultiNotifier) SendBackupDeletion(fileName, fileURL string) error {
	var lastErr error
	for _, notifier := range m.notifiers {
		if err := notifier.SendBackupDeletion(fileName, fileURL); err != nil {
			log.Printf("Failed to send deletion notification: %v", err)
			lastErr = err
		}
	}
	return lastErr
}
