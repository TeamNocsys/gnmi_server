package member

import (
    "fmt"
    "github.com/spf13/cobra"
    "gnmi_server/cmd/command"
    "gnmi_server/cmd/command/config/utils"
)

var (
    ErrMemberExists = fmt.Errorf("member already exists")
    ErrInvaildPattern = fmt.Errorf("interface pattern like %s or %s", utils.PORT_PATTERN, utils.PORT_CHANNEL_PATTERN)
)

func NewMemberCommand(gnmiCli command.Client) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "member",
        Short: "Configure port channel member",
        Args:  cobra.NoArgs,
    }

    cmd.AddCommand(
        NewAddCommand(gnmiCli),
        NewDelCommand(gnmiCli),
    )

    return cmd
}
