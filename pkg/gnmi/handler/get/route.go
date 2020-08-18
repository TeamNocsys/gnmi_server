package get

import (
    "gnmi_server/pkg/gnmi"
)

func GetServeMux() *gnmi.GetServeMux {
    mux := gnmi.NewGetServeMux()
    route(mux)
    return mux
}

func route(mux *gnmi.GetServeMux) {
    mux.AddRouter("/test", Test).
        AddRouter("/sonic-platform/platform", ComponentInfoHandler).
        AddRouter("/sonic-platform/platform/component-list/fan", FanInfoHandler).
        AddRouter("/sonic-platform/platform/component-list/power-supply", PowerSupplyInfoHandler).
        AddRouter("/sonic-platform/platform/component-list/state/temperature", TemperatureInfoHandler).
        AddRouter("/sonic-platform/platform/component-list/system", SystemInfoHandler).
        AddRouter("/sonic-lldp/lldp", LLDPHandler).
        AddRouter("/sonic-port/port/port-state-list", PortStateHandler).
        AddRouter("/sonic-port/port/port-state-list/counters", PortStateHandler)
}
