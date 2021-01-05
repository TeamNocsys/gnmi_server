package portchannel

import (
    "fmt"
    "github.com/spf13/cobra"
    "gnmi_server/cmd/command"
    "gnmi_server/cmd/command/config/portchannel/member"
    "gnmi_server/cmd/command/config/utils"
)

var (
    ErrInvaildPattern = fmt.Errorf("port channel pattern like %s ", utils.PORT_CHANNEL_PATTERN)
)

func NewPortChannelCommand(gnmiCli command.Client) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "portchannel",
        Short: "Configure portchannel",
        Args:  cobra.NoArgs,
    }

    cmd.AddCommand(
        NewAddCommand(gnmiCli),
        NewDelCommand(gnmiCli),
        member.NewMemberCommand(gnmiCli),
    )

    return cmd
}