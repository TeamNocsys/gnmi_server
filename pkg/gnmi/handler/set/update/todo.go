package update

import (
    "context"
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/golang/protobuf/proto"
    gpb "github.com/openconfig/gnmi/proto/gnmi"
    "github.com/sirupsen/logrus"
    "gnmi_server/cmd/command"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func TodoHandler(ctx context.Context, value *gpb.TypedValue, db command.Client) error {
    info := &sonicpb.SonicTodo{}
    if bytes := value.GetBytesVal(); bytes == nil {
        return status.Error(codes.Internal, ErrProtobufType)
    } else if err := proto.Unmarshal(bytes, info); err != nil {
        return err
    } else {
        if info.Todo != nil {
            if info.Todo.TodoList != nil {
                for _, v := range info.Todo.TodoList {
                    logrus.Debug("SET|" + v.Name + "|" + v.TodoList.Json.Value)
                }
            }
        }
    }

    return nil
}