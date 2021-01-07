package set

import (
    "gnmi_server/pkg/gnmi"
)

func SetServeMux() *gnmi.SetServeMux {
    mux := gnmi.NewSetServeMux()
    route(mux)
    return mux
}

func route(mux *gnmi.SetServeMux) {
    mux.AddUpdateRouter("/sonic-port/port/port-list", PortListHandler)
    mux.AddUpdateRouter("/sonic-loopback-interface/loopback-interface", LoopbackUpdateHandler)
    mux.AddDeleteRouter("/sonic-loopback-interface/loopback-interface", LoopbackDeleteHandler)
    mux.AddUpdateRouter("/sonic", SonicUpdateHandler)
    mux.AddReplaceRouter("/sonic-portchannel/portchannel", PortChannelReplaceHandler)
    mux.AddDeleteRouter("/sonic-portchannel/portchannel/portchannel-list/portchannel-name",
        PortChannelDeleteHandler)
    mux.AddUpdateRouter("/sonic-portchannel/portchannel/portchannel-list/members", PortChannelMemberUpdateHandler)
    mux.AddDeleteRouter("/sonic-portchannel/portchannel/portchannel-list/members", PortChannelMemberDeleteHandler)
}
