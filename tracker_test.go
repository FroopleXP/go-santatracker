package santatracker

import "testing"

func TestNewTracker(t *testing.T) {
	_, err := NewTracker(false)
	if err != nil {
		t.Error("Failed to create tracker: ", err)
	}
}
