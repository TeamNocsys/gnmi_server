package get

import (
    my_gnmi "gnmi_server/pkg/gnmi"
)

func GetServeMux() *my_gnmi.GetServeMux {
    mux := my_gnmi.NewGetServeMux()
    route(mux)
    return mux
}

func route(mux *my_gnmi.GetServeMux) {
    mux.AddRouter("/test", Test).
        AddRouter("/sonic-platform/platform", ComponentInfoHandler).
        AddRouter("/sonic-platform/platform/component-list/fan", FanInfoHandler).
        AddRouter("/sonic-platform/platform/component-list/power-supply", PowerSupplyInfoHandler).
        AddRouter("/sonic-platform/platform/component-list/temperature", TemperatureInfoHandler).
        AddRouter("/sonic-lldp/lldp", LLDPHandler).
        AddRouter("/sonic-port/port/port-state-list", PortStateHandler).
        AddRouter("/sonic-port/port/port-state-list/counters", PortStateHandler)
}
