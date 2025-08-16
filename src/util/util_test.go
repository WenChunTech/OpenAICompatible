package util

import (
	"testing"
)

func TestGenerateUUID(t *testing.T) {
	uuid := GenerateUUID()
	t.Logf("Generated UUID: %s", uuid)
}
