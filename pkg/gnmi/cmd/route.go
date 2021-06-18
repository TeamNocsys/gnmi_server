package cmd

import (
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/gnmi/proto/gnmi"
    "github.com/openconfig/ygot/proto/ywrapper"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"

    "strings"
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
        cmdstr += " add"
    } else if oper == DEL {
        cmdstr += " del"
    } else {
        return ErrInvalidOperType
    }

    cmdstr += " prefix vrf " + adpt.vrf + " " + adpt.dest

    /* cmd example:
     * config route add [OPTIONS] prefix [vrf <vrf_name>] <A.B.C.D/M> nexthop
     *                  <[vrf <vrf_name>] <A.B.C.D>>|<dev <dev_name>>
     */
    nh_ar := []string {}
    if_ar := []string {}

    if data.Nexthop != nil {
        nh_ar = strings.Split(data.Nexthop.Value, ",")
    }

    if data.Ifname != nil {
        if_ar = strings.Split(data.Ifname.Value, ",")
    }

    if len(nh_ar) == 0 {
        if len(if_ar) != 0 {
            nh_ar = make ([]string, len(if_ar))
        }
    } else {
        if len (if_ar) == 0 {
            if_ar = make ([]string, len(nh_ar))
        }
    }

    if len(nh_ar) != len (if_ar) || len(nh_ar) == 0 {
        return ErrUnknown
    }

    for idx := range nh_ar {
        exec_str := cmdstr
        if nh_ar [idx] == "0.0.0.0" || nh_ar [idx] == "" {
            exec_str += " nexthop dev " + if_ar[idx]
        } else {
            exec_str += " nexthop vrf " + adpt.vrf + " " + nh_ar[idx]
        }

        if err := adpt.exec(exec_str); err != nil {
            return err
        }
    }

    return nil
}

