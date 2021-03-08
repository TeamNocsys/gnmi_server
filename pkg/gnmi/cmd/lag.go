package cmd

import (
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/gnmi/proto/gnmi"
    "github.com/openconfig/ygot/proto/ywrapper"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
    "strconv"
)

type LagAdapter struct {
    Adapter
    name string
}

func NewLagAdapter(name string, cli command.Client) *LagAdapter {
    return &LagAdapter{
        Adapter: Adapter{
            client: cli,
        },
        name:  name,
    }
}

func (adpt *LagAdapter) Show(dataType gnmi.GetRequest_DataType) (*sonicpb.NocsysPortchannel_Portchannel_PortchannelList, error) {
    conn := adpt.client.Config()
    if conn == nil {
        return nil, swsssdk.ErrConnNotExist
    }

    if data, err := conn.GetAll(swsssdk.CONFIG_DB, []string{"PORTCHANNEL", adpt.name}); err != nil {
        return nil, err
    } else {
        retval := &sonicpb.NocsysPortchannel_Portchannel_PortchannelList{}
        for k, v := range data {
            switch k {
            case "mtu":
                if i, err := strconv.ParseUint(v, 10, 64); err != nil {
                    return nil, err
                } else {
                    retval.Mtu = &ywrapper.UintValue{Value: i}
                }
            case "admin_status":
                switch v {
                case "up":
                    retval.AdminStatus = sonicpb.NocsysTypesAdminStatus_NOCSYSTYPESADMINSTATUS_up
                case "down":
                    retval.AdminStatus = sonicpb.NocsysTypesAdminStatus_NOCSYSTYPESADMINSTATUS_down
                }
            case "min_links":
                if i, err := strconv.ParseUint(v, 10, 64); err != nil {
                    return nil, err
                } else {
                    retval.MinLinks = &ywrapper.UintValue{Value: i}
                }
            }
        }
        return retval, nil
    }
}

func (adpt *LagAdapter) Config(data *sonicpb.NocsysPortchannel_Portchannel_PortchannelList, oper OperType) error {
    cmdstr := "config portchannel"
    if oper == ADD {
        conn := adpt.client.Config()
        if conn == nil {
            return swsssdk.ErrConnNotExist
        }
        if ok, err := conn.HasEntry("PORTCHANNEL", adpt.name); err != nil {
            return err
        } else if ok {
            return nil
        }

        cmdstr += " add " + adpt.name
    } else if oper == DEL {
        conn := adpt.client.Config()
        if conn == nil {
            return swsssdk.ErrConnNotExist
        }
        if ok, err := conn.HasEntry("PORTCHANNEL", adpt.name); err != nil {
            return err
        } else if !ok {
            return nil
        }

        cmdstr += " del " + adpt.name
    } else {
        return ErrInvalidOperType
    }

    return adpt.exec(cmdstr)
}