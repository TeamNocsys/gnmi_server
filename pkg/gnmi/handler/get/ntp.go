package get

import (
    "context"
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/gnmi/proto/gnmi"
    "gnmi_server/cmd/command"
    "gnmi_server/pkg/gnmi/cmd"
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

    sn := &sonicpb.AcctonNtp{
        Ntp: &sonicpb.AcctonNtp_Ntp{},
    }
    if hkeys, err := conn.GetKeys("NTP_SERVER", spec); err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    } else {
        for _, hkey := range hkeys {
            keys := conn.SplitKeys(hkey)
            c := cmd.NewNtpAdapter(keys[0], db)
            if data, err := c.Show(r.Type); err != nil {
                return nil, err
            } else {
                sn.Ntp.NtpList = append(sn.Ntp.NtpList,
                    &sonicpb.AcctonNtp_Ntp_NtpListKey{
                        Ip: keys[0],
                        NtpList: data,
                    })
            }
        }
    }

    response, err := handler_utils.CreateGetResponse(ctx, r, sn)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return response, nil
}