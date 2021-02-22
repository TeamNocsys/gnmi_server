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
        AddRouter("/nocsys-platform/platform", ComponentInfoHandler).
        AddRouter("/nocsys-platform/platform/component-list/fan", FanInfoHandler).
        AddRouter("/nocsys-platform/platform/component-list/power-supply", PowerSupplyInfoHandler).
        AddRouter("/nocsys-platform/platform/component-list/state/temperature", TemperatureInfoHandler).
        AddRouter("/nocsys-platform/platform/component-list/system", SystemInfoHandler).
        AddRouter("/nocsys-lldp/lldp/lldp-list", LLDPHandler).
        AddRouter("/nocsys-port/port/port-list", PortHandler).
        AddRouter("/nocsys-port/port/port-statistics-list", PortStatisticsHandler).
        AddRouter("/nocsys-portchannel/portchannel/portchannel-list", PortChannelHandler).
        AddRouter("/nocsys-portchannel/portchannel-member/portchannel-member-list", PortChannelMemberHandler).
        AddRouter("/nocsys-vlan/vlan/vlan-list", VlanHandler).
        AddRouter("/nocsys-vlan/vlan-member/vlan-member-list", VlanMemberHandler).
        AddRouter("/nocsys-vlan/vlan-interface/vlan-interface-list", VlanInterfaceHandler).
        AddRouter("/nocsys-vlan/vlan-interface/vlan-interface-ipprefix-list", VlanInterfaceIPPrefixHandler).
        AddRouter("/nocsys-interface/interface/interface-list", InterfaceHandler).
        AddRouter("/nocsys-interface/interface/interface-ipprefix-list", InterfaceIPPrefixHandler).
        AddRouter("/nocsys-loopback-interface/loopback-interface/loopback-interface-list", LoopbackInterfaceHandler).
        AddRouter("/nocsys-loopback-interface/loopback-interface/loopback-interface-ipprefix-list", LoopbackInterfaceIPPrefixHandler).
        AddRouter("/nocsys-acl/acl-rule/acl-rule-list", AclRuleHandler).
        AddRouter("/nocsys-acl/acl-table/acl-table-list", AclTableHandler).
        AddRouter("/nocsys-fdb/fdb/fdb-list", FdbHandler).
        AddRouter("/nocsys-route/route/route-list", IpRouteHandler).
        AddRouter("/nocsys-vrf/vrf/vrf-list", VrfHandler).
        AddRouter("/nocsys-ntp/ntp/ntp-list", NtpHandler).
        AddRouter("/nocsys-neighor/neighor/neighor-list", NeighborHandler)
}
