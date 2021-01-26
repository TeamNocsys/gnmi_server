package member

import (
    "fmt"
    "github.com/spf13/cobra"
    "gnmi_server/cmd/command"
    "gnmi_server/cmd/command/config/utils"
    "gnmi_server/internal/pkg/swsssdk"
    "regexp"
)

type delOptions struct {
    name     string
    port     string
}

func NewDelCommand(gnmiCli command.Client) *cobra.Command {
    var opts delOptions

    cmd := &cobra.Command{
        Use:   "del <name> <port_name>",
        Short: "Remove port from port channel",
        Args:  cobra.ExactArgs(2),
        RunE: func(cmd *cobra.Command, args []string) error {
            opts.name = args[0]
            if ok, err := regexp.MatchString(utils.PORT_CHANNEL_PATTERN, opts.name); err != nil {
                return err
            } else if !ok {
                return ErrInvaildPattern
            }
            opts.port = args[1]
            if ok, err := regexp.MatchString(utils.PORT_PATTERN, opts.name); err != nil {
                return err
            } else if !ok {
                return ErrInvaildPattern
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
        info, err := conn.GetEntry("PORTCHANNEL", opts.name)
        if err != nil {
            return err
        }
        if len(info) == 0 {
            return fmt.Errorf("%s doesn't exist", opts.name)
        }

        ok, err := conn.HasEntry("PORT", opts.port)
        if err != nil {
            return err
        }
        if !ok {
            return nil
        }

        var members = info["members"]
        if members != nil {
            for index, member := range members.([]string) {
                if member == opts.port {
                    members = append(members.([]string)[:index], members.([]string)[index+1:]...)
                }
            }
            members = append(members.([]string), opts.port)
            if len(members.([]string)) > 0 {
                info["members"] = members
            } else {
                delete(info, "members")
            }
            if _, err := conn.SetEntry("PORTCHANNEL", opts.name, info); err != nil {
                return err
            }
        }
        return nil
    }
}
