package ip

import (
    "github.com/spf13/cobra"
    "gnmi_server/cmd/command"
    "gnmi_server/cmd/command/config/utils"
    "gnmi_server/internal/pkg/swsssdk"
    "gnmi_server/internal/pkg/swsssdk/helper/config_db"
    "regexp"
)

type delOptions struct {
    name     string
    addr     string
}

func NewDelCommand(gnmiCli command.Client) *cobra.Command {
    var opts delOptions

    cmd := &cobra.Command{
        Use:   "del <name> <address>",
        Short: "Remove ip address from interface",
        Args:  cobra.ExactArgs(2),
        RunE: func(cmd *cobra.Command, args []string) error {
            opts.name = args[0]
            opts.addr = args[1]
            return runDel(gnmiCli, &opts)
        },
    }

    return cmd
}

func runDel(gnmiCli command.Client, opts *delOptions) error {
    if conn := gnmiCli.Config(); conn == nil {
        return swsssdk.ErrDatabaseNotExist
    } else {
        if ok, err := regexp.MatchString(utils.PORT_PATTERN, opts.name); err != nil {
            return err
        } else if ok {
            _, err := conn.SetEntry(config_db.PORT_TABLE, []string{opts.name, opts.addr}, nil)
            return err
        }
        if ok, err := regexp.MatchString(utils.PORT_CHANNEL_PATTERN, opts.name); err != nil {
            return err
        } else if ok {
            _, err := conn.SetEntry(config_db.PORTCHANNEL_TABLE, []string{opts.name, opts.addr}, nil)
            return err
        }
        if ok, err := regexp.MatchString(utils.VLAN_PATTERN, opts.name); err != nil {
            return err
        } else if ok {
            _, err := conn.SetEntry(config_db.VLAN_TABLE, []string{opts.name, opts.addr}, nil)
            return err
        }
        return utils.ErrUnknowInterface
    }
}
