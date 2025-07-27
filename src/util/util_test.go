package util

import (
	"testing"
)

func TestGenerateUUID(t *testing.T) {
	uuid, err := GenerateUUID()
	if err != nil {
		t.Errorf("Failed to generate UUID: %v", err)
	}
	t.Logf("Generated UUID: %s", uuid)
}
