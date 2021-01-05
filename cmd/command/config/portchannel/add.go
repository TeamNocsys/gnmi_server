package portchannel

import (
    "github.com/spf13/cobra"
    "gnmi_server/cmd/command"
    "gnmi_server/cmd/command/config/utils"
    "gnmi_server/internal/pkg/swsssdk"
    "gnmi_server/internal/pkg/swsssdk/helper/config_db"
    "regexp"
)

type addOptions struct {
    name string
}

func NewAddCommand(gnmiCli command.Client) *cobra.Command {
    var opts addOptions

    cmd := &cobra.Command{
        Use:   "add <name>",
        Short: "Add port channel to the switch",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            opts.name = args[0]
            return runAdd(gnmiCli, &opts)
        },
    }

    return cmd
}

func runAdd(gnmiCli command.Client, opts *addOptions) error {
    if conn := gnmiCli.Config(); conn == nil {
        return swsssdk.ErrDatabaseNotExist
    } else {
        if ok, err := regexp.MatchString(utils.PORT_CHANNEL_PATTERN, opts.name); err != nil {
            return err
        } else if ok {
            _, err := conn.SetEntry(config_db.PORTCHANNEL_TABLE, opts.name, map[string]interface{}{})
            return err
        }
        return ErrInvaildPattern
    }
}