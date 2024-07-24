//go:build !windows && !darwin

package cmd

import (
	"context"
	"fmt"

	"github.com/ollama/ollama/api"
)

func startApp(ctx context.Context, client *api.Client) error {
	return fmt.Errorf("could not connect to AiAdmin server, run 'AiAdmin serve' to start it")
}
