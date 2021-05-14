// +build broadcom

package cmd

import (
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
)

func (adpt *VrfRouteAdapter) Config(data *sonicpb.NocsysRoute_Route_RouteList, oper OperType) error {
    cmdstr := "config route"
    if oper == ADD {
        cmdstr += " add"
    } else if oper == DEL {
        cmdstr += " del"
    } else {
        return ErrInvalidOperType
    }

    cmdstr += " prefix vrf " + adpt.vrf + " " + adpt.dest

    if data.Nexthop != nil {
        cmdstr += " nexthop vrf " + adpt.vrf + " " + data.Nexthop.Value
    } else if data.Ifname != nil {
        cmdstr += " nexthop vrf " + adpt.vrf + " dev " + data.Ifname.Value
    } else {
        return ErrUnknown
    }

    return adpt.exec(cmdstr)
}