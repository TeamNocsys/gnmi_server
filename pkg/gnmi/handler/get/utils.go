package get

import (
    "context"
    "github.com/getlantern/deepcopy"
    "github.com/golang/glog"
    "github.com/openconfig/gnmi/proto/gnmi"
    "time"
)

func createResponse(ctx context.Context, req *gnmi.GetRequest, bytes []byte) (*gnmi.GetResponse, error) {
    var prefix gnmi.Path
    var path gnmi.Path

    err := deepcopy.Copy(&prefix, req.Prefix)
    if err != nil {
        glog.Errorf("deep copy struct Prefix failed: %s", err.Error())
        return nil, err
    }

    err = deepcopy.Copy(&path, req.Path[0])
    if err != nil {
        glog.Errorf("deep copy struct Path failed: %s", err.Error())
        return nil, err
    }

    notification := gnmi.Notification{
        Timestamp: time.Now().Unix(),
        Prefix:    &prefix,
        Update: []*gnmi.Update{
            &gnmi.Update{
            Path: &path,
            Val: &gnmi.TypedValue{
                Value: &gnmi.TypedValue_BytesVal{
                    BytesVal: bytes,
                },
            },
            },
        },
    }

    response := &gnmi.GetResponse{}
    response.Notification = append(response.Notification, &notification)
    return response, nil
}
