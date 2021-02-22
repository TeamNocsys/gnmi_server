package get

import (
    "context"
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/gnmi/proto/gnmi"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
    "gnmi_server/internal/pkg/swsssdk/helper"
    "gnmi_server/pkg/gnmi/handler"
    handler_utils "gnmi_server/pkg/gnmi/handler/utils"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func AclRuleHandler(ctx context.Context, r *gnmi.GetRequest, db command.Client) (*gnmi.GetResponse, error) {
    conn := db.State()
    if conn == nil {
        return nil, status.Error(codes.Internal, "")
    }
    // 获取指定ACL Rule或全部ACL Rule
    kvs := handler.FetchPathKey(r)
    var spec []string
    if v, ok := kvs["table-name"]; ok {
        spec = append(spec, v)
    } else {
        spec = append(spec, "*")
    }
    if v, ok := kvs["rule-name"]; ok {
        spec = append(spec, v)
    } else {
        spec = append(spec, "*")
    }

    sa := &sonicpb.NocsysAcl{
        AclRule: &sonicpb.NocsysAcl_AclRule{},
    }
    if hkeys, err := conn.GetKeys(swsssdk.APPL_DB, append([]string{"ACL_RULE"}, spec...)); err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    } else {
        for _, hkey := range hkeys {
            keys := conn.SplitKeys(swsssdk.APPL_DB, hkey)
            c := helper.AclRule{
                Keys: keys,
                Client: db,
                Data: nil,
            }
            if err := c.LoadFromDB(); err != nil {
                return nil, status.Errorf(codes.Internal, err.Error())
            }
            sa.AclRule.AclRuleList = append(sa.AclRule.AclRuleList,
                &sonicpb.NocsysAcl_AclRule_AclRuleListKey{
                    TableName: keys[0],
                    RuleName: keys[1],
                    AclRuleList: c.Data,
                })
        }
    }

    response, err := handler_utils.CreateGetResponse(ctx, r, sa)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return response, nil
}

func AclTableHandler(ctx context.Context, r *gnmi.GetRequest, db command.Client) (*gnmi.GetResponse, error) {
    conn := db.State()
    if conn == nil {
        return nil, status.Error(codes.Internal, "")
    }
    // 获取指定ACL Table或全部ACL Table
    kvs := handler.FetchPathKey(r)
    spec := "*"
    if v, ok := kvs["table-name"]; ok {
        spec = v
    }

    sa := &sonicpb.NocsysAcl{
        AclTable: &sonicpb.NocsysAcl_AclTable{},
    }
    if hkeys, err := conn.GetKeys(swsssdk.APPL_DB, []string{"ACL_TABLE", spec}); err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    } else {
        for _, hkey := range hkeys {
            keys := conn.SplitKeys(swsssdk.APPL_DB, hkey)
            c := helper.AclTable{
                Key: keys[0],
                Client: db,
                Data: nil,
            }
            if err := c.LoadFromDB(); err != nil {
                return nil, status.Errorf(codes.Internal, err.Error())
            }
            sa.AclTable.AclTableList = append(sa.AclTable.AclTableList,
                &sonicpb.NocsysAcl_AclTable_AclTableListKey{
                    TableName: keys[0],
                    AclTableList: c.Data,
                })
        }
    }

    response, err := handler_utils.CreateGetResponse(ctx, r, sa)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return response, nil
}