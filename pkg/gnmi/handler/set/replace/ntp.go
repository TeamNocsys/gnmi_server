package replace

import (
    "context"
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/golang/protobuf/proto"
    gpb "github.com/openconfig/gnmi/proto/gnmi"
    "gnmi_server/cmd/command"
    "gnmi_server/pkg/gnmi/cmd"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func NtpHandler(ctx context.Context, value *gpb.TypedValue, db command.Client) error {
    info := &sonicpb.NocsysNtp{}
    if bytes := value.GetBytesVal(); bytes == nil {
        return status.Error(codes.Internal, ErrProtobufType)
    } else if err := proto.Unmarshal(bytes, info); err != nil {
        return err
    } else {
        if info.Ntp != nil {
            if info.Ntp.NtpList != nil {
                for _, v := range info.Ntp.NtpList {
                    c := cmd.NewNtpAdapter(v.Ip, db)
                    if err := c.Config(v.NtpList, cmd.ADD); err != nil {
                        return nil
                    }
                }
            }
        }
    }

    return nil
}