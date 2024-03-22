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
    //mux.AddDeleteRouter("/accton-port/port/port-list", delete.PortHandler).
    mux.AddDeleteRouter("/accton-vlan/vlan/vlan-list", delete.VlanHandler).
        AddDeleteRouter("/accton-vlan/vlan-member/vlan-member-list", delete.VlanMemberHandler).
        AddDeleteRouter("/accton-vlan/vlan-interface/vlan-interface-list", delete.VlanInterfaceHandler).
        AddDeleteRouter("/accton-vlan/vlan-interface/vlan-interface-ipprefix-list", delete.VlanInterfaceIPPrefixHandler).
        AddDeleteRouter("/accton-portchannel/portchannel/portchannel-list", delete.PortChannelHandler).
        AddDeleteRouter("/accton-portchannel/portchannel-member/portchannel-member-list", delete.PortChannelMemberHandler).
        AddDeleteRouter("/accton-interface/interface/interface-list", delete.InterfaceHandler).
        AddDeleteRouter("/accton-interface/interface/interface-ipprefix-list", delete.InterfaceIPPrefixHandler).
        AddDeleteRouter("/accton-loopback-interface/loopback-interface/loopback-interface-list", delete.LoopbackInterfaceHandler).
        AddDeleteRouter("/accton-loopback-interface/loopback-interface/loopback-interface-ipprefix-list", delete.LoopbackInterfaceIPPrefixHandler).
        //AddDeleteRouter("/accton-acl/acl-table/acl-table-list", delete.AclTableHandler).
        //AddDeleteRouter("/accton-acl/acl-rule/acl-rule-list", delete.AclRuleHandler).
        AddDeleteRouter("/accton-mirror-session/mirror-session/mirror-session-list", delete.MirrorSessionHandler).
        //AddDeleteRouter("/accton-fdb/fdb/fdb-list", delete.FdbHandler).
        AddDeleteRouter("/accton-route/route/route-list", delete.IpRouteHandler).
        AddDeleteRouter("/accton-vrf/vrf/vrf-list", delete.VrfHandler).
        AddDeleteRouter("/accton-ntp/ntp/ntp-list", delete.NtpHandler)
        //AddDeleteRouter("/accton-todo/todo/todo-list", delete.TodoHandler).
        //AddDeleteRouter("/accton-neighor/neighor/neighor-list", delete.NeighborHandler)

    // 设置替换路由
    mux.AddReplaceRouter("/accton-port/port/port-list", replace.PortHandler).
        AddReplaceRouter("/accton-vlan/vlan/vlan-list", replace.VlanHandler).
        AddReplaceRouter("/accton-vlan/vlan-member/vlan-member-list", replace.VlanMemberHandler).
        AddReplaceRouter("/accton-vlan/vlan-interface/vlan-interface-list", replace.VlanInterfaceHandler).
        AddReplaceRouter("/accton-vlan/vlan-interface/vlan-interface-ipprefix-list", replace.VlanInterfaceIPPrefixHandler).
        AddReplaceRouter("/accton-portchannel/portchannel/portchannel-list", replace.PortChannelHandler).
        AddReplaceRouter("/accton-portchannel/portchannel-member/portchannel-member-list", replace.PortChannelMemberHandler).
        AddReplaceRouter("/accton-interface/interface/interface-list", replace.InterfaceHandler).
        AddReplaceRouter("/accton-interface/interface/interface-ipprefix-list", replace.InterfaceIPPrefixHandler).
        AddReplaceRouter("/accton-loopback-interface/loopback-interface/loopback-interface-list", replace.LoopbackInterfaceHandler).
        AddReplaceRouter("/accton-loopback-interface/loopback-interface/loopback-interface-ipprefix-list", replace.LoopbackInterfaceIPPrefixHandler).
        AddReplaceRouter("/accton-acl/acl-table/acl-table-list", replace.AclTableHandler).
        AddReplaceRouter("/accton-acl/acl-rule/acl-rule-list", replace.AclRuleHandler).
        AddReplaceRouter("/accton-mirror-session/mirror-session/mirror-session-list", replace.MirrorSessionHandler).
        //AddReplaceRouter("/accton-fdb/fdb/fdb-list", replace.FdbHandler).
        AddReplaceRouter("/accton-route/route/route-list", replace.IpRouteHandler).
        AddReplaceRouter("/accton-vrf/vrf/vrf-list", replace.VrfHandler).
        AddReplaceRouter("/accton-ntp/ntp/ntp-list", replace.NtpHandler).
        AddReplaceRouter("/accton-todo/todo/todo-list", replace.TodoHandler)
        //AddReplaceRouter("/accton-neighor/neighor/neighor-list", replace.NeighborHandler)

    // 设置更新路由
    mux.AddUpdateRouter("/accton-port/port/port-list", update.PortHandler).
        AddUpdateRouter("/accton-vlan/vlan/vlan-list", update.VlanHandler).
        AddUpdateRouter("/accton-vlan/vlan-member/vlan-member-list", update.VlanMemberHandler).
        //AddUpdateRouter("/accton-vlan/vlan-interface/vlan-interface-list", update.VlanInterfaceHandler).
        //AddUpdateRouter("/accton-vlan/vlan-interface/vlan-interface-ipprefix-list", update.VlanInterfaceIPPrefixHandler).
        AddUpdateRouter("/accton-portchannel/portchannel/portchannel-list", update.PortChannelHandler).
        AddUpdateRouter("/accton-portchannel/portchannel-member/portchannel-member-list", update.PortChannelMemberHandler).
        //AddUpdateRouter("/accton-interface/interface/interface-list", update.InterfaceHandler).
        //AddUpdateRouter("/accton-interface/interface/interface-ipprefix-list", update.InterfaceIPPrefixHandler).
        //AddUpdateRouter("/accton-loopback-interface/loopback-interface/loopback-interface-list", update.LoopbackInterfaceHandler).
        //AddUpdateRouter("/accton-loopback-interface/loopback-interface/loopback-interface-ipprefix-list", update.LoopbackInterfaceIPPrefixHandler).
        AddUpdateRouter("/accton-acl/acl-table/acl-table-list", update.AclTableHandler).
        AddUpdateRouter("/accton-acl/acl-rule/acl-rule-list", update.AclRuleHandler).
        //AddUpdateRouter("/accton-mirror-session/mirror-session/mirror-session-list", update.MirrorSessionHandler).
        //AddUpdateRouter("/accton-fdb/fdb/fdb-list", update.FdbHandler).
        //AddUpdateRouter("/accton-route/route/route-list", update.IpRouteHandler).
        //AddUpdateRouter("/accton-vrf/vrf/vrf-list", update.VrfHandler).
        AddUpdateRouter("/accton-ntp/ntp/ntp-list", update.NtpHandler).
        AddUpdateRouter("/accton-todo/todo/todo-list", update.TodoHandler)
        //AddUpdateRouter("/accton-neighor/neighor/neighor-list", update.NeighborHandler)
}
