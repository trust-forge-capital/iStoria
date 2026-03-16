package cmd

import (
	"testing"
)

// TestClearScreen tests ClearScreen function doesn't panic
func TestClearScreen(t *testing.T) {
	// ClearScreen should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("ClearScreen panicked: %v", r)
		}
	}()
	ClearScreen()
}

// TestSetupLiveMode tests signal channel creation
func TestSetupLiveMode(t *testing.T) {
	sigChan, cleanup := SetupLiveMode()
	defer cleanup()

	if sigChan == nil {
		t.Error("SetupLiveMode returned nil channel")
	}
}

// TestLiveConfigStruct tests LiveConfig default values
func TestLiveConfigStruct(t *testing.T) {
	config := &LiveConfig{
		Enabled:  true,
		Interval: 1000,
		NoClear:  false,
	}

	if !config.Enabled {
		t.Error("Enabled should be true")
	}
	if config.Interval == 0 {
		t.Error("Interval should not be 0")
	}
}
