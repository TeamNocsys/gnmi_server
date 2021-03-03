package cmd

import (
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/gnmi/proto/gnmi"
    "github.com/openconfig/ygot/proto/ywrapper"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
)

type VrfRouteAdapter struct {
    Adapter
    vrf string
    dest string
}

func NewVrfRouteAdapter(vrf, dest string, cli command.Client) *VrfRouteAdapter {
    return &VrfRouteAdapter{
        Adapter: Adapter{
            client: cli,
        },
        vrf: vrf,
        dest: dest,
    }
}

func (adpt *VrfRouteAdapter) Show(dataType gnmi.GetRequest_DataType) (*sonicpb.NocsysRoute_Route_RouteList, error) {
    conn := adpt.client.State()
    if conn == nil {
        return nil, swsssdk.ErrConnNotExist
    }

    if data, err := conn.GetAll(swsssdk.APPL_DB, append([]string{"ROUTE_TABLE"}, adpt.vrf, adpt.dest)); err != nil {
        return nil, err
    } else {
        retval := &sonicpb.NocsysRoute_Route_RouteList{}
        for k, v := range data {
            switch k {
            case "nexthop":
                retval.Nexthop = &ywrapper.StringValue{Value: v}
            case "ifname":
                retval.Ifname = &ywrapper.StringValue{Value: v}
            }
        }
        return retval, nil
    }
}

func (adpt *VrfRouteAdapter) Config(data *sonicpb.NocsysRoute_Route_RouteList, oper OperType) error {
    cmdstr := "config route"
    if oper == ADD {
        cmdstr += " add "
    } else if oper == DEL {
        cmdstr += " del "
    } else {
        return ErrInvalidOperType
    }

    if data.Nexthop != nil {
        cmdstr += "nexthop vrf " + adpt.vrf + " " + data.Nexthop.Value
    } else if data.Ifname != nil {
        cmdstr += "nexthop vrf " + adpt.vrf + " dev " + data.Ifname.Value
    } else {
        return ErrUnknown
    }

    return adpt.exec(cmdstr)
}