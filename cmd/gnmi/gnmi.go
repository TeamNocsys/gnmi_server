package main

import (
    "fmt"
    "github.com/spf13/cobra"
    "gnmi_server/cmd/command"
    "gnmi_server/cmd/command/config"
    "gnmi_server/cmd/command/run"
    "gnmi_server/cmd/command/show"
    "gnmi_server/internal/pkg/swsssdk"
    "os"
)

func exec(gnmiCli *command.GnmiClient) error {
    var cfg string
    cmd := &cobra.Command{
        Use:              "gnmi [OPTIONS] COMMAND [ARG...]",
        Short:            "A implementation of gnmi service for sonic",
        Args:             cobra.NoArgs,
        SilenceUsage:     true,
        SilenceErrors:    true,
        TraverseChildren: true,
        PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
            if cfg != "" {
                swsssdk.LoadConfig(cfg)
            }
            return gnmiCli.Connect()
        },
    }

    cmd.PersistentFlags().StringVar(&cfg, "config", "", "sonic database configuration file path")

    cmd.AddCommand(
        config.NewConfigCommand(gnmiCli),
        run.NewRunCommand(gnmiCli),
        show.NewShowCommand(gnmiCli),
    )

    return cmd.Execute()
}

func main() {
    gnmiCli:= command.NewGnmiClient()
    defer gnmiCli.Close()

    if err := exec(gnmiCli); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
