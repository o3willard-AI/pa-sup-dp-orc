//go:build !linux

package terminal

// Stub for non-Linux: implements interface with noops/errors.
type ATSPITerminalAdapter struct{}

func (a *ATSPITerminalAdapter) Name() string {
	return "AT-SPI2 (disabled: !linux)"
}

func (a *ATSPITerminalAdapter) Available() bool {
	return false
}

func (a *ATSPITerminalAdapter) ListSessions() ([]string, error) {
	return nil, nil
}

func (a *ATSPITerminalAdapter) Capture(_ string) (string, error) {
	return "", nil
}

func (a *ATSPITerminalAdapter) Subscribe(_ string, _ func(string)) error {
	return nil
}

func (a *ATSPITerminalAdapter) GetDimensions(_ string) (int, int, error) {
	return 0, 0, nil
}
