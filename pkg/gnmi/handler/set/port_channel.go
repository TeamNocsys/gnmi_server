package set

import (
    "context"
    "fmt"
    "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/getlantern/deepcopy"
    "github.com/golang/protobuf/proto"
    "github.com/openconfig/gnmi/proto/gnmi"
    "github.com/sirupsen/logrus"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
    "gnmi_server/internal/pkg/swsssdk/helper/config_db"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func PortChannelReplaceHandler(ctx context.Context, req *gnmi.SetRequest, db command.Client) (*gnmi.SetResponse, error) {
    conn := db.Config()
    if conn == nil {
        logrus.Error("Database connection is null")
        return nil, status.Error(codes.Internal, "Database connection is null")
    }

    replaces := req.GetReplace()
    if len(replaces) != 1 {
        message := "Invalid argument for path: " + generalPrefixPath(req.Prefix)
        logrus.Error(message)
        return nil, status.Error(codes.InvalidArgument, message)
    }

    var results []*gnmi.UpdateResult

    replace := replaces[0]
    val := replace.GetVal()
    sonicPortChannel := &sonic.SonicPortchannel{}
    if err := proto.Unmarshal(val.GetBytesVal(), sonicPortChannel); err != nil {
        message := "Can not unmarshal the bytes value to SonicPortChannel " + err.Error()
        logrus.Error(message)
        return nil, status.Error(codes.InvalidArgument, message)
    }

    if err := createPortChannel(sonicPortChannel, conn); err != nil {
        message := fmt.Sprintf("Create port channel failed for %s: %s",
            sonicPortChannel.String(), err.Error())
        logrus.Error(message)
        return nil, status.Error(codes.Internal, message)
    }

    var path gnmi.Path
    if err := deepcopy.Copy(&path, replace.GetPath()); err != nil {
        message := "Deep copy struct path failed: " + err.Error()
        logrus.Error(message)
        return nil, status.Error(codes.Internal, message)
    }

    result := gnmi.UpdateResult{
        Path: replace.GetPath(),
        Op: gnmi.UpdateResult_REPLACE,
    }
    results = append(results, &result)

    if response, err := CreateSetResponse(ctx, req, results); err != nil {
        message := "Create set response failed: " + err.Error()
        logrus.Error(message)
        return nil, status.Error(codes.Internal, message)
    } else {
        return response, nil
    }
}

func PortChannelDeleteHandler(ctx context.Context, req *gnmi.SetRequest, db command.Client) (*gnmi.SetResponse, error) {
    conn := db.Config()
    if conn == nil {
        message := "Datebase connection is null"
        logrus.Error(message)
        return nil, status.Error(codes.Internal, message)
    }

    deletes := req.GetDelete()
    if len(deletes) != 1 {
        message := "Invalid argument for path: " + generalPrefixPath(req.Prefix)
        logrus.Error(message)
        return nil, status.Error(codes.InvalidArgument, message)
    }

    var results []*gnmi.UpdateResult

    delete := deletes[0]
    name := delete.GetTarget()
    if err := deletePortChannel(name, conn); err != nil {
        message := fmt.Sprintf("Delete port channel failed for %s: %s", name, err.Error())
        logrus.Error(message)
        return nil, status.Error(codes.Internal, message)
    }

    result := gnmi.UpdateResult {
        Path: delete,
        Op: gnmi.UpdateResult_DELETE,
    }

    results = append(results, &result)

    if response, err := CreateSetResponse(ctx, req, results); err != nil {
        logrus.Errorf("Delete portchannel-%s failed: %s", name, err.Error())
        return nil, status.Error(codes.Internal, "Delete portchannel failed")
    } else {
        return response, nil
    }
}

func createPortChannel(sonicPortChannel *sonic.SonicPortchannel, conn *swsssdk.ConfigDBConnector) error {
    portChannel := sonicPortChannel.GetPortchannel()

    for _, key := range portChannel.GetPortchannelList() {
        name := key.GetPortchannelName()
        minLinks := key.GetPortchannelList().GetMinLinks()
        mtu := key.GetPortchannelList().GetMtu()

        memberList, err := conn.GetAllByPattern(swsssdk.CONFIG_DB, []string{config_db.PORTCHANNEL_MEMBER_TABLE, name, "*"})
        if err != nil {
            return err
        } else {
            for entry, _ := range memberList {
                keys := splitConfigDBKey(entry)
                if len(keys) == 3 {
                    if _, err := conn.SetEntry(config_db.PORTCHANNEL_MEMBER_TABLE, []string{keys[1], keys[2]}, nil); err != nil {
                        return err
                    }
                }
            }
        }

        if entry, err := conn.GetEntry(config_db.PORTCHANNEL_TABLE, name); err != nil {
            return err
        } else if entry != nil {
            conn.SetEntry(config_db.PORTCHANNEL_TABLE, name, nil)
        }

        values := make(map[string]interface{})
        values["admin_status"] = "up"
        values["min_links"] = minLinks.String()
        values["mtu"] = mtu.String()
        if _, err := conn.SetEntry(config_db.PORTCHANNEL_TABLE, name, values); err != nil {
            return err
        }

        members := key.GetPortchannelList().GetMembers()
        for _, member := range members {
            memberName := member.GetValue()
            if _, err := conn.SetEntry(config_db.PORTCHANNEL_MEMBER_TABLE, []string{name, memberName},
                make(map[string]interface{})); err != nil {
                return err
            }
        }
    }

    return nil
}

func deletePortChannel(name string, conn *swsssdk.ConfigDBConnector) error {
    if memberList, err := conn.GetAllByPattern(swsssdk.CONFIG_DB, []string{config_db.PORTCHANNEL_MEMBER_TABLE, name, "*"}); err != nil {
        return err
    } else {
        for entry, _ := range memberList {
            keys := splitConfigDBKey(entry)
            if len(keys) == 3 {
                if _, err := conn.SetEntry(config_db.PORTCHANNEL_MEMBER_TABLE, []string{keys[1], keys[2]}, nil); err != nil {
                    return err
                }
            }
        }
    }

    if entry, err := conn.GetEntry(config_db.PORTCHANNEL_TABLE, name); err != nil {
        return err
    } else if entry != nil {
        conn.SetEntry(config_db.PORTCHANNEL_TABLE, name, nil)
    }

    return nil
}

func PortChannelMemberUpdateHandler(ctx context.Context, req *gnmi.SetRequest, db command.Client) (*gnmi.SetResponse, error) {
    conn := db.Config()
    if conn == nil {
        message := "Database connection is null"
        logrus.Error(message)
        return nil, status.Error(codes.Internal, message)
    }

    updates := req.GetUpdate()
    if len(updates) != 1 {
        message := "Invalid argument for path: " + generalPrefixPath(req.Prefix)
        logrus.Error(message)
        return nil, status.Error(codes.InvalidArgument, message)
    }

    var results []*gnmi.UpdateResult

    update := updates[0]
    val := update.GetVal()
    sonicPortChannel := &sonic.SonicPortchannel{}
    if err := proto.Unmarshal(val.GetBytesVal(), sonicPortChannel); err != nil {
        message := "Can not unmarshal the bytes value to SonicPortChannel: " + err.Error()
        logrus.Error(message)
        return nil, status.Error(codes.Internal, message)
    }

    if err := portChannelAddMember(sonicPortChannel, conn); err != nil {
        logrus.Errorf("Add portchannel member failed: %s", err.Error())
        return nil, status.Error(codes.Internal, "Add portchannel member failed")
    }

    result := gnmi.UpdateResult {
        Path: update.GetPath(),
        Op: gnmi.UpdateResult_UPDATE,
    }
    results = append(results, &result)

    if response, err := CreateSetResponse(ctx, req, results); err != nil {
        logrus.Errorf("Create response failed: %s", err.Error())
        return nil, status.Error(codes.Internal, "Create response failed")
    } else {
        return response, nil
    }
}

func PortChannelMemberDeleteHandler(ctx context.Context, req *gnmi.SetRequest, db command.Client)(*gnmi.SetResponse, error) {
    conn := db.Config()
    if conn == nil {
        message := "Database connection is null"
        logrus.Error(message)
        return nil, status.Error(codes.Internal, message)
    }

    deletes := req.GetDelete()
    if len(deletes) != 1 {
        message := "Invalid argument for path: " + generalPrefixPath(req.Prefix)
        logrus.Error(message)
        return nil, status.Error(codes.InvalidArgument, message)
    }

    var results []*gnmi.UpdateResult

    delete := deletes[0]
    memberName := delete.GetTarget()

    if err := portChannelDelMember(memberName, conn); err != nil {
        logrus.Errorf("Delete port-%s from portchannel failed: %s", memberName, err.Error())
        return nil, status.Error(codes.InvalidArgument, "Delete member from portchannel failed")
    }

    result := gnmi.UpdateResult {
        Path: delete,
        Op: gnmi.UpdateResult_DELETE,
    }
    results = append(results, &result)

    if response, err := CreateSetResponse(ctx, req, results); err != nil {
        logrus.Errorf("Create response failed: %s", err.Error())
        return nil, status.Error(codes.Internal, "Create response failed")
    } else {
        return response, nil
    }
}

func portChannelAddMember(sonicPortChannel *sonic.SonicPortchannel, conn *swsssdk.ConfigDBConnector) error {
    portChannel := sonicPortChannel.GetPortchannel()

    values := make(map[string]interface{})
    values["NULL"] = "NULL"

    for _, key := range portChannel.GetPortchannelList() {
        name := key.GetPortchannelName()
        members := key.GetPortchannelList().GetMembers()

        for _, member := range members {
            memberName := member.String()
            _, err := conn.SetEntry(config_db.PORTCHANNEL_MEMBER_TABLE, []string{name, memberName}, values)
            if err != nil {
                return err
            }
        }
    }

    return nil
}

func portChannelDelMember(member string, conn *swsssdk.ConfigDBConnector) error {
    values := make(map[string]interface{})
    values["NULL"] = "NULL"

    if memberTable, err := conn.GetAllByPattern(config_db.PORTCHANNEL_MEMBER_TABLE, []string{"*", member}); err != nil {
        return err
    } else {
        for entry := range memberTable {
            keys := splitConfigDBKey(entry)
            if len(keys) == 3 {
                if _, err := conn.SetEntry(config_db.PORTCHANNEL_MEMBER_TABLE, []string{keys[1], keys[2]}, nil); err != nil {
                    return err
                }
            }
        }
    }

    return nil
}
