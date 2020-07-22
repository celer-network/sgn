package main

import (
	"github.com/celer-network/sgn/cmd"
	"github.com/spf13/cobra"
)

func main() {
	cobra.EnableCommandSorting = false

	executor := cmd.GetSgnopsExecutor()
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}
