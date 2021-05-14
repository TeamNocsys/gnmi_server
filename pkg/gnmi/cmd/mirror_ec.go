// +build ec

package cmd

import (
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "gnmi_server/internal/pkg/swsssdk"
)

func (adpt *MirrorAdapter) Config(data *sonicpb.NocsysMirrorSession_MirrorSession_MirrorSessionList, oper OperType) error {
    var cmdstr string
    if oper == ADD {
        conn := adpt.client.Config()
        if conn == nil {
            return swsssdk.ErrConnNotExist
        }
        if ok, err := conn.HasEntry("MIRROR_SESSION", adpt.name); err != nil {
            return err
        } else if ok {
            return nil
        }

        cmdstr = "config mirror_session add"

        if data.Type == sonicpb.NocsysMirrorSession_MirrorSession_MirrorSessionList_TYPE_SPAN {
            cmdstr += " span " + adpt.name
        } else {
            return ErrUnknown
        }

        if data.DstPort != nil && data.SrcPort != nil {
            cmdstr += " " + data.DstPort.Value + " " + data.SrcPort.Value
        }

        if data.Direction == sonicpb.NocsysMirrorSession_MirrorSession_MirrorSessionList_DIRECTION_RX {
            cmdstr += " rx"
        } else if data.Direction == sonicpb.NocsysMirrorSession_MirrorSession_MirrorSessionList_DIRECTION_TX {
            cmdstr += " tx"
        } else {
            cmdstr += " both"
        }
    } else if oper == DEL {
        conn := adpt.client.Config()
        if conn == nil {
            return swsssdk.ErrConnNotExist
        }
        if ok, err := conn.HasEntry("MIRROR_SESSION", adpt.name); err != nil {
            return err
        } else if !ok {
            return nil
        }

        cmdstr = "config mirror_session remove " + adpt.name
    }

    return adpt.exec(cmdstr)
}