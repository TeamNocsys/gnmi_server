// +build debug

package handler

import (
    "context"
    "encoding/json"
    "github.com/getlantern/deepcopy"
    "github.com/golang/protobuf/proto"
    "github.com/openconfig/gnmi/proto/gnmi"
    "github.com/sirupsen/logrus"
    "time"
)

func CreateResponse(ctx context.Context, req *gnmi.GetRequest, message proto.Message) (*gnmi.GetResponse, error) {
    var prefix gnmi.Path
    var path gnmi.Path

    err := deepcopy.Copy(&prefix, req.Prefix)
    if err != nil {
        logrus.Errorf("deep copy struct Prefix failed: %s", err.Error())
        return nil, err
    }

    err = deepcopy.Copy(&path, req.Path[0])
    if err != nil {
        logrus.Errorf("deep copy struct Path failed: %s", err.Error())
        return nil, err
    }

    bytes, err := json.Marshal(message)
    if err != nil {
        logrus.Errorf("serialization failed: %s", err.Error())
        return nil, err
    }

    notification := gnmi.Notification{
        Timestamp: time.Now().Unix(),
        Prefix:    &prefix,
        Update: []*gnmi.Update{
            &gnmi.Update{
                Path: &path,
                Val: &gnmi.TypedValue{
                    Value: &gnmi.TypedValue_JsonVal{
                        JsonVal: bytes,
                    },
                },
            },
        },
    }

    response := &gnmi.GetResponse{}
    response.Notification = append(response.Notification, &notification)
    return response, nil
}