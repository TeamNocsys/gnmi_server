package command

import (
    "fmt"
    "gnmi_server/internal/pkg/swsssdk"
)

type Client interface {
    Close()
    Config() *swsssdk.ConfigDBConnector
    State() *swsssdk.Connector
}

type GnmiClient struct {
    config *swsssdk.ConfigDBConnector
    state *swsssdk.Connector
}

func NewGnmiClient() (*GnmiClient, error) {
    cli := &GnmiClient{
        config: swsssdk.NewConfigDBConnector(),
        state: swsssdk.NewConnector(),
    }

    if !cli.config.Connect() {
        return nil, fmt.Errorf("Connect to database %s failed!", swsssdk.CONFIG_DB_NAME)
    }

    if !cli.state.Connect(swsssdk.APPLICATION_DB_NAME) {
        return nil, fmt.Errorf("Connect to database %s failed!", swsssdk.APPLICATION_DB_NAME)
    }

    return cli, nil
}

func (gc *GnmiClient) Config() *swsssdk.ConfigDBConnector {
    return gc.config
}

func (gc *GnmiClient) State() *swsssdk.Connector {
    return gc.state
}

func (gc *GnmiClient) Close() {
    gc.config.Close()
    gc.state.Close()
}
