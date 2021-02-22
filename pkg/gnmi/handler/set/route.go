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
    mux.AddDeleteRouter("/nocsys-port/port/port-list", delete.PortHandler).
        AddDeleteRouter("/nocsys-vlan/vlan/vlan-list", delete.VlanHandler).
        AddDeleteRouter("/nocsys-vlan/vlan-member/vlan-member-list", delete.VlanMemberHandler).
        AddDeleteRouter("/nocsys-vlan/vlan-interface/vlan-interface-list", delete.VlanInterfaceHandler).
        AddDeleteRouter("/nocsys-vlan/vlan-interface/vlan-interface-ipprefix-list", delete.VlanInterfaceIPPrefixHandler).
        AddDeleteRouter("/nocsys-portchannel/portchannel/portchannel-list", delete.PortChannelHandler).
        AddDeleteRouter("/nocsys-portchannel/portchannel-member/portchannel-member-list", delete.PortChannelMemberHandler).
        AddDeleteRouter("/nocsys-interface/interface/interface-list", delete.InterfaceHandler).
        AddDeleteRouter("/nocsys-interface/interface/interface-ipprefix-list", delete.InterfaceIPPrefixHandler).
        AddDeleteRouter("/nocsys-loopback-interface/loopback-interface/loopback-interface-list", delete.LoopbackInterfaceHandler).
        AddDeleteRouter("/nocsys-loopback-interface/loopback-interface/loopback-interface-ipprefix-list", delete.LoopbackInterfaceIPPrefixHandler).
        AddDeleteRouter("/nocsys-acl/acl-table/acl-table-list", delete.AclTableHandler).
        AddDeleteRouter("/nocsys-acl/acl-rule/acl-rule-list", delete.AclRuleHandler).
        AddDeleteRouter("/nocsys-mirror-session/mirror-session/mirror-session-list", delete.MirrorSessionHandler).
        AddDeleteRouter("/nocsys-fdb/fdb/fdb-list", delete.FdbHandler).
        AddDeleteRouter("/nocsys-route/route/route-list", delete.IpRouteHandler).
        AddDeleteRouter("/nocsys-vrf/vrf/vrf-list", delete.VrfHandler).
        AddDeleteRouter("/nocsys-ntp/ntp/ntp-list", delete.NtpHandler).
        AddDeleteRouter("/nocsys-todo/todo/todo-list", delete.TodoHandler).
        AddDeleteRouter("/nocsys-neighor/neighor/neighor-list", delete.NeighborHandler)

    // 设置替换路由
    mux.AddReplaceRouter("/nocsys-port/port/port-list", replace.PortHandler).
        AddReplaceRouter("/nocsys-vlan/vlan/vlan-list", replace.VlanHandler).
        AddReplaceRouter("/nocsys-vlan/vlan-member/vlan-member-list", replace.VlanMemberHandler).
        AddReplaceRouter("/nocsys-vlan/vlan-interface/vlan-interface-list", replace.VlanInterfaceHandler).
        AddReplaceRouter("/nocsys-vlan/vlan-interface/vlan-interface-ipprefix-list", replace.VlanInterfaceIPPrefixHandler).
        AddReplaceRouter("/nocsys-portchannel/portchannel/portchannel-list", replace.PortChannelHandler).
        AddReplaceRouter("/nocsys-portchannel/portchannel-member/portchannel-member-list", replace.PortChannelMemberHandler).
        AddReplaceRouter("/nocsys-interface/interface/interface-list", replace.InterfaceHandler).
        AddReplaceRouter("/nocsys-interface/interface/interface-ipprefix-list", replace.InterfaceIPPrefixHandler).
        AddReplaceRouter("/nocsys-loopback-interface/loopback-interface/loopback-interface-list", replace.LoopbackInterfaceHandler).
        AddReplaceRouter("/nocsys-loopback-interface/loopback-interface/loopback-interface-ipprefix-list", replace.LoopbackInterfaceIPPrefixHandler).
        AddReplaceRouter("/nocsys-acl/acl-table/acl-table-list", replace.AclTableHandler).
        AddReplaceRouter("/nocsys-acl/acl-rule/acl-rule-list", replace.AclRuleHandler).
        AddReplaceRouter("/nocsys-mirror-session/mirror-session/mirror-session-list", replace.MirrorSessionHandler).
        AddReplaceRouter("/nocsys-fdb/fdb/fdb-list", replace.FdbHandler).
        AddReplaceRouter("/nocsys-route/route/route-list", replace.IpRouteHandler).
        AddReplaceRouter("/nocsys-vrf/vrf/vrf-list", replace.VrfHandler).
        AddReplaceRouter("/nocsys-ntp/ntp/ntp-list", replace.NtpHandler).
        AddReplaceRouter("/nocsys-todo/todo/todo-list", replace.TodoHandler).
        AddReplaceRouter("/nocsys-neighor/neighor/neighor-list", replace.NeighborHandler)

    // 设置更新路由
    //mux.AddUpdateRouter("/sonic", SonicUpdateHandler).
    mux.AddUpdateRouter("/nocsys-port/port/port-list", update.PortHandler).
        AddUpdateRouter("/nocsys-vlan/vlan/vlan-list", update.VlanHandler).
        AddUpdateRouter("/nocsys-vlan/vlan-member/vlan-member-list", update.VlanMemberHandler).
        AddUpdateRouter("/nocsys-vlan/vlan-interface/vlan-interface-list", update.VlanInterfaceHandler).
        AddUpdateRouter("/nocsys-vlan/vlan-interface/vlan-interface-ipprefix-list", update.VlanInterfaceIPPrefixHandler).
        AddUpdateRouter("/nocsys-portchannel/portchannel/portchannel-list", update.PortChannelHandler).
        AddUpdateRouter("/nocsys-portchannel/portchannel-member/portchannel-member-list", update.PortChannelMemberHandler).
        AddUpdateRouter("/nocsys-interface/interface/interface-list", update.InterfaceHandler).
        AddUpdateRouter("/nocsys-interface/interface/interface-ipprefix-list", update.InterfaceIPPrefixHandler).
        AddUpdateRouter("/nocsys-loopback-interface/loopback-interface/loopback-interface-list", update.LoopbackInterfaceHandler).
        AddUpdateRouter("/nocsys-loopback-interface/loopback-interface/loopback-interface-ipprefix-list", update.LoopbackInterfaceIPPrefixHandler).
        AddUpdateRouter("/nocsys-acl/acl-table/acl-table-list", update.AclTableHandler).
        AddUpdateRouter("/nocsys-acl/acl-rule/acl-rule-list", update.AclRuleHandler).
        AddUpdateRouter("/nocsys-mirror-session/mirror-session/mirror-session-list", update.MirrorSessionHandler).
        AddUpdateRouter("/nocsys-fdb/fdb/fdb-list", update.FdbHandler).
        AddUpdateRouter("/nocsys-route/route/route-list", update.IpRouteHandler).
        AddUpdateRouter("/nocsys-vrf/vrf/vrf-list", update.VrfHandler).
        AddUpdateRouter("/nocsys-ntp/ntp/ntp-list", update.NtpHandler).
        AddUpdateRouter("/nocsys-todo/todo/todo-list", update.TodoHandler).
        AddUpdateRouter("/nocsys-neighor/neighor/neighor-list", update.NeighborHandler)
}
