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

type MirrorAdapter struct {
    Adapter
    name string
}

func NewMirrorAdapter(name string, cli command.Client) *MirrorAdapter {
    return &MirrorAdapter{
        Adapter: Adapter{
            client: cli,
        },
        name:  name,
    }
}

func (adpt *MirrorAdapter) Show(dataType gnmi.GetRequest_DataType) (*sonicpb.NocsysMirrorSession_MirrorSession_MirrorSessionList, error) {
    conn := adpt.client.State()
    if conn == nil {
        return nil, swsssdk.ErrConnNotExist
    }

    if data, err := conn.GetAll(swsssdk.APPL_DB, []string{"MIRROR_SESSION_TABLE", adpt.name}); err != nil {
        return nil, err
    } else {
        retval := &sonicpb.NocsysMirrorSession_MirrorSession_MirrorSessionList{}
        for k, v := range data {
            switch k {
            case "status":
                switch v {
                case "active":
                    retval.Status = sonicpb.NocsysMirrorSession_MirrorSession_MirrorSessionList_STATUS_active
                case "inactive":
                    retval.Status = sonicpb.NocsysMirrorSession_MirrorSession_MirrorSessionList_STATUS_inactive
                }
            case "src_ip":
                retval.SrcIp = &ywrapper.StringValue{Value: v}
            case "dst_ip":
                retval.DstIp = &ywrapper.StringValue{Value: v}
            case "gre_type":
                // mellanox => 0x8949 other=> 0x88be
                if i, err := strconv.ParseUint(v, 10, 64); err != nil {
                    return nil, err
                } else {
                    retval.GreType = &ywrapper.UintValue{Value: i}
                }
            case "dscp":
                if i, err := strconv.ParseUint(v, 10, 64); err != nil {
                    return nil, err
                } else {
                    retval.Dscp = &ywrapper.UintValue{Value: i}
                }
            case "ttl":
                if i, err := strconv.ParseUint(v, 10, 64); err != nil {
                    return nil, err
                } else {
                    retval.Ttl = &ywrapper.UintValue{Value: i}
                }
            case "queue":
                if i, err := strconv.ParseUint(v, 10, 64); err != nil {
                    return nil, err
                } else {
                    retval.Queue = &ywrapper.UintValue{Value: i}
                }
            case "dst_port":
                retval.DstPort = &ywrapper.StringValue{Value: v}
            case "src_port":
                retval.SrcPort = &ywrapper.StringValue{Value: v}
            case "direction":
                switch strings.ToUpper(v) {
                case "RX":
                    retval.Direction = sonicpb.NocsysMirrorSession_MirrorSession_MirrorSessionList_DIRECTION_RX
                case "TX":
                    retval.Direction = sonicpb.NocsysMirrorSession_MirrorSession_MirrorSessionList_DIRECTION_TX
                case "BOTH":
                    retval.Direction = sonicpb.NocsysMirrorSession_MirrorSession_MirrorSessionList_DIRECTION_BOTH
                }
            case "type":
                switch strings.ToUpper(v) {
                case "SPAN":
                    retval.Type = sonicpb.NocsysMirrorSession_MirrorSession_MirrorSessionList_TYPE_SPAN
                case "ERSPAN":
                    retval.Type = sonicpb.NocsysMirrorSession_MirrorSession_MirrorSessionList_TYPE_ERSPAN
                }
            }
        }
        return retval, nil
    }
}

func (adpt *MirrorAdapter) Config(data *sonicpb.NocsysMirrorSession_MirrorSession_MirrorSessionList, oper OperType) error {
    var cmdstr string
    if oper == ADD {
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
        cmdstr = "config mirror_session remove " + adpt.name
    }

    return adpt.exec(cmdstr)
}