package cmd

import (
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/gnmi/proto/gnmi"
    "github.com/openconfig/ygot/proto/ywrapper"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
    "strconv"
)

type PortAdapter struct {
    Adapter
    ifname string
}

func NewPortAdapter(ifname string, cli command.Client) *PortAdapter {
    return &PortAdapter{
        Adapter: Adapter{
            client: cli,
        },
        ifname:  ifname,
    }
}

func (adpt *PortAdapter) Show(dataType gnmi.GetRequest_DataType) (*sonicpb.AcctonPort_Port_PortList, error) {
    retval := &sonicpb.AcctonPort_Port_PortList{}
    if dataType == gnmi.GetRequest_ALL || dataType == gnmi.GetRequest_CONFIG {
        conn := adpt.client.Config()
        if conn == nil {
            return nil, swsssdk.ErrConnNotExist
        }

        if data, err := conn.GetAll(swsssdk.CONFIG_DB, []string{"PORT", adpt.ifname}); err != nil {
            return nil, err
        } else {
            for k, v := range data {
                switch k {
                case "index":
                    if i, err := strconv.ParseUint(v, 10, 64); err != nil {
                        return nil, err
                    } else {
                        retval.Index = &ywrapper.UintValue{Value: i}
                    }
                case "lanes":
                    retval.Lanes = &ywrapper.StringValue{Value: v}
                case "mtu":
                    if i, err := strconv.ParseUint(v, 10, 64); err != nil {
                        return nil, err
                    } else {
                        retval.Mtu = &ywrapper.UintValue{Value: i}
                    }
                case "alias":
                    retval.Alias = &ywrapper.StringValue{Value: v}
                case "admin_status":
                    switch v {
                    case "up":
                        retval.AdminStatus = sonicpb.AcctonTypesAdminStatus_ACCTONTYPESADMINSTATUS_up
                    case "down":
                        retval.AdminStatus = sonicpb.AcctonTypesAdminStatus_ACCTONTYPESADMINSTATUS_down
                    }
                case "speed":
                    if i, err := strconv.ParseUint(v, 10, 64); err != nil {
                        return nil, err
                    } else {
                        retval.Speed = &ywrapper.UintValue{Value: i}
                    }
                }
            }
        }
    }

    if dataType == gnmi.GetRequest_ALL || dataType == gnmi.GetRequest_STATE {
        retval.State = &sonicpb.AcctonPort_Port_PortList_State{}
        conn := adpt.client.State()
        if conn == nil {
            return nil, swsssdk.ErrConnNotExist
        }

        if data, err := conn.GetAll(swsssdk.APPL_DB, []string{"PORT_TABLE", adpt.ifname}); err != nil {
            return nil, err
        } else {
            for k, v := range data {
                switch k {
                case "oper_status":
                    switch v {
                    case "up":
                        retval.State.OperStatus = sonicpb.AcctonTypesOperStatus_ACCTONTYPESOPERSTATUS_up
                    case "down":
                        retval.State.OperStatus = sonicpb.AcctonTypesOperStatus_ACCTONTYPESOPERSTATUS_down
                    }
                }
            }
        }
    }
    return retval, nil
}

func (adpt *PortAdapter) Config(data *sonicpb.AcctonPort_Port_PortList, oper OperType) error {
    if oper == ADD || oper == UPDATE {
        if data.Mtu != nil {
            cmdstr := "config interface mtu " + adpt.ifname + " " + strconv.FormatUint(data.Mtu.Value, 10)
            if err := adpt.exec(cmdstr); err != nil {
                return err
            }
        }

        if data.AdminStatus != sonicpb.AcctonTypesAdminStatus_ACCTONTYPESADMINSTATUS_UNSET {
            var cmdstr string
            if data.AdminStatus == sonicpb.AcctonTypesAdminStatus_ACCTONTYPESADMINSTATUS_up {
                cmdstr = "config interface startup " + adpt.ifname
            } else if data.AdminStatus == sonicpb.AcctonTypesAdminStatus_ACCTONTYPESADMINSTATUS_down {
                cmdstr = "config interface shutdown " + adpt.ifname
            }
            if err := adpt.exec(cmdstr); err != nil {
                return err
            }
        }

        if data.Speed != nil {
            cmdstr := "config interface speed " + adpt.ifname + " " + strconv.FormatUint(data.Mtu.Value, 10)
            if err := adpt.exec(cmdstr); err != nil {
                return err
            }
        }

        return nil
    }

    return ErrInvalidOperType
}