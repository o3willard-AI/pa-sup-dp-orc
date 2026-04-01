package llm

import (
	"testing"
)

func TestGatewayInterface(t *testing.T) {
	// This test ensures the interface is defined correctly.
	var gw Gateway
	_ = gw
	t.Log("Gateway interface defined")
}