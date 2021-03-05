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

func (adpt *LagMemberAdapter) Config(data *sonicpb.NocsysPortchannel_PortchannelMember_PortchannelMemberList, oper OperType) error {
    cmdstr := "config portchannel member"
    if oper == ADD || oper == UPDATE {
        conn := adpt.client.Config()
        if conn == nil {
            return swsssdk.ErrConnNotExist
        }
        if ok, err := conn.HasEntry("PORTCHANNEL_MEMBER", []string{adpt.name, adpt.ifname}); err != nil {
            return err
        } else if ok {
            return nil
        }

        cmdstr = " add "
    } else if oper == DEL {
        cmdstr = " del "
    } else {
        return ErrInvalidOperType
    }
    conn := adpt.client.Config()
    if conn == nil {
        return swsssdk.ErrConnNotExist
    }
    if ok, err := conn.HasEntry("PORTCHANNEL_MEMBER", []string{adpt.name, adpt.ifname}); err != nil {
        return err
    } else if !ok {
        return nil
    }

    cmdstr += " " + adpt.name + " " + adpt.ifname
    return adpt.exec(cmdstr)
}