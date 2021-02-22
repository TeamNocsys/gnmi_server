package get

import (
    "context"
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/gnmi/proto/gnmi"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk/helper"
    "gnmi_server/pkg/gnmi/handler"
    handler_utils "gnmi_server/pkg/gnmi/handler/utils"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func NtpHandler(ctx context.Context, r *gnmi.GetRequest, db command.Client) (*gnmi.GetResponse, error) {
    conn := db.Config()
    if conn == nil {
        return nil, status.Error(codes.Internal, "")
    }
    kvs := handler.FetchPathKey(r)
    spec := "*"
    if v, ok := kvs["ip"]; ok {
        spec = v
    }

    sn := &sonicpb.NocsysNtp{
        Ntp: &sonicpb.NocsysNtp_Ntp{},
    }
    if hkeys, err := conn.GetKeys("NTP_SERVER", spec); err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    } else {
        for _, hkey := range hkeys {
            keys := conn.SplitKeys(hkey)
            c := helper.Ntp{
                Key: keys[0],
                Client: db,
                Data: nil,
            }
            if err := c.LoadFromDB(); err != nil {
                return nil, status.Errorf(codes.Internal, err.Error())
            }
            sn.Ntp.NtpList = append(sn.Ntp.NtpList,
                &sonicpb.NocsysNtp_Ntp_NtpListKey{
                    Ip: keys[0],
                    NtpList: c.Data,
                })
        }
    }

    response, err := handler_utils.CreateGetResponse(ctx, r, sn)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return response, nil
}