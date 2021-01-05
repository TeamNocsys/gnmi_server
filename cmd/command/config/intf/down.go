package intf

import (
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/spf13/cobra"
    "gnmi_server/cmd/command"
    "gnmi_server/cmd/command/config/utils"
    "gnmi_server/internal/pkg/swsssdk"
    "gnmi_server/internal/pkg/swsssdk/helper/config_db"
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
            _, err := conn.SetEntry(config_db.PORT_TABLE, opts.name, map[string]interface{}{
                config_db.PORT_ADMIN_STATUS: config_db.AdminStatusToString(sonicpb.SonicPortAdminStatus_SONICPORTADMINSTATUS_down),
            })
            return err
        }
        if ok, err := regexp.MatchString(utils.PORT_CHANNEL_PATTERN, opts.name); err != nil {
            return err
        } else if ok {
            _, err := conn.SetEntry(config_db.PORTCHANNEL_TABLE, opts.name, map[string]interface{}{
                config_db.PORTCHANNEL_ADMIN_STATUS: config_db.AdminStatusToString(sonicpb.SonicPortchannelAdminStatus_SONICPORTCHANNELADMINSTATUS_down),
            })
            return err
        }
        if ok, err := regexp.MatchString(utils.VLAN_PATTERN, opts.name); err != nil {
            return err
        } else if ok {
            _, err := conn.SetEntry(config_db.VLAN_TABLE, opts.name, map[string]interface{}{
                config_db.VLAN_ADMIN_STATUS: config_db.AdminStatusToString(sonicpb.SonicVlanAdminStatus_SONICVLANADMINSTATUS_down),
            })
            return err
        }

        return utils.ErrUnknowInterface
    }
}
