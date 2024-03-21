package cmd

import (
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/gnmi/proto/gnmi"
    "github.com/openconfig/ygot/proto/ywrapper"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
    "strings"
)

type FdbAdapter struct {
    Adapter
    mac string
}

func NewFdbAdapter(mac string, cli command.Client) *FdbAdapter {
    return &FdbAdapter{
        Adapter: Adapter{
            client: cli,
        },
        mac:  mac,
    }
}

func (adpt *FdbAdapter) Show(dataType gnmi.GetRequest_DataType) (*sonicpb.AcctonFdb_Fdb_FdbList, error) {
    conn := adpt.client.State()
    if conn == nil {
        return nil, swsssdk.ErrConnNotExist
    }

    if data, err := conn.GetAll(swsssdk.STATE_DB, []string{"FDB_TABLE", adpt.mac}); err != nil {
        return nil, err
    } else {
        retval := &sonicpb.AcctonFdb_Fdb_FdbList{}
        for k, v := range data {
            switch k {
            case "type":
                switch strings.ToUpper(v) {
                case "STATIC":
                    retval.Type = sonicpb.AcctonFdb_Fdb_FdbList_TYPE_STATIC
                case "DYNAMIC":
                    retval.Type = sonicpb.AcctonFdb_Fdb_FdbList_TYPE_DYNAMIC
                }
            case "port":
                retval.Port = &ywrapper.StringValue{Value: v}
            }
        }
        return retval, nil
    }
}