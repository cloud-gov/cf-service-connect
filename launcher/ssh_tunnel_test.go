package launcher

import (
	"os"
	"testing"
)

func TestGetCFBinaryName(t *testing.T) {
	cfBinaryName := getCFBinaryName()
	if cfBinaryName != "cf" {
		t.Fatalf("expected CF binary named 'cf', got %s", cfBinaryName)
	}

	os.Setenv("CF_BINARY_NAME", "cf7")
	cfBinaryName = getCFBinaryName()
	if cfBinaryName != "cf7" {
		t.Fatalf("expected CF binary named 'cf7', got %s", cfBinaryName)
	}
}
