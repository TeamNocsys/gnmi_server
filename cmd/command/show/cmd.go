package show

import (
    "github.com/spf13/cobra"
    "gnmi_server/cmd/command"
    "gnmi_server/cmd/command/show/interfaces"
)

func NewShowCommand(gnmiCli command.Client) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "show",
        Short: "Show states",
        Args:  cobra.NoArgs,
    }

    cmd.AddCommand(
        interfaces.NewInterfacesCommand(gnmiCli),
    )

    return cmd
}
