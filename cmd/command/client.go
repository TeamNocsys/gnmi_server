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
        return nil, fmt.Errorf("Connect to database %s failed!", swsssdk.CONFIG_DB)
    }

    if !cli.state.Connect(swsssdk.APPL_DB) {
        return nil, fmt.Errorf("Connect to database %s failed!", swsssdk.APPL_DB)
    }
    if !cli.state.Connect(swsssdk.ASIC_DB) {
        return nil, fmt.Errorf("Connect to database %s failed!", swsssdk.ASIC_DB)
    }
    if !cli.state.Connect(swsssdk.COUNTERS_DB) {
        return nil, fmt.Errorf("Connect to database %s failed!", swsssdk.COUNTERS_DB)
    }
    if !cli.state.Connect(swsssdk.LOGLEVEL_DB) {
        return nil, fmt.Errorf("Connect to database %s failed!", swsssdk.LOGLEVEL_DB)
    }
    if !cli.state.Connect(swsssdk.PFC_WD_DB) {
        return nil, fmt.Errorf("Connect to database %s failed!", swsssdk.PFC_WD_DB)
    }
    if !cli.state.Connect(swsssdk.FLEX_COUNTER_DB) {
        return nil, fmt.Errorf("Connect to database %s failed!", swsssdk.FLEX_COUNTER_DB)
    }
    if !cli.state.Connect(swsssdk.STATE_DB) {
        return nil, fmt.Errorf("Connect to database %s failed!", swsssdk.STATE_DB)
    }
    if !cli.state.Connect(swsssdk.SNMP_OVERLAY_DB) {
        return nil, fmt.Errorf("Connect to database %s failed!", swsssdk.SNMP_OVERLAY_DB)
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
