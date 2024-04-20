package cmd

import (
	"net/rpc"

	"github.com/spf13/cobra"

	"github.com/basecamp/parachute/internal/server"
)

type resumeCommand struct {
	cmd  *cobra.Command
	host string
}

func newResumeCommand() *resumeCommand {
	resumeCommand := &resumeCommand{}
	resumeCommand.cmd = &cobra.Command{
		Use:   "resume",
		Short: "Resume a service",
		RunE:  resumeCommand.run,
		Args:  cobra.NoArgs,
	}

	resumeCommand.cmd.Flags().StringVar(&resumeCommand.host, "host", "", "Host to resume (empty for wildcard)")

	return resumeCommand
}

func (c *resumeCommand) run(cmd *cobra.Command, args []string) error {
	return withRPCClient(globalConfig.SocketPath(), func(client *rpc.Client) error {
		var response bool
		args := server.ResumeArgs{
			Host: c.host,
		}

		err := client.Call("parachute.Resume", args, &response)

		return err
	})
}