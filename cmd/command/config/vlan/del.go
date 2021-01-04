package vlan

import (
    "github.com/spf13/cobra"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
    "gnmi_server/internal/pkg/swsssdk/helper"
    "gnmi_server/internal/pkg/swsssdk/helper/config_db"
    "strconv"
)

type delOptions struct {
    vid int
}

func NewDelCommand(gnmiCli command.Client) *cobra.Command {
    var opts delOptions

    cmd := &cobra.Command{
        Use:   "del <vid>",
        Short: "Remove vlan from the switch",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            var err error
            if opts.vid, err = strconv.Atoi(args[0]); err != nil {
                return err
            }
            return runDel(gnmiCli, &opts)
        },
    }

    return cmd
}

func runDel(gnmiCli command.Client, opts *delOptions) error {
    if conn := gnmiCli.Config(); conn == nil {
        return swsssdk.ErrDatabaseNotExist
    } else {
        _, err := conn.SetEntry(config_db.VLAN_TABLE, helper.VID(opts.vid), nil)
        return err
    }
}
