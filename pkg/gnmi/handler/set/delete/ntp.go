package delete

import (
    "context"
    "gnmi_server/cmd/command"
    "gnmi_server/pkg/gnmi/cmd"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func NtpHandler(ctx context.Context, kvs map[string]string, db command.Client) error {
    if v, ok := kvs["ip"]; !ok {
        return status.Error(codes.Internal, ErrNoKey)
    } else {
        c := cmd.NewNtpAdapter(v, db)
        return c.Config(nil, cmd.DEL)
    }
}