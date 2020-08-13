package interfaces

import (
    "github.com/spf13/cobra"
    "gnmi_server/cmd/command"
)

func NewInterfacesCommand(gnmiCli command.Client) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "interfaces",
        Short: "Show details of the network interfaces",
        Args:  cobra.NoArgs,
    }

    cmd.AddCommand(
        NewCountersCommand(gnmiCli),
    )

    return cmd
}
