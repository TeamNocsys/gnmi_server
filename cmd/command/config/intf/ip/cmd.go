package ip

import (
    "github.com/spf13/cobra"
    "gnmi_server/cmd/command"
)

func NewIpCommand(gnmiCli command.Client) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "ip",
        Short: "Configure interface ip address",
        Args:  cobra.NoArgs,
    }

    cmd.AddCommand(
        NewAddCommand(gnmiCli),
        NewDelCommand(gnmiCli),
    )

    return cmd
}
