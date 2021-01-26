package delete

import (
    "context"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk/helper"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func NtpHandler(ctx context.Context, kvs map[string]string, db command.Client) error {
    if v, ok := kvs["ip"]; !ok {
        return status.Error(codes.Internal, ErrNoKey)
    } else {
        c := &helper.Ntp{
            Key:    v,
            Client: db,
            Data:   nil,
        }
        if err := c.RemoveFromDB(); err != nil {
            return err
        }
    }

    return nil
}