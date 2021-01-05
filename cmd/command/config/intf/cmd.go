package intf

import (
    "github.com/spf13/cobra"
    "gnmi_server/cmd/command"
    "gnmi_server/cmd/command/config/intf/ip"
)

func NewInterfaceCommand(gnmiCli command.Client) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "interface",
        Short: "Configure interface",
        Args:  cobra.NoArgs,
    }

    cmd.AddCommand(
        NewUpCommand(gnmiCli),
        NewDownCommand(gnmiCli),
        ip.NewIpCommand(gnmiCli),
    )

    return cmd
}