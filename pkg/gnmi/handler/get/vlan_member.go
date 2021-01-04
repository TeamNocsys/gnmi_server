package get

import (
    "context"
    "errors"
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/gnmi/proto/gnmi"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
    "gnmi_server/internal/pkg/swsssdk/helper/config_db"
    "gnmi_server/pkg/gnmi/handler"
    handler_utils "gnmi_server/pkg/gnmi/handler/utils"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    "strings"
)

func VlanMemberHandler(ctx context.Context, r *gnmi.GetRequest, db command.Client) (*gnmi.GetResponse, error) {
    conn := db.Config()
    if conn == nil {
        return nil, status.Error(codes.Internal, "")
    }
    // 获取指定VLAN Member或全部VLAN Member
    kvs := handler.FetchPathKey(r)
    spec := []string{}
    if v, ok := kvs["vlan-name"]; ok {
        spec = append(spec, v)
    } else {
        spec = append(spec, "*")
    }
    if v, ok := kvs["port"]; ok {
        spec = append(spec, v)
    } else {
        spec = append(spec, "*")
    }

    infos, err := conn.GetAllByPattern(swsssdk.CONFIG_DB, append([]string{config_db.VLAN_MEMBER_TABLE}, spec...))
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }
    sv := &sonicpb.SonicVlan{
        VlanMember: &sonicpb.SonicVlan_VlanMember{},
    }
    s := swsssdk.Config().GetDBSeparator(swsssdk.CONFIG_DB)
    for hash, info := range infos {
        keys := strings.Split(hash, s)
        if len(keys) != 3 {
            continue
        }
        v, err := getVlanMemberList(info)
        if err != nil {
            return nil, status.Error(codes.Internal, err.Error())
        }
        sv.VlanMember.VlanMemberList = append(sv.VlanMember.VlanMemberList, &sonicpb.SonicVlan_VlanMember_VlanMemberListKey{
            VlanName: keys[1],
            Port: keys[2],
            VlanMemberList: v,
        })
    }

    response, err := handler_utils.CreateGetResponse(ctx, r, sv)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return response, nil
}

func getVlanMemberList(info map[string]string) (*sonicpb.SonicVlan_VlanMember_VlanMemberList, error) {
    r := &sonicpb.SonicVlan_VlanMember_VlanMemberList{}

    if v, ok := info[config_db.VLAN_MEMBER_TAGGING_MODE]; ok {
        if v == "tagged" {
            r.TaggingMode = sonicpb.SonicVlanVlanTaggingMode_SONICVLANVLANTAGGINGMODE_tagged
        } else {
            r.TaggingMode = sonicpb.SonicVlanVlanTaggingMode_SONICVLANVLANTAGGINGMODE_untagged
        }
    } else  {
        return nil, errors.New("missing " + config_db.VLAN_TABLE_DESCRIPTION + " field")
    }

    return r, nil
}