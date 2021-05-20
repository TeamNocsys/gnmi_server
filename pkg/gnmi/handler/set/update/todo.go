package update

import (
    "context"
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/golang/protobuf/jsonpb"
    "github.com/golang/protobuf/proto"
    gpb "github.com/openconfig/gnmi/proto/gnmi"
    "github.com/sirupsen/logrus"
    "gnmi_server/cmd/command"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func TodoHandler(ctx context.Context, value *gpb.TypedValue, db command.Client) error {
    info := &sonicpb.NocsysTodo{}
    if bytes := value.GetBytesVal(); bytes == nil {
        return status.Error(codes.Internal, ErrProtobufType)
    } else if err := proto.Unmarshal(bytes, info); err != nil {
        return err
    } else {
        m := jsonpb.Marshaler{}
        s, _ := m.MarshalToString(info)
        logrus.Tracef("UPDATE|%s", s)
        if info.Todo != nil {
            if info.Todo.TodoList != nil {
                for _, v := range info.Todo.TodoList {
                    logrus.Debug("UPDATE|" + v.Name + "|" + v.TodoList.Json.Value)
                }
            }
        }
    }

    return nil
}