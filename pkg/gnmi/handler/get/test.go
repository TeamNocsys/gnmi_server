package get

import (
    "context"
    gpb "github.com/openconfig/gnmi/proto/gnmi"
    "gnmi_server/cmd/command"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func Test(ctx context.Context, r *gpb.GetRequest, db command.Client) (*gpb.GetResponse, error) {
    return nil, status.Errorf(codes.Unimplemented, "test")
}
