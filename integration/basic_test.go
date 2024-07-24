//go:build integration

package integration

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/AllSage/AiAdmin/api"
	"github.com/stretchr/testify/require"
)

func TestOrcaMiniBlueSky(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	// Set up the test data
	req := api.GenerateRequest{
		Model:  "orca-mini",
		Prompt: "why is the sky blue?",
		Stream: &stream,
		Options: map[string]interface{}{
			"temperature": 0,
			"seed":        123,
		},
	}
	GenerateTestHelper(ctx, t, req, []string{"rayleigh", "scattering"})
}

func TestUnicodeModelDir(t *testing.T) {
	// This is only useful for Windows with utf-16 characters, so skip this test for other platforms
	if runtime.GOOS != "windows" {
		t.Skip("Unicode test only applicable to windows")
	}
	// Only works for local testing
	if os.Getenv("AiAdmin_TEST_EXISTING") != "" {
		t.Skip("TestUnicodeModelDir only works for local testing, skipping")
	}

	modelDir, err := os.MkdirTemp("", "AiAdmin_埃")
	require.NoError(t, err)
	defer os.RemoveAll(modelDir)
	slog.Info("unicode", "AiAdmin_MODELS", modelDir)

	oldModelsDir := os.Getenv("AiAdmin_MODELS")
	if oldModelsDir == "" {
		defer os.Unsetenv("AiAdmin_MODELS")
	} else {
		defer os.Setenv("AiAdmin_MODELS", oldModelsDir)
	}
	err = os.Setenv("AiAdmin_MODELS", modelDir)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	req := api.GenerateRequest{
		Model:  "orca-mini",
		Prompt: "why is the sky blue?",
		Stream: &stream,
		Options: map[string]interface{}{
			"temperature": 0,
			"seed":        123,
		},
	}
	GenerateTestHelper(ctx, t, req, []string{"rayleigh", "scattering"})
}
