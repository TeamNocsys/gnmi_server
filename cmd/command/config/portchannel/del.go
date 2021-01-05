package portchannel

import (
    "github.com/spf13/cobra"
    "gnmi_server/cmd/command"
    "gnmi_server/cmd/command/config/utils"
    "gnmi_server/internal/pkg/swsssdk"
    "gnmi_server/internal/pkg/swsssdk/helper/config_db"
    "regexp"
)

type delOptions struct {
    name string
}

func NewDelCommand(gnmiCli command.Client) *cobra.Command {
    var opts delOptions

    cmd := &cobra.Command{
        Use:   "del <name>",
        Short: "Remove port channel from the switch",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            opts.name = args[0]
            return runDel(gnmiCli, &opts)
        },
    }

    return cmd
}

func runDel(gnmiCli command.Client, opts *delOptions) error {
    if conn := gnmiCli.Config(); conn == nil {
        return swsssdk.ErrDatabaseNotExist
    } else {
        if ok, err := regexp.MatchString(utils.PORT_CHANNEL_PATTERN, opts.name); err != nil {
            return err
        } else if ok {
            _, err := conn.SetEntry(config_db.PORTCHANNEL_TABLE, opts.name, nil)
            return err
        }
        return ErrInvaildPattern
    }
}