package config

import (
	"testing"
)

func TestKeychain_SetAndGet(t *testing.T) {
	// Skip test in CI environments without keychain access
	if testing.Short() {
		t.Skip("Skipping keychain test in short mode")
	}

	// Set up environment for file-based keyring during tests
	tempDir := t.TempDir()
	t.Setenv("KEYRING_BACKEND", "file")
	t.Setenv("KEYRING_FILE_DIR", tempDir)

	kc, err := NewKeychain("pairadmin-test")
	if err != nil {
		t.Fatalf("NewKeychain failed: %v", err)
	}
	defer func() {
		_ = kc.Delete("test-key")
	}()

	// Set
	err = kc.Set("test-key", "secret-value")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Get
	val, err := kc.Get("test-key")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if val != "secret-value" {
		t.Errorf("Expected 'secret-value', got %q", val)
	}

	// Delete
	err = kc.Delete("test-key")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify deleted
	_, err = kc.Get("test-key")
	if err == nil {
		t.Error("Expected error after deletion")
	}
}