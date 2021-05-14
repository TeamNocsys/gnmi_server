package cmd

import (
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/gnmi/proto/gnmi"
    "github.com/openconfig/ygot/proto/ywrapper"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
    "strings"
)

type VrfAdapter struct {
    Adapter
    name string
}

func NewVrfAdapter(name string, cli command.Client) *VrfAdapter {
    return &VrfAdapter{
        Adapter: Adapter{
            client: cli,
        },
        name:  name,
    }
}

func (adpt *VrfAdapter) Show(dataType gnmi.GetRequest_DataType) (*sonicpb.NocsysVrf_Vrf_VrfList, error) {
    conn := adpt.client.Config()
    if conn == nil {
        return nil, swsssdk.ErrConnNotExist
    }

    if data, err := conn.GetAll(swsssdk.CONFIG_DB, []string{"VRF", adpt.name}); err != nil {
        return nil, err
    } else {
        retval := &sonicpb.NocsysVrf_Vrf_VrfList{}
        for k, v := range data {
            switch k {
            case "fallback":
                switch strings.ToLower(v) {
                case "true":
                    retval.Fallback = &ywrapper.BoolValue{Value: true}
                case "false":
                    retval.Fallback = &ywrapper.BoolValue{Value: false}
                }
            }
        }
        return retval, nil
    }
}