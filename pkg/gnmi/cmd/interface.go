package cmd

import (
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/gnmi/proto/gnmi"
    "github.com/openconfig/ygot/proto/ywrapper"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
)

type IfType int32

const (
    INTERFACE           IfType = 0
    LOOPBACK_INTERFACE  IfType = 1
    VLAN_INTERFACE      IfType = 2
)

var IfType_table = map[int32]string{
    0: "INTERFACE",
    1: "LOOPBACK_INTERFACE",
    2: "VLAN_INTERFACE",
}

type IfAdapter struct {
    Adapter
    ifType IfType
    ifname string
}

func NewIfAdapter(ifType IfType, ifname string, cli command.Client) *IfAdapter {
    return &IfAdapter{
        Adapter: Adapter{
            client: cli,
        },
        ifType:  ifType,
        ifname:  ifname,
    }
}

func (adpt *IfAdapter) Show(dataType gnmi.GetRequest_DataType) (interface{}, error) {
    conn := adpt.client.Config()
    if conn == nil {
        return nil, swsssdk.ErrConnNotExist
    }

    if data, err := conn.GetAll(swsssdk.CONFIG_DB, []string{IfType_table[int32(adpt.ifType)], adpt.ifname}); err != nil {
        return nil, err
    } else {
        var vrf string
        for k, v := range data {
            switch k {
            case "vrf_name":
                vrf = v
            }
        }

        if adpt.ifType == INTERFACE {
            return &sonicpb.NocsysInterface_Interface_InterfaceList{
                VrfName: &ywrapper.StringValue{Value: vrf},
            }, nil
        } else if adpt.ifType == VLAN_INTERFACE {
            return &sonicpb.NocsysVlan_VlanInterface_VlanInterfaceList{
                VrfName: &ywrapper.StringValue{Value: vrf},
            }, nil
        } else if adpt.ifType == LOOPBACK_INTERFACE {
            return &sonicpb.NocsysLoopbackInterface_LoopbackInterface_LoopbackInterfaceList{
                VrfName: &ywrapper.StringValue{Value: vrf},
            }, nil
        }
    }

    return nil, ErrUnknown
}

func (adpt *IfAdapter) Config(data interface{}, oper OperType) error {
    var vrf string
    if adpt.ifType == INTERFACE {
        if v, ok := data.(*sonicpb.NocsysInterface_Interface_InterfaceList); !ok {
            return ErrTypeConversion
        } else {
            if v.VrfName != nil {
                vrf = v.VrfName.Value
            } else {
                return nil
            }
        }
    } else if adpt.ifType == VLAN_INTERFACE {
        if v, ok := data.(*sonicpb.NocsysVlan_VlanInterface_VlanInterfaceList); !ok {
            return ErrTypeConversion
        } else {
            if v.VrfName != nil {
                vrf = v.VrfName.Value
            } else {
                return nil
            }
        }
    } else if adpt.ifType == LOOPBACK_INTERFACE {
        if v, ok := data.(*sonicpb.NocsysLoopbackInterface_LoopbackInterface_LoopbackInterfaceList); !ok {
            return ErrTypeConversion
        } else {
            if v.VrfName != nil {
                vrf = v.VrfName.Value
            } else {
                return nil
            }
        }
    }

    var cmdstr string
    if oper == ADD {
        cmdstr = "config interface vrf bind " + adpt.ifname + " " + vrf
    } else if oper == DEL {
        cmdstr = "config interface vrf unbind " + adpt.ifname + " " + vrf
    } else {
        return ErrInvalidOperType
    }
    return adpt.exec(cmdstr)
}