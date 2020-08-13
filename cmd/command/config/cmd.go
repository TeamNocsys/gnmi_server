package config

import (
    "github.com/spf13/cobra"
    "gnmi_server/cmd/command"
    "gnmi_server/cmd/command/config/vlan"
)

func NewConfigCommand(gnmiCli command.Client) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "config",
        Short: "Manage configs",
        Args:  cobra.NoArgs,
    }

    cmd.AddCommand(
        vlan.NewVLANCommand(gnmiCli),
    )

    return cmd
}
