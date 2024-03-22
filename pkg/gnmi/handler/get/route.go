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
        AddRouter("/accton-mdns", MdnsInfoHandler).
        AddRouter("/accton-system-top", SysTopInfoHandler).
        AddRouter("/accton-system-top/cpus", SysTopInfoCpuHandler).
        AddRouter("/accton-system-top/memory", SysTopInfoMemHandler).
        AddRouter("/accton-system-top/disk", SysTopInfoDiskHandler).
        AddRouter("/accton-platform/platform", ComponentInfoHandler).
        AddRouter("/accton-platform/platform/component-list/fan", FanInfoHandler).
        AddRouter("/accton-platform/platform/component-list/power-supply", PowerSupplyInfoHandler).
        AddRouter("/accton-platform/platform/component-list/state/temperature", TemperatureInfoHandler).
        AddRouter("/accton-platform/platform/component-list/system", SystemInfoHandler).
        AddRouter("/accton-lldp/lldp/lldp-list", LLDPHandler).
        AddRouter("/accton-port/port/port-list", PortHandler).
        AddRouter("/accton-port/port/port-statistics-list", PortStatisticsHandler).
        AddRouter("/accton-portchannel/portchannel/portchannel-list", PortChannelHandler).
        AddRouter("/accton-portchannel/portchannel-member/portchannel-member-list", PortChannelMemberHandler).
        AddRouter("/accton-vlan/vlan/vlan-list", VlanHandler).
        AddRouter("/accton-vlan/vlan-member/vlan-member-list", VlanMemberHandler).
        AddRouter("/accton-vlan/vlan-interface/vlan-interface-list", VlanInterfaceHandler).
        AddRouter("/accton-vlan/vlan-interface/vlan-interface-ipprefix-list", VlanInterfaceIPPrefixHandler).
        AddRouter("/accton-interface/interface/interface-list", InterfaceHandler).
        AddRouter("/accton-interface/interface/interface-ipprefix-list", InterfaceIPPrefixHandler).
        AddRouter("/accton-loopback-interface/loopback-interface/loopback-interface-list", LoopbackInterfaceHandler).
        AddRouter("/accton-loopback-interface/loopback-interface/loopback-interface-ipprefix-list", LoopbackInterfaceIPPrefixHandler).
        AddRouter("/accton-acl/acl-rule/acl-rule-list", AclRuleHandler).
        AddRouter("/accton-acl/acl-table/acl-table-list", AclTableHandler).
        AddRouter("/accton-fdb/fdb/fdb-list", FdbHandler).
        AddRouter("/accton-route/route/route-list", IpRouteHandler).
        AddRouter("/accton-vrf/vrf/vrf-list", VrfHandler).
        AddRouter("/accton-ntp/ntp/ntp-list", NtpHandler).
        AddRouter("/accton-neighor/neighor/neighor-list", NeighborHandler)
}
