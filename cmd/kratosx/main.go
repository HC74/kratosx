package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/HC74/kratos/cmd/kratos/v2/internal/change"
	"github.com/HC74/kratos/cmd/kratos/v2/internal/project"
	"github.com/HC74/kratos/cmd/kratos/v2/internal/proto"
	"github.com/HC74/kratos/cmd/kratos/v2/internal/run"
	"github.com/HC74/kratos/cmd/kratos/v2/internal/upgrade"
)

// release is the current kratosx tool version.
const release = "v1.0.0"

var rootCmd = &cobra.Command{
	Use:     "kratosx",
	Short:   "Kratosx: An elegant toolkit for Go microservices.",
	Long:    `Kratosx: An elegant toolkit for Go microservices.`,
	Version: release,
}

func init() {
	rootCmd.AddCommand(project.CmdNew)
	rootCmd.AddCommand(proto.CmdProto)
	rootCmd.AddCommand(upgrade.CmdUpgrade)
	rootCmd.AddCommand(change.CmdChange)
	rootCmd.AddCommand(run.CmdRun)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
