package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/AllSage/AiAdmin/api"
)

func startApp(ctx context.Context, client *api.Client) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	link, err := os.Readlink(exe)
	if err != nil {
		return err
	}
	if !strings.Contains(link, "AiAdmin.app") {
		return fmt.Errorf("could not find AiAdmin app")
	}
	path := strings.Split(link, "AiAdmin.app")
	if err := exec.Command("/usr/bin/open", "-a", path[0]+"AiAdmin.app").Run(); err != nil {
		return err
	}
	return waitForServer(ctx, client)
}
