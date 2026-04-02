package config

import (
	"strings"
	"testing"
)

func skipIfKeychainUnavailable(t *testing.T, err error) {
	t.Helper()
	if err != nil && strings.Contains(err.Error(), "No directory provided for file keyring") {
		t.Skip("Keychain file backend unavailable, skipping test")
	}
}

func TestKeychain_SetAndGet(t *testing.T) {
	// Skip test in CI environments without keychain access
	if testing.Short() {
		t.Skip("Skipping keychain test in short mode")
	}

	kc, err := NewKeychain("pairadmin-test")
	if err != nil {
		t.Skipf("Keychain unavailable: %v", err)
	}
	defer func() {
		_ = kc.Delete("test-key")
	}()

	// Set
	err = kc.Set("test-key", "secret-value")
	if err != nil {
		skipIfKeychainUnavailable(t, err)
		t.Fatalf("Set failed: %v", err)
	}

	// Get
	val, err := kc.Get("test-key")
	if err != nil {
		skipIfKeychainUnavailable(t, err)
		t.Fatalf("Get failed: %v", err)
	}
	if val != "secret-value" {
		t.Errorf("Expected 'secret-value', got %q", val)
	}

	// Delete
	err = kc.Delete("test-key")
	if err != nil {
		skipIfKeychainUnavailable(t, err)
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify deleted
	_, err = kc.Get("test-key")
	if err != nil {
		skipIfKeychainUnavailable(t, err)
		// Expected error - key not found
	} else {
		t.Error("Expected error after deletion")
	}
}
