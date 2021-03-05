package cmd

import (
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/gnmi/proto/gnmi"
    "github.com/openconfig/ygot/proto/ywrapper"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
    "strconv"
    "strings"
)

type VlanAdapter struct {
    Adapter
    name string
}

func NewVlanAdapter(name string, cli command.Client) *VlanAdapter {
    return &VlanAdapter{
        Adapter: Adapter{
            client: cli,
        },
        name:  name,
    }
}

func (adpt *VlanAdapter) Show(dataType gnmi.GetRequest_DataType) (*sonicpb.NocsysVlan_Vlan_VlanList, error) {
    conn := adpt.client.Config()
    if conn == nil {
        return nil, swsssdk.ErrConnNotExist
    }

    if data, err := conn.GetEntry("VLAN", adpt.name); err != nil {
        return nil, err
    } else {
        retval := &sonicpb.NocsysVlan_Vlan_VlanList{}
        for k, v := range data {
            switch k {
            case "description":
                retval.Description = &ywrapper.StringValue{Value: v.(string)}
            case "dhcp_servers":
                for _, server := range v.([]string) {
                    retval.DhcpServers = append(retval.DhcpServers, &ywrapper.StringValue{Value: server})
                }
            case "mtu":
                if i, err := strconv.ParseUint(v.(string), 10, 64); err != nil {
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
            }
        }
        return retval, nil
    }
}

func (adpt *VlanAdapter) Config(data *sonicpb.NocsysVlan_Vlan_VlanList, oper OperType) error {
    if oper == ADD || oper == UPDATE {
        conn := adpt.client.Config()
        if conn == nil {
            return swsssdk.ErrConnNotExist
        }
        if ok, err := conn.HasEntry("VLAN", adpt.name); err != nil {
            return err
        } else if !ok {
            cmdstr := "config vlan add " + strings.TrimLeft(adpt.name, "Vlan")
            if err := adpt.exec(cmdstr); err != nil {
                return err
            }
        }

        if data.Description != nil {
            cmdstr := "config interface description " + adpt.name + " " +  data.Description.Value
            if err := adpt.exec(cmdstr); err != nil {
                return err
            }
        }

        if data.Mtu != nil {
            cmdstr := "config interface mtu " + adpt.name + " " +  strconv.FormatUint(data.Mtu.Value, 10)
            if err := adpt.exec(cmdstr); err != nil {
                return err
            }
        }

        if data.AdminStatus != sonicpb.NocsysTypesAdminStatus_NOCSYSTYPESADMINSTATUS_UNSET {
            var cmdstr string
            if data.AdminStatus == sonicpb.NocsysTypesAdminStatus_NOCSYSTYPESADMINSTATUS_up {
                cmdstr = "config interface startup " + adpt.name
            } else if data.AdminStatus == sonicpb.NocsysTypesAdminStatus_NOCSYSTYPESADMINSTATUS_down {
                cmdstr = "config interface shutdown " + adpt.name
            }
            if err := adpt.exec(cmdstr); err != nil {
                return err
            }
        }

        if data.DhcpServers != nil {
            // TODO: cmdstr = "config interface ip dhcp-relay add " + vid
        }
    } else if oper == DEL {
        cmdstr := "config vlan del " + strings.TrimLeft(adpt.name, "Vlan")
        if err := adpt.exec(cmdstr); err != nil {
            return err
        }
    }
    return nil
}