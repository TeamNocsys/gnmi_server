package intf

import (
    "github.com/spf13/cobra"
    "gnmi_server/cmd/command"
    "gnmi_server/cmd/command/config/utils"
    "gnmi_server/internal/pkg/swsssdk"
    "regexp"
)

type downOptions struct {
    name string
}

func NewDownCommand(gnmiCli command.Client) *cobra.Command {
    var opts downOptions

    cmd := &cobra.Command{
        Use:   "down <name>",
        Short: "Disable switch port",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            opts.name = args[0]
            return runDown(gnmiCli, &opts)
        },
    }

    return cmd
}

func runDown(gnmiCli command.Client, opts *downOptions) error {
    if conn := gnmiCli.Config(); conn == nil {
        return swsssdk.ErrDatabaseNotExist
    } else {
        if ok, err := regexp.MatchString(utils.PORT_PATTERN, opts.name); err != nil {
            return err
        } else if ok {
            _, err := conn.SetEntry("PORT", opts.name, map[string]interface{}{
                "admin_status": "down",
            })
            return err
        }
        if ok, err := regexp.MatchString(utils.PORT_CHANNEL_PATTERN, opts.name); err != nil {
            return err
        } else if ok {
            _, err := conn.SetEntry("PORTCHANNEL", opts.name, map[string]interface{}{
                "admin_status": "down",
            })
            return err
        }
        if ok, err := regexp.MatchString(utils.VLAN_PATTERN, opts.name); err != nil {
            return err
        } else if ok {
            _, err := conn.SetEntry("VLAN", opts.name, map[string]interface{}{
                "admin_status": "down",
            })
            return err
        }

        return utils.ErrUnknowInterface
    }
}
