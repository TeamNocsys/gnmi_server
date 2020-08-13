package vlan

import (
    "github.com/spf13/cobra"
    "gnmi_server/cmd/command"
    "gnmi_server/cmd/command/config/vlan/member"
)

func NewVLANCommand(gnmiCli command.Client) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "vlan",
        Short: "Configure vlan",
        Args:  cobra.NoArgs,
    }

    cmd.AddCommand(
        NewAddCommand(gnmiCli),
        NewDelCommand(gnmiCli),
        member.NewMemberCommand(gnmiCli),
    )

    return cmd
}
