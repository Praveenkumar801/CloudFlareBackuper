package notification

import (
	"errors"
	"testing"
)

// TestNotifierInterface verifies that all notifiers implement the Notifier interface
func TestNotifierInterface(t *testing.T) {
	// Verify Discord notifier implements interface
	var _ Notifier = &DiscordNotifier{}

	// Verify Telegram notifier implements interface
	var _ Notifier = &TelegramNotifier{}

	// Verify Multi notifier implements interface
	var _ Notifier = &MultiNotifier{}
}

// TestMultiNotifier verifies that MultiNotifier sends to all notifiers
func TestMultiNotifier(t *testing.T) {
	// Create mock notifiers
	mock1 := &mockNotifier{}
	mock2 := &mockNotifier{}

	multi := NewMultiNotifier(mock1, mock2)

	// Test SendBackupSuccess
	err := multi.SendBackupSuccess("test.tar.gz", "http://example.com/test.tar.gz", 1024)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !mock1.successCalled {
		t.Error("Expected mock1.SendBackupSuccess to be called")
	}
	if !mock2.successCalled {
		t.Error("Expected mock2.SendBackupSuccess to be called")
	}

	// Test SendBackupFailure
	testErr := errors.New("test error")
	err = multi.SendBackupFailure(testErr)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !mock1.failureCalled {
		t.Error("Expected mock1.SendBackupFailure to be called")
	}
	if !mock2.failureCalled {
		t.Error("Expected mock2.SendBackupFailure to be called")
	}

	// Test SendBackupDeletion
	err = multi.SendBackupDeletion("old.tar.gz", "http://example.com/old.tar.gz")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !mock1.deletionCalled {
		t.Error("Expected mock1.SendBackupDeletion to be called")
	}
	if !mock2.deletionCalled {
		t.Error("Expected mock2.SendBackupDeletion to be called")
	}
}

// mockNotifier is a mock implementation of Notifier for testing
type mockNotifier struct {
	successCalled  bool
	failureCalled  bool
	deletionCalled bool
}

func (m *mockNotifier) SendBackupSuccess(fileName, fileURL string, fileSize int64) error {
	m.successCalled = true
	return nil
}

func (m *mockNotifier) SendBackupFailure(err error) error {
	m.failureCalled = true
	return nil
}

func (m *mockNotifier) SendBackupDeletion(fileName, fileURL string) error {
	m.deletionCalled = true
	return nil
}
