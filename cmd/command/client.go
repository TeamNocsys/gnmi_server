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

func NewGnmiClient() *GnmiClient {
    cli := &GnmiClient{
        config: swsssdk.NewConfigDBConnector(),
        state: swsssdk.NewConnector(),
    }
    return cli
}

func (gc *GnmiClient) Connect() error {
    if !gc.config.Connect() {
        return fmt.Errorf("Connect to database %s failed!", swsssdk.CONFIG_DB)
    }

    if !gc.state.Connect(swsssdk.APPL_DB) {
        return fmt.Errorf("Connect to database %s failed!", swsssdk.APPL_DB)
    }
    if !gc.state.Connect(swsssdk.ASIC_DB) {
        return fmt.Errorf("Connect to database %s failed!", swsssdk.ASIC_DB)
    }
    if !gc.state.Connect(swsssdk.COUNTERS_DB) {
        return fmt.Errorf("Connect to database %s failed!", swsssdk.COUNTERS_DB)
    }
    if !gc.state.Connect(swsssdk.LOGLEVEL_DB) {
        return fmt.Errorf("Connect to database %s failed!", swsssdk.LOGLEVEL_DB)
    }
    if !gc.state.Connect(swsssdk.PFC_WD_DB) {
        return fmt.Errorf("Connect to database %s failed!", swsssdk.PFC_WD_DB)
    }
    if !gc.state.Connect(swsssdk.FLEX_COUNTER_DB) {
        return fmt.Errorf("Connect to database %s failed!", swsssdk.FLEX_COUNTER_DB)
    }
    if !gc.state.Connect(swsssdk.STATE_DB) {
        return fmt.Errorf("Connect to database %s failed!", swsssdk.STATE_DB)
    }
    if !gc.state.Connect(swsssdk.SNMP_OVERLAY_DB) {
        return fmt.Errorf("Connect to database %s failed!", swsssdk.SNMP_OVERLAY_DB)
    }
    return nil
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
