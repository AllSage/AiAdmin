package main

import (
	"context"

	"github.com/AllSage/AiAdmin/cmd"
	"github.com/spf13/cobra"
)

func main() {
	cobra.CheckErr(cmd.NewCLI().ExecuteContext(context.Background()))
}
