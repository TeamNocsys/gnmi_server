package utils

import (
    "github.com/golang/protobuf/jsonpb"
    "github.com/golang/protobuf/proto"
    gpb "github.com/openconfig/gnmi/proto/gnmi"
    "github.com/sirupsen/logrus"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

var (
        ErrProtobufType = "Only bytes and jsonietf types are supported"
)

func UnmarshalGpbValue(in_value *gpb.TypedValue, m proto.Message) error {
    if in_json := in_value.GetJsonIetfVal(); in_json != nil {
        if err := jsonpb.UnmarshalString(string(in_json), m); err != nil {
            return err
        }

        logrus.Tracef("PARSEJSON|json-%s", in_json)
    } else if bytes := in_value.GetBytesVal(); bytes != nil {
        if err := proto.Unmarshal(bytes, m); err != nil {
            return err
        }
    } else {
        return status.Error(codes.Internal, ErrProtobufType)
    }

    return nil
}


