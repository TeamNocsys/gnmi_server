package cmd

import (
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/gnmi/proto/gnmi"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
)

type IfAddrAdapter struct {
    Adapter
    ifType IfType
    ifname string
    ipaddr string
}

func NewIfAddrAdapter(ifType IfType, ifname , ipaddr string, cli command.Client) *IfAddrAdapter {
    return &IfAddrAdapter{
        Adapter: Adapter{
            client: cli,
        },
        ifType: ifType,
        ifname: ifname,
        ipaddr: ipaddr,
    }
}

func (adpt *IfAddrAdapter) Show(dataType gnmi.GetRequest_DataType) (interface{}, error) {
    conn := adpt.client.Config()
    if conn == nil {
        return nil, swsssdk.ErrConnNotExist
    }

    if data, err := conn.GetAll(swsssdk.CONFIG_DB, []string{IfType_table[int32(adpt.ifType)], adpt.ifname, adpt.ipaddr}); err != nil {
        return nil, err
    } else {
        if adpt.ifType == INTERFACE {
            retval := &sonicpb.AcctonInterface_Interface_InterfaceIpprefixList{}
            for k, v := range data {
                switch k {
                case "scope":
                    switch v {
                    case "local":
                        retval.Scope = sonicpb.AcctonInterface_Interface_InterfaceIpprefixList_SCOPE_local
                    case "global":
                        retval.Scope = sonicpb.AcctonInterface_Interface_InterfaceIpprefixList_SCOPE_global
                    }
                case "family":
                    switch v {
                    case "IPv4":
                        retval.Family = sonicpb.AcctonTypesIpFamily_ACCTONTYPESIPFAMILY_IPv4
                    case "IPv6":
                        retval.Family = sonicpb.AcctonTypesIpFamily_ACCTONTYPESIPFAMILY_IPv6
                    }
                }
            }
            return retval, nil
        } else if adpt.ifType == VLAN_INTERFACE {
            retval := &sonicpb.AcctonVlan_VlanInterface_VlanInterfaceIpprefixList{}
            for k, v := range data {
                switch k {
                case "scope":
                    switch v {
                    case "local":
                        retval.Scope = sonicpb.AcctonVlan_VlanInterface_VlanInterfaceIpprefixList_SCOPE_local
                    case "global":
                        retval.Scope = sonicpb.AcctonVlan_VlanInterface_VlanInterfaceIpprefixList_SCOPE_global
                    }
                case "family":
                    switch v {
                    case "IPv4":
                        retval.Family = sonicpb.AcctonTypesIpFamily_ACCTONTYPESIPFAMILY_IPv4
                    case "IPv6":
                        retval.Family = sonicpb.AcctonTypesIpFamily_ACCTONTYPESIPFAMILY_IPv6
                    }
                }
            }
            return retval, nil
        } else if adpt.ifType == LOOPBACK_INTERFACE {
            retval := &sonicpb.AcctonLoopbackInterface_LoopbackInterface_LoopbackInterfaceIpprefixList{}
            for k, v := range data {
                switch k {
                case "scope":
                    switch v {
                    case "local":
                        retval.Scope = sonicpb.AcctonLoopbackInterface_LoopbackInterface_LoopbackInterfaceIpprefixList_SCOPE_local
                    case "global":
                        retval.Scope = sonicpb.AcctonLoopbackInterface_LoopbackInterface_LoopbackInterfaceIpprefixList_SCOPE_global
                    }
                case "family":
                    switch v {
                    case "IPv4":
                        retval.Family = sonicpb.AcctonTypesIpFamily_ACCTONTYPESIPFAMILY_IPv4
                    case "IPv6":
                        retval.Family = sonicpb.AcctonTypesIpFamily_ACCTONTYPESIPFAMILY_IPv6
                    }
                }
            }
            return retval, nil
        }
    }

    return nil, ErrUnknown
}

func (adpt *IfAddrAdapter) Config(data interface{}, oper OperType) error {
    conn := adpt.client.Config()
    if conn == nil {
        return swsssdk.ErrConnNotExist
    }

    var cmdstr string
    if oper == ADD {
        // 如果存在则跳过重复设置
        if ok, err := conn.HasEntry(IfType_table[int32(adpt.ifType)], []string{adpt.ifname, adpt.ipaddr}); err != nil {
            return err
        } else if ok {
            return nil
        }

        cmdstr = "config interface ip add " + adpt.ifname + " " + adpt.ipaddr
    } else if oper == DEL {
        // 如果不存在则跳过删除
        if ok, err := conn.HasEntry(IfType_table[int32(adpt.ifType)], []string{adpt.ifname, adpt.ipaddr}); err != nil {
            return err
        } else if !ok {
            return nil
        }

        cmdstr = "config interface ip remove " + adpt.ifname + " " + adpt.ipaddr
    } else {
        return ErrInvalidOperType
    }
    return adpt.exec(cmdstr)
}