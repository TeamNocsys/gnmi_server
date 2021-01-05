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
        AddRouter("/sonic-port/port/port-list", PortHandler).
        AddRouter("/sonic-port/port/port-state-list", PortStateHandler).
        AddRouter("/sonic-port/port/port-state-list/counters", PortStateHandler).
        AddRouter("/sonic-portchannel/portchannel/portchannel-list", PortChannelHandler).
        AddRouter("/sonic-vlan/vlan/vlan-list", VlanHandler).
        AddRouter("/sonic-vlan/vlan-member/vlan-member-list", VlanMemberHandler).
        AddRouter("/sonic-vlan/vlan-interface/vlan-interface-list", VlanInterfaceHandler).
        AddRouter("/sonic-vlan/vlan-interface/vlan-interface-ipprefix-list", VlanInterfaceIPPrefixHandler).
        AddRouter("/sonic-interface/interface/interface-list", InterfaceHandler).
        AddRouter("/sonic-interface/interface/interface-ipprefix-list", InterfaceIPPrefixHandler).
        AddRouter("/sonic-acl/acl-rule/acl-rule-list", AclRuleHandler).
        AddRouter("/sonic-acl/acl-table/acl-table-list", AclTableHandler)
}
