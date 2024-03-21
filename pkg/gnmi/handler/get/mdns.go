package get

import (
    "context"

    "gnmi_server/cmd/command"
    handler_utils "gnmi_server/pkg/gnmi/handler/utils"
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/gnmi/proto/gnmi"
    "github.com/openconfig/ygot/proto/ywrapper"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"

    "gnmi_server/pkg/mdns"
    "time"
)

func MdnsInfoHandler(
    ctx context.Context, r *gnmi.GetRequest, db command.Client,
) (*gnmi.GetResponse, error) {

    expireFlag := false

    nocsysMdns := &sonicpb.AcctonMdns{
        Mdns : &sonicpb.AcctonMdns_Mdns {},
    }

    curTime := time.Now()
    mdns.MdnsResolver.Entries.Range(func(key, value interface{}) bool {

        if !curTime.After(value.(mdns.HostEntry).Expiration) {
            tmp_e := &sonicpb.AcctonMdns_Mdns_MdnsListKey {
                IpPrefix : key.(string),
                MdnsList : &sonicpb.AcctonMdns_Mdns_MdnsList {
                    Hostname : &ywrapper.StringValue {
                        Value : value.(mdns.HostEntry).HostName,
                    },
                },
            }
            nocsysMdns.Mdns.MdnsList = append (nocsysMdns.Mdns.MdnsList, tmp_e)
        } else {
            expireFlag = true
        }
        return true
    })

    if expireFlag {
        mdns.MdnsResolver.NotifyExpireEntries()
    }

    response, err := handler_utils.CreateGetResponse(ctx, r, nocsysMdns)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }

    return response, nil
}
