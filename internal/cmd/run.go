package cmd

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/kevinmcconnell/mproxy/internal/server"
)

type runCommand struct {
	cmd              *cobra.Command
	config           server.Config
	debugLogsEnabled bool
}

func newRunCommand() *runCommand {
	runCommand := &runCommand{}
	runCommand.cmd = &cobra.Command{
		Use:   "run",
		Short: "Run the server",
		RunE:  runCommand.runServer,
	}

	runCommand.cmd.Flags().BoolVar(&runCommand.debugLogsEnabled, "debug", false, "Include debugging logs")
	runCommand.cmd.Flags().BoolVar(&runCommand.config.ACMEUseStaging, "ssl-staging", false, "Use Let's Encrypt staging environment for SSL certificates")
	runCommand.cmd.Flags().IntVar(&runCommand.config.HttpPort, "http-port", server.DefaultHttpPort, "Port to serve HTTP traffic on")
	runCommand.cmd.Flags().IntVar(&runCommand.config.HttpsPort, "https-port", server.DefaultHttpsPort, "Port to serve HTTPS traffic on")
	runCommand.cmd.Flags().Int64Var(&runCommand.config.MaxRequestBodySize, "max-request-body", 0, "Max size of request body (default of 0 means unlimited)")
	runCommand.cmd.Flags().DurationVar(&runCommand.config.HttpIdleTimeout, "http-idle-timeout", server.DefaultHttpIdleTimeout, "Timeout before idle connection is closed")
	runCommand.cmd.Flags().DurationVar(&runCommand.config.HttpReadTimeout, "http-read-timeout", server.DefaultHttpReadTimeout, "Tiemout for client to send a request")
	runCommand.cmd.Flags().DurationVar(&runCommand.config.HttpWriteTimeout, "http-write-timeout", server.DefaultHttpWriteTimeout, "Timeout for client to receive a response")

	return runCommand
}

func (c *runCommand) runServer(cmd *cobra.Command, args []string) error {
	c.setLogger()
	c.config.ConfigDir = configDir

	s := server.NewServer(&c.config, server.NewRouter())
	s.Start()
	defer s.Stop()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch

	return nil
}

func (c *runCommand) setLogger() {
	level := slog.LevelInfo
	if c.debugLogsEnabled {
		level = slog.LevelDebug
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})))
}