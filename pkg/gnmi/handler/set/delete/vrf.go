package delete

import (
    "context"
    "gnmi_server/cmd/command"
    "gnmi_server/pkg/gnmi/cmd"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func VrfHandler(ctx context.Context, kvs map[string]string, db command.Client) error {
    if v, ok := kvs["vrf-name"]; !ok {
        return status.Error(codes.Internal, ErrNoKey)
    } else {
        c := cmd.NewVrfAdapter(v, db)
        return c.Config(nil, cmd.DEL)
    }
}