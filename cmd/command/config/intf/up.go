package intf

import (
    "github.com/spf13/cobra"
    "gnmi_server/cmd/command"
    "gnmi_server/cmd/command/config/utils"
    "gnmi_server/internal/pkg/swsssdk"
    "regexp"
)

type upOptions struct {
    name string
}

func NewUpCommand(gnmiCli command.Client) *cobra.Command {
    var opts upOptions

    cmd := &cobra.Command{
        Use:   "up <name>",
        Short: "Enable switch port",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            opts.name = args[0]
            return runUp(gnmiCli, &opts)
        },
    }

    return cmd
}

func runUp(gnmiCli command.Client, opts *upOptions) error {
    if conn := gnmiCli.Config(); conn == nil {
        return swsssdk.ErrDatabaseNotExist
    } else {
        if ok, err := regexp.MatchString(utils.PORT_PATTERN, opts.name); err != nil {
            return err
        } else if ok {
            _, err := conn.SetEntry("PORT", opts.name, map[string]interface{}{
                "admin_status": "up",
            })
            return err
        }
        if ok, err := regexp.MatchString(utils.PORT_CHANNEL_PATTERN, opts.name); err != nil {
            return err
        } else if ok {
            _, err := conn.SetEntry("PORTCHANNEL", opts.name, map[string]interface{}{
                "admin_status": "up",
            })
            return err
        }
        if ok, err := regexp.MatchString(utils.VLAN_PATTERN, opts.name); err != nil {
            return err
        } else if ok {
            _, err := conn.SetEntry("VLAN", opts.name, map[string]interface{}{
                "admin_status": "up",
            })
            return err
        }
        return utils.ErrUnknowInterface
    }
}
