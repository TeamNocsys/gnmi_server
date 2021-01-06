package set

import (
    "context"
    "fmt"
    "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/getlantern/deepcopy"
    "github.com/openconfig/gnmi/proto/gnmi"
    "github.com/sirupsen/logrus"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
    "gnmi_server/internal/pkg/swsssdk/helper/config_db"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    "google.golang.org/protobuf/proto"
)

func PortListHandler(ctx context.Context, req *gnmi.SetRequest, db command.Client) (*gnmi.SetResponse, error) {
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

    update := updates[0]
    elements := update.GetPath().GetElem()
    if len(elements) != 1 && elements[0].GetName() != "admin-status" {
        message := "Unimplemented request: " + req.String()
        logrus.Error(message)
        return nil, status.Error(codes.Unimplemented, "Unimplemented request")
    }

    val := update.GetVal()

    sonicPort := &sonic.SonicPort{}
    if err := proto.Unmarshal(val.GetBytesVal(), sonicPort); err != nil {
        logrus.Error("Can not unmarshal the bytes value to SonicPort, request is: " + req.String())
        return nil, status.Error(codes.InvalidArgument, "Can't unmarshal bytes value to SonicPort")
    }

    port := sonicPort.GetPort()
    portListKeys := port.GetPortList()
    for _, portListKey := range portListKeys {
        portName := portListKey.GetPortName()
        portList := portListKey.GetPortList()
        adminStatus := portList.GetAdminStatus()

        var newStatus = ""
        if adminStatus == sonic.SonicPortAdminStatus_SONICPORTADMINSTATUS_up {
            newStatus= "up"
        } else {
            newStatus= "down"
        }

        oldValues, err := conn.GetEntry(config_db.PORT_TABLE, []string{config_db.PORT_TABLE, portName})
        if err == swsssdk.ErrDatabaseNotExist {
            message := "Configuration database is not exists"
            logrus.Errorf(message)
            return nil, status.Error(codes.Internal, message)
        }

        oldStatus, ok := oldValues["admin_status"]
        if ok == true && oldStatus == newStatus {
            continue
        }

        values := make(map[string]interface{})
        values["admin_status"] = "up"

        result, err := conn.ModEntry(config_db.PORT_TABLE, []string{config_db.PORT_TABLE, portName}, values)
        if err != nil {
            message := fmt.Sprintf("set admin_status failed for port-%s", portName)
            return nil, status.Error(codes.Internal, message)
        } else if result == false {
            message := fmt.Sprintf("port-%s is not exists", portName)
            return nil, status.Error(codes.NotFound, message)
        }
    }

    var path gnmi.Path
    if err := deepcopy.Copy(&path, update.GetPath()); err != nil {
        message := "Deep copy struct path failed"
        logrus.Errorf(message + ": " + err.Error())
        return nil, status.Error(codes.Internal, message)
    }

    result := gnmi.UpdateResult{
        Path: &path,
        Op: gnmi.UpdateResult_UPDATE,
    }

    var results []*gnmi.UpdateResult
    results = append(results, &result)

    if response, err := CreateSetResponse(ctx, req, results); err != nil {
        message := fmt.Sprintf("Create set response failed: %s", err.Error())
        logrus.Error(message)
        return nil, status.Error(codes.Internal, message)
    } else {
        return response, nil
    }
}

