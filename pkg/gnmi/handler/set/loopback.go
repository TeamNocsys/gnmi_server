package set

import (
    "context"
    "errors"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
    "gnmi_server/internal/pkg/swsssdk/helper/config_db"
    "strconv"
    "strings"

    "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/getlantern/deepcopy"
    "github.com/golang/protobuf/proto"
    "github.com/openconfig/gnmi/proto/gnmi"
    "github.com/sirupsen/logrus"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

const LoopbackPrefix string = "Loopback"
const LoopbackPrefixLen int = len(LoopbackPrefix)
const TotalLoopbackPrefixLen int = 11

func LoopbackUpdateHandler(ctx context.Context, req *gnmi.SetRequest, db command.Client) (*gnmi.SetResponse, error) {
    conn := db.Config()
    if conn == nil {
        message := "Database connection is NULL"
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
    elements := update.GetPath().GetElem()
    if len(elements) != 0 {
        logrus.Errorf("Unimplemented request: " + req.String())
        return nil, status.Error(codes.Unimplemented, "Unimplemented request")
    }

    val := update.GetVal()
    sonicLoopback := &sonic.SonicLoopbackInterface{}
    err := proto.Unmarshal(val.GetBytesVal(), sonicLoopback)
    if err != nil {
        message := "Can not unmarshal the bytes value to SonicLoopbackInterface"
        logrus.Error(message)
        return nil, status.Error(codes.Unimplemented, message)
    }

    if err := addLoopback(sonicLoopback, conn); err != nil {
        logrus.Errorf("Add loopback interface failed: %s", err.Error())
        return nil, status.Error(codes.Internal, "Add loopback interface failed")
    }

    var path gnmi.Path
    err = deepcopy.Copy(&path, update.GetPath())
    if err != nil {
        message := "Deep copy struct path failed: " + err.Error()
        logrus.Errorf(message)
        return nil, status.Error(codes.Internal, message)
    }

    result := gnmi.UpdateResult{
        Path: &path,
        Op:   gnmi.UpdateResult_UPDATE,
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

func LoopbackDeleteHandler(ctx context.Context, req *gnmi.SetRequest, db command.Client) (*gnmi.SetResponse, error) {
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
    if len(delete.GetElem()) != 0 {
        logrus.Error("Unimplemented request: " + req.String())
        return nil, status.Error(codes.Unimplemented, "Unimplemented request")
    }

    name := delete.GetTarget()
    if err := deleteLoopback(name, conn); err != nil {
        logrus.Errorf("Delete loopback interface failed: %s", err.Error())
        return nil, status.Error(codes.Internal, "Delete loopback interface failed")
    }

    result := gnmi.UpdateResult{
        Path: delete,
        Op:   gnmi.UpdateResult_DELETE,
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

func addLoopback(sonicLoopback *sonic.SonicLoopbackInterface, conn *swsssdk.ConfigDBConnector) error {
    loopback := sonicLoopback.GetLoopbackInterface()
    ipPrefixList := loopback.GetLoopbackInterfaceIpprefixList()

    loopbackTable, err := conn.GetTable(config_db.LOOPBACK_INTERFACE_TABLE)
    if err != nil {
        return err
    }

    values := make(map[string]interface{})

    for _, key := range ipPrefixList {
        name := key.GetLoopbackInterfaceName()
        if !verifyLoopbackName(name) {
            return errors.New("invalid loopback name")
        }

        ipPrefix := key.GetIpPrefix()
        //ipPrefixList := key.GetLoopbackInterfaceIpprefixList()

        if _, ok := loopbackTable[name]; !ok {
            logrus.Warnf("%s already exist in configuration database", name)
        } else {
            if _, err := conn.SetEntry(config_db.LOOPBACK_INTERFACE_TABLE, name, values); err != nil {
                return err
            }
        }

        if _, err := conn.SetEntry(config_db.LOOPBACK_INTERFACE_TABLE, []string{name, ipPrefix}, values); err != nil {
            return err
        }
    }

    return nil
}

func deleteLoopback(name string, conn *swsssdk.ConfigDBConnector) error {
    if !verifyLoopbackName(name) {
        return errors.New("invalid loopback name")
    }

    loopbackTables, err := conn.GetAllByPattern(swsssdk.CONFIG_DB, []string{config_db.LOOPBACK_INTERFACE_TABLE, name})
    if err != nil {
        return err
    }

    for entry := range loopbackTables {
        keys := splitConfigDBKey(entry)
        if len(keys) == 3 && keys[1] == name {
            if _, err := conn.SetEntry(config_db.LOOPBACK_INTERFACE_TABLE, []string{keys[1], keys[2]}, nil); err != nil {
                return err
            }
        }
    }

    if _, err := conn.SetEntry(config_db.LOOPBACK_INTERFACE_TABLE, name, nil); err != nil {
        return err
    }

    return nil
}

func verifyLoopbackName(name string) bool {
    if len(name) > TotalLoopbackPrefixLen {
        return false
    }

    if !strings.HasPrefix(name, LoopbackPrefix) {
        return false
    }

    str := name[LoopbackPrefixLen:]
    if _, err := strconv.Atoi(str); err != nil {
        return false
    }

    return true
}
