package member

import (
    "fmt"
    "github.com/spf13/cobra"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
    "gnmi_server/internal/pkg/swsssdk/helper"
    "strconv"
)

type delOptions struct {
    vid  int
    intf string
}

func NewDelCommand(gnmiCli command.Client) *cobra.Command {
    var opts delOptions

    cmd := &cobra.Command{
        Use:   "del <vid> <interface_name>",
        Short: "Remove interface from vlan",
        Args:  cobra.ExactArgs(2),
        RunE: func(cmd *cobra.Command, args []string) error {
            var err error
            if opts.vid, err = strconv.Atoi(args[0]); err != nil {
                return err
            }
            opts.intf = args[1]
            return runDel(gnmiCli, &opts)
        },
    }

    return cmd
}

func runDel(gnmiCli command.Client, opts *delOptions) error {
    if conn := gnmiCli.Connector(); conn == nil {
        return swsssdk.ErrDatabaseNotExist
    } else {
        info, err := conn.GetEntry(helper.VLAN_TABLE_NAME, helper.VID(opts.vid))
        if err != nil {
            return err
        }
        if len(info) == 0 {
            return fmt.Errorf("%s doesn't exist", helper.VID(opts.vid))
        }
        var members = info["members"]
        if members != nil {
            for index, member := range members.([]string) {
                if member == opts.intf {
                    members = append(members.([]string)[:index], members.([]string)[index+1:]...)
                }
            }
            members = append(members.([]string), opts.intf)
            if len(members.([]string)) > 0 {
                info["members"] = members
            } else {
                delete(info, "members")
            }
            if _, err := conn.SetEntry(helper.VLAN_TABLE_NAME, helper.VID(opts.vid), info); err != nil {
                return err
            }
        }
        if _, err := conn.SetEntry(helper.VLAN_MEMBERTABLE_NAME, []string{helper.VID(opts.vid), opts.intf}, nil); err != nil {
            return err
        }
        return nil
    }
}
