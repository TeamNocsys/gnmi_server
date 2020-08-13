package member

import (
    "github.com/spf13/cobra"
    "gnmi_server/cmd/command"
)

func NewMemberCommand(gnmiCli command.Client) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "member",
        Short: "Configure vlan member",
        Args:  cobra.NoArgs,
    }

    cmd.AddCommand(
        NewAddCommand(gnmiCli),
        NewDelCommand(gnmiCli),
    )

    return cmd
}
