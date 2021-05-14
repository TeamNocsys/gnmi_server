package cmd

import (
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/gnmi/proto/gnmi"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
)

type LagMemberAdapter struct {
    Adapter
    name string
    ifname string
}

func NewLagMemberAdapter(name, ifname string, cli command.Client) *LagMemberAdapter {
    return &LagMemberAdapter{
        Adapter: Adapter{
            client: cli,
        },
        name:    name,
        ifname:  ifname,
    }
}

func (adpt *LagMemberAdapter) Show(dataType gnmi.GetRequest_DataType) (*sonicpb.NocsysPortchannel_PortchannelMember_PortchannelMemberList, error) {
    conn := adpt.client.Config()
    if conn == nil {
        return nil, swsssdk.ErrConnNotExist
    }

    if _, err := conn.HasEntry("PORTCHANNEL_MEMBER", []string{adpt.name, adpt.ifname}); err != nil {
        return nil, err
    } else {
        return &sonicpb.NocsysPortchannel_PortchannelMember_PortchannelMemberList{}, nil
    }
}