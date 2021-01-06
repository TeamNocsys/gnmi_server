package set

import (
    "context"
    "fmt"
    "github.com/openconfig/gnmi/proto/gnmi"
    "github.com/sirupsen/logrus"
    "gnmi_server/cmd/command"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    "os/exec"
)

func SonicUpdateHandler(ctx context.Context, req *gnmi.SetRequest, db command.Client) (*gnmi.SetResponse, error) {
    updates := req.GetUpdate()
    if len(updates) != 1 {
        message := "Invalid argument for path: " + generalPrefixPath(req.Prefix)
        logrus.Error(message)
        return nil, status.Error(codes.InvalidArgument, message)
    }

    var results []*gnmi.UpdateResult

    update := updates[0]
    if len(update.GetPath().GetElem()) != 0 {
        logrus.Errorf("Unimplemented requesst: " + req.String())
        return nil, status.Error(codes.Unimplemented, "Unimplemented request")
    }

    arg := update.GetVal().GetStringVal()

    cmd := exec.Command("sonic-cfggen", "-a", arg, "--write-to-db")
    _, err := cmd.Output()
    if err != nil {
        logrus.Errorf("Execute sonic-cfggen  -a %s --write-to-db failed", arg)
        return nil, status.Error(codes.Internal, "Execute sonic-cfggen failed")
    }

    //var path gnmi.Path
    //err = deepcopy.Copy(&path, update.GetPath())
    //if err != nil {
    //    logrus.Errorf("Deep copy struct path failed: %s", err.Error())
    //    return nil, err
    //}

    result := gnmi.UpdateResult{
    //    Path: &path,
        Path: update.Path,
        Op: gnmi.UpdateResult_UPDATE,
    }

    results = append(results, &result)

    if response, err := CreateSetResponse(ctx, req, results); err != nil {
        message := fmt.Sprintf("Create set response failed: %s", err.Error())
        logrus.Error(message)
        return nil, status.Error(codes.Internal, message)
    } else {
        return response, nil
    }
}
