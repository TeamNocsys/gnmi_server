package interfaces

import (
    "fmt"
    "github.com/spf13/cobra"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
)

type countersOptions struct {
}

func NewCountersCommand(gnmiCli command.Client) *cobra.Command {
    var opts countersOptions

    cmd := &cobra.Command{
        Use:   "counters",
        Short: "Show interface counters",
        Args:  cobra.NoArgs,
        RunE: func(cmd *cobra.Command, args []string) error {
            return runCounters(gnmiCli, &opts)
        },
    }

    return cmd
}

func runCounters(gnmiCli command.Client, opts *countersOptions) error {
    if conn := gnmiCli.State(); conn == nil {
        return swsssdk.ErrDatabaseNotExist
    } else {
        content, err := conn.GetAll(swsssdk.COUNTERS_DB, "COUNTERS_PORT_NAME_MAP")
        if err != nil {
            return err
        }
        fmt.Println(content)
        return nil
    }
}
