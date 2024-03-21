package cmd

import (
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/gnmi/proto/gnmi"
    "github.com/openconfig/ygot/proto/ywrapper"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
)

type NeighborAdapter struct {
    Adapter
    ifname string
    ipaddr string
}

func NewNeighborAdapter(ifname, ipaddr string, cli command.Client) *NeighborAdapter {
    return &NeighborAdapter{
        Adapter: Adapter{
            client: cli,
        },
        ifname: ifname,
        ipaddr: ipaddr,
    }
}

func (adpt *NeighborAdapter) Show(dataType gnmi.GetRequest_DataType) (*sonicpb.AcctonNeighor_Neighor_NeighorList, error) {
    conn := adpt.client.State()
    if conn == nil {
        return nil, swsssdk.ErrConnNotExist
    }

    if data, err := conn.GetAll(swsssdk.APPL_DB, append([]string{"NEIGH_TABLE"}, adpt.ifname, adpt.ipaddr)); err != nil {
        return nil, err
    } else {
        retval := &sonicpb.AcctonNeighor_Neighor_NeighorList{}
        for k, v := range data {
            switch k {
            case "neigh":
                retval.Neigh = &ywrapper.StringValue{Value: v}
            }
        }
        return retval, nil
    }
}