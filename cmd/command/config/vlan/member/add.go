package member

import (
    "fmt"
    "github.com/spf13/cobra"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
    "gnmi_server/internal/pkg/swsssdk/helper"
    "strconv"
)

var (
    ErrMemberExists = fmt.Errorf("member already exists in vlan")
)

type addOptions struct {
    vid      int
    intf     string
    untagged bool
}

func NewAddCommand(gnmiCli command.Client) *cobra.Command {
    var opts addOptions

    cmd := &cobra.Command{
        Use:   "add <vid> <interface_name>",
        Short: "Add interface to vlan",
        Args:  cobra.ExactArgs(2),
        RunE: func(cmd *cobra.Command, args []string) error {
            var err error
            if opts.vid, err = strconv.Atoi(args[0]); err != nil {
                return err
            }
            opts.intf = args[1]
            return runAdd(gnmiCli, &opts)
        },
    }

    flags := cmd.Flags()
    flags.BoolVarP(&opts.untagged, "untagged", "u", true, "Indicates whether the port is untagged")

    return cmd
}

func runAdd(gnmiCli command.Client, opts *addOptions) error {
    if conn := gnmiCli.Config(); conn == nil {
        return swsssdk.ErrDatabaseNotExist
    } else {
        info, err := conn.GetEntry("VLAN", helper.VID(opts.vid))
        if err != nil {
            return err
        }
        if len(info) == 0 {
            return fmt.Errorf("%s doesn't exist", helper.VID(opts.vid))
        }

        var members = info["members"]
        if members != nil {
            for _, member := range members.([]string) {
                if member == opts.intf {
                    return ErrMemberExists
                }
            }
            members = append(members.([]string), opts.intf)
        } else {
            members = []string{opts.intf}
        }
        info["members"] = members
        conn.SetEntry("VLAN", helper.VID(opts.vid), info)
        var mode string
        if opts.untagged {
            mode = "untagged"
        } else {
            mode = "tagged"
        }
        conn.SetEntry("VLAN_MEMBER", []string{helper.VID(opts.vid), opts.intf}, map[string]interface{}{
            "tagging_mode": mode,
        })
        return nil
    }
}
