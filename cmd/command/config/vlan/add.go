package vlan

import (
    "github.com/spf13/cobra"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
    "gnmi_server/internal/pkg/swsssdk/helper"
    "strconv"
)

type addOptions struct {
    vid int
}

func NewAddCommand(gnmiCli command.Client) *cobra.Command {
    var opts addOptions

    cmd := &cobra.Command{
        Use:   "add <vid>",
        Short: "Add vlan to the switch",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            var err error
            if opts.vid, err = strconv.Atoi(args[0]); err != nil {
                return err
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
        _, err := conn.SetEntry(helper.VLAN_TABLE_NAME, helper.VID(opts.vid), map[string]interface{}{
            "vlanid": opts.vid,
        })
        return err
    }
}
