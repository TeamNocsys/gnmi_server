package set

import (
    "gnmi_server/pkg/gnmi"
    "gnmi_server/pkg/gnmi/handler/set/delete"
    "gnmi_server/pkg/gnmi/handler/set/replace"
    "gnmi_server/pkg/gnmi/handler/set/update"
)

func SetServeMux() *gnmi.SetServeMux {
    mux := gnmi.NewSetServeMux()
    route(mux)
    return mux
}

func route(mux *gnmi.SetServeMux) {
    // 设置删除路由
    mux.AddDeleteRouter("/sonic-port/port/port-list", delete.PortHandler).
        AddDeleteRouter("/sonic-vlan/vlan/vlan-list", delete.VlanHandler).
        AddDeleteRouter("/sonic-vlan/vlan-member/vlan-member-list", delete.VlanMemberHandler).
        AddDeleteRouter("/sonic-vlan/vlan-interface/vlan-interface-list", delete.VlanInterfaceHandler).
        AddDeleteRouter("/sonic-vlan/vlan-interface/vlan-interface-ipprefix-list", delete.VlanInterfaceIPPrefixHandler).
        AddDeleteRouter("/sonic-portchannel/portchannel/portchannel-list", delete.PortChannelHandler).
        AddDeleteRouter("/sonic-portchannel/portchannel-member/portchannel-member-list", delete.PortChannelMemberHandler).
        AddDeleteRouter("/sonic-interface/interface/interface-list", delete.InterfaceHandler).
        AddDeleteRouter("/sonic-interface/interface/interface-ipprefix-list", delete.InterfaceIPPrefixHandler).
        AddDeleteRouter("/sonic-loopback-interface/loopback-interface/loopback-interface-list", delete.LoopbackInterfaceHandler).
        AddDeleteRouter("/sonic-loopback-interface/loopback-interface/loopback-interface-ipprefix-list", delete.LoopbackInterfaceIPPrefixHandler).
        AddDeleteRouter("/sonic-acl/acl-table/acl-table-list", delete.AclTableHandler).
        AddDeleteRouter("/sonic-acl/acl-rule/acl-rule-list", delete.AclRuleHandler).
        AddDeleteRouter("/sonic-mirror-session/mirror-session/mirror-session-list", delete.MirrorSessionHandler).
        AddDeleteRouter("/sonic-fdb/fdb/fdb-list", delete.FdbHandler).
        AddDeleteRouter("/sonic-route/route/route-list", delete.IpRouteHandler).
        AddDeleteRouter("/sonic-route/route/global-route-list", delete.GlobalIpRouteHandler).
        AddDeleteRouter("/sonic-vrf/vrf/vrf-list", delete.VrfHandler).
        AddDeleteRouter("/sonic-ntp/ntp/ntp-list", delete.NtpHandler).
        AddDeleteRouter("/sonic-todo/todo/todo-list", delete.TodoHandler).
        AddDeleteRouter("/sonic-neighor/neighor/neighor-list", delete.NeighborHandler)

    // 设置替换路由
    mux.AddReplaceRouter("/sonic-port/port/port-list", replace.PortHandler).
        AddReplaceRouter("/sonic-vlan/vlan/vlan-list", replace.VlanHandler).
        AddReplaceRouter("/sonic-vlan/vlan-member/vlan-member-list", replace.VlanMemberHandler).
        AddReplaceRouter("/sonic-vlan/vlan-interface/vlan-interface-list", replace.VlanInterfaceHandler).
        AddReplaceRouter("/sonic-vlan/vlan-interface/vlan-interface-ipprefix-list", replace.VlanInterfaceIPPrefixHandler).
        AddReplaceRouter("/sonic-portchannel/portchannel/portchannel-list", replace.PortChannelHandler).
        AddReplaceRouter("/sonic-portchannel/portchannel-member/portchannel-member-list", replace.PortChannelMemberHandler).
        AddReplaceRouter("/sonic-interface/interface/interface-list", replace.InterfaceHandler).
        AddReplaceRouter("/sonic-interface/interface/interface-ipprefix-list", replace.InterfaceIPPrefixHandler).
        AddReplaceRouter("/sonic-loopback-interface/loopback-interface/loopback-interface-list", replace.LoopbackInterfaceHandler).
        AddReplaceRouter("/sonic-loopback-interface/loopback-interface/loopback-interface-ipprefix-list", replace.LoopbackInterfaceIPPrefixHandler).
        AddReplaceRouter("/sonic-acl/acl-table/acl-table-list", replace.AclTableHandler).
        AddReplaceRouter("/sonic-acl/acl-rule/acl-rule-list", replace.AclRuleHandler).
        AddReplaceRouter("/sonic-mirror-session/mirror-session/mirror-session-list", replace.MirrorSessionHandler).
        AddReplaceRouter("/sonic-fdb/fdb/fdb-list", replace.FdbHandler).
        AddReplaceRouter("/sonic-route/route/route-list", replace.IpRouteHandler).
        AddReplaceRouter("/sonic-route/route/global-route-list", replace.GlobalIpRouteHandler).
        AddReplaceRouter("/sonic-vrf/vrf/vrf-list", replace.VrfHandler).
        AddReplaceRouter("/sonic-ntp/ntp/ntp-list", replace.NtpHandler).
        AddReplaceRouter("/sonic-todo/todo/todo-list", replace.TodoHandler).
        AddReplaceRouter("/sonic-neighor/neighor/neighor-list", replace.NeighborHandler)

    // 设置更新路由
    //mux.AddUpdateRouter("/sonic", SonicUpdateHandler).
    mux.AddUpdateRouter("/sonic-port/port/port-list", update.PortHandler).
        AddUpdateRouter("/sonic-vlan/vlan/vlan-list", update.VlanHandler).
        AddUpdateRouter("/sonic-vlan/vlan-member/vlan-member-list", update.VlanMemberHandler).
        AddUpdateRouter("/sonic-vlan/vlan-interface/vlan-interface-list", update.VlanInterfaceHandler).
        AddUpdateRouter("/sonic-vlan/vlan-interface/vlan-interface-ipprefix-list", update.VlanInterfaceIPPrefixHandler).
        AddUpdateRouter("/sonic-portchannel/portchannel/portchannel-list", update.PortChannelHandler).
        AddUpdateRouter("/sonic-portchannel/portchannel-member/portchannel-member-list", update.PortChannelMemberHandler).
        AddUpdateRouter("/sonic-interface/interface/interface-list", update.InterfaceHandler).
        AddUpdateRouter("/sonic-interface/interface/interface-ipprefix-list", update.InterfaceIPPrefixHandler).
        AddUpdateRouter("/sonic-loopback-interface/loopback-interface/loopback-interface-list", update.LoopbackInterfaceHandler).
        AddUpdateRouter("/sonic-loopback-interface/loopback-interface/loopback-interface-ipprefix-list", update.LoopbackInterfaceIPPrefixHandler).
        AddUpdateRouter("/sonic-acl/acl-table/acl-table-list", update.AclTableHandler).
        AddUpdateRouter("/sonic-acl/acl-rule/acl-rule-list", update.AclRuleHandler).
        AddUpdateRouter("/sonic-mirror-session/mirror-session/mirror-session-list", update.MirrorSessionHandler).
        AddUpdateRouter("/sonic-fdb/fdb/fdb-list", update.FdbHandler).
        AddUpdateRouter("/sonic-route/route/route-list", update.IpRouteHandler).
        AddUpdateRouter("/sonic-route/route/global-route-list", update.GlobalIpRouteHandler).
        AddUpdateRouter("/sonic-vrf/vrf/vrf-list", update.VrfHandler).
        AddUpdateRouter("/sonic-ntp/ntp/ntp-list", update.NtpHandler).
        AddUpdateRouter("/sonic-todo/todo/todo-list", update.TodoHandler).
        AddUpdateRouter("/sonic-neighor/neighor/neighor-list", update.NeighborHandler)
}
