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
            _, err := conn.SetEntry(config_db.PORT_TABLE, opts.name, map[string]interface{}{
                config_db.PORT_ADMIN_STATUS: config_db.AdminStatusToString(sonicpb.SonicPortAdminStatus_SONICPORTADMINSTATUS_up),
            })
            return err
        }
        if ok, err := regexp.MatchString(utils.PORT_CHANNEL_PATTERN, opts.name); err != nil {
            return err
        } else if ok {
            _, err := conn.SetEntry(config_db.PORTCHANNEL_TABLE, opts.name, map[string]interface{}{
                config_db.PORTCHANNEL_ADMIN_STATUS: config_db.AdminStatusToString(sonicpb.SonicPortchannelAdminStatus_SONICPORTCHANNELADMINSTATUS_up),
            })
            return err
        }
        if ok, err := regexp.MatchString(utils.VLAN_PATTERN, opts.name); err != nil {
            return err
        } else if ok {
            _, err := conn.SetEntry(config_db.VLAN_TABLE, opts.name, map[string]interface{}{
                config_db.VLAN_ADMIN_STATUS: config_db.AdminStatusToString(sonicpb.SonicVlanAdminStatus_SONICVLANADMINSTATUS_up),
            })
            return err
        }
        return utils.ErrUnknowInterface
    }
}
