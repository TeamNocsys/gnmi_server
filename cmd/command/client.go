package command

import (
    "fmt"
    "gnmi_server/internal/pkg/swsssdk"
)

type Client interface {
    Close()
    Connector() *swsssdk.ConfigDBConnector
}

type GnmiClient struct {
    connector *swsssdk.ConfigDBConnector
}

func NewGnmiClient() (*GnmiClient, error) {
    cli := &GnmiClient{
        connector: swsssdk.NewConfigDBConnector(),
    }

    if !cli.connector.Connect() {
        return nil, fmt.Errorf("Connect to database %s failed!", swsssdk.CONFIG_DB_NAME)
    }

    if !cli.connector.Connector.Connect(swsssdk.APPLICATION_DB_NAME) {
        return nil, fmt.Errorf("Connect to database %s failed!", swsssdk.APPLICATION_DB_NAME)
    }

    return cli, nil
}

func (gc *GnmiClient) Connector() *swsssdk.ConfigDBConnector {
    return gc.connector
}

func (gc *GnmiClient) Close() {
    gc.connector.Close()
}
