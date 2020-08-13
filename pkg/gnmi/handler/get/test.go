package get

import (
    "context"
    gpb "github.com/openconfig/gnmi/proto/gnmi"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func Test(ctx context.Context, r *gpb.GetRequest) (*gpb.GetResponse, error) {
    return nil, status.Errorf(codes.Unimplemented, "test")
}
