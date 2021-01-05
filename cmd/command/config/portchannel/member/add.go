package member

import (
    "fmt"
    "github.com/spf13/cobra"
    "gnmi_server/cmd/command"
    "gnmi_server/cmd/command/config/utils"
    "gnmi_server/internal/pkg/swsssdk"
    "gnmi_server/internal/pkg/swsssdk/helper/config_db"
    "regexp"
)

type addOptions struct {
    name     string
    port     string
}

func NewAddCommand(gnmiCli command.Client) *cobra.Command {
    var opts addOptions

    cmd := &cobra.Command{
        Use:   "add <name> <port_name>",
        Short: "Add port to port channel",
        Args:  cobra.ExactArgs(2),
        RunE: func(cmd *cobra.Command, args []string) error {
            opts.name = args[0]
            if ok, err := regexp.MatchString(utils.PORT_CHANNEL_PATTERN, opts.name); err != nil {
                return err
            } else if !ok {
                return ErrInvaildPattern
            }
            opts.port = args[1]
            if ok, err := regexp.MatchString(utils.PORT_PATTERN, opts.port); err != nil {
                return err
            } else if !ok {
                return ErrInvaildPattern
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
        info, err := conn.GetEntry(config_db.PORTCHANNEL_TABLE, opts.name)
        if err != nil {
            return err
        }
        if len(info) == 0 {
            return fmt.Errorf("%s doesn't exist", opts.name)
        }

        ok, err := conn.HasEntry(config_db.PORT_TABLE, opts.port)
        if err != nil {
            return err
        }
        if !ok {
            return nil
        }

        var members = info["members"]
        if members != nil {
            for _, member := range members.([]string) {
                if member == opts.port {
                    return ErrMemberExists
                }
            }
            members = append(members.([]string), opts.port)
        } else {
            members = []string{opts.port}
        }
        info["members"] = members
        conn.SetEntry(config_db.PORTCHANNEL_TABLE, opts.name, info)
        return nil
    }
}
