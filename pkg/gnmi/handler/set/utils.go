package set

import (
    "context"
    "github.com/getlantern/deepcopy"
    "github.com/openconfig/gnmi/proto/gnmi"
    gpb "github.com/openconfig/gnmi/proto/gnmi"
    "github.com/sirupsen/logrus"
    "time"
)

func CreateSetResponse(ctx context.Context, req *gnmi.SetRequest, response []*gnmi.UpdateResult) (*gnmi.SetResponse, error) {
    var prefix gnmi.Path

    err := deepcopy.Copy(&prefix, req.Prefix)
    if err != nil {
        logrus.Errorf("Deep copy struct Prefix failed: %s", err.Error())
        return nil, err
    }

    setResponse := gnmi.SetResponse{
        Prefix: &prefix,
        Response: response,
        Timestamp: time.Now().Unix(),
        Extension: nil,
    }

    return &setResponse, nil
}


func generalPrefixPath(path *gpb.Path) string {
    var gPath string

    for _, pathElem := range path.GetElem() {
        gPath += "/" + pathElem.GetName()
    }
    return gPath
}

