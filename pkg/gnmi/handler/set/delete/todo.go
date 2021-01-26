package delete

import (
    "context"
    "github.com/sirupsen/logrus"
    "gnmi_server/cmd/command"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func TodoHandler(ctx context.Context, kvs map[string]string, db command.Client) error {
    if v, ok := kvs["name"]; !ok {
        return status.Error(codes.Internal, ErrNoKey)
    } else {
        logrus.Debugf("DEL|" + v)
    }

    return nil
}