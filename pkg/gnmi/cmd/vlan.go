package cmd

import (
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/gnmi/proto/gnmi"
    "github.com/openconfig/ygot/proto/ywrapper"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
    "strconv"
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