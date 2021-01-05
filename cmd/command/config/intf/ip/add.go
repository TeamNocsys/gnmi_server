package ip

import (
    "github.com/spf13/cobra"
    "gnmi_server/cmd/command"
    "gnmi_server/cmd/command/config/utils"
    "gnmi_server/internal/pkg/swsssdk"
    "gnmi_server/internal/pkg/swsssdk/helper/config_db"
    "regexp"
)

type addOptions struct {
    name     string
    addr     string
    gw       string
}

func NewAddCommand(gnmiCli command.Client) *cobra.Command {
    var opts addOptions

    cmd := &cobra.Command{
        Use:   "add <name> <address> <gateway>",
        Short: "Add ip address to interface",
        Args:  cobra.RangeArgs(2,3),
        RunE: func(cmd *cobra.Command, args []string) error {
            opts.name = args[0]
            opts.addr = args[1]
            if len(args) == 3 {
                opts.gw = args[2]
            }
            return runAdd(gnmiCli, &opts)
        },
    }

    return cmd
}

func runAdd(gnmiCli command.Client, opts *addOptions) error {
    if conn := gnmiCli.Config(); conn == nil {
        return swsssdk.ErrDatabaseNotExist
    } else {
        if ok, err := regexp.MatchString(utils.PORT_PATTERN, opts.name); err != nil {
            return err
        } else if ok {
            _, err := conn.SetEntry(config_db.PORT_TABLE, []string{opts.name, opts.addr}, map[string]interface{}{})
            return err
        }
        if ok, err := regexp.MatchString(utils.PORT_CHANNEL_PATTERN, opts.name); err != nil {
            return err
        } else if ok {
            _, err := conn.SetEntry(config_db.PORTCHANNEL_TABLE, []string{opts.name, opts.addr}, map[string]interface{}{})
            return err
        }
        if ok, err := regexp.MatchString(utils.VLAN_PATTERN, opts.name); err != nil {
            return err
        } else if ok {
            _, err := conn.SetEntry(config_db.VLAN_TABLE, []string{opts.name, opts.addr}, map[string]interface{}{})
            return err
        }
        return utils.ErrUnknowInterface
    }
}
