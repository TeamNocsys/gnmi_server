package cmd

import (
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/gnmi/proto/gnmi"
    "github.com/openconfig/ygot/proto/ywrapper"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/utils"
    "strconv"
    "strings"
)

type NtpAdapter struct {
    Adapter
    ipaddr string
}

func NewNtpAdapter(ipaddr string, cli command.Client) *NtpAdapter {
    return &NtpAdapter{
        Adapter: Adapter{
            client: cli,
        },
        ipaddr:  ipaddr,
    }
}

func (adpt *NtpAdapter) Show(dataType gnmi.GetRequest_DataType) (*sonicpb.NocsysNtp_Ntp_NtpList, error) {
    if err, data := utils.Utils_execute_cmd("show", "ntp"); err != nil {
        return nil, err
    } else {
        retval := &sonicpb.NocsysNtp_Ntp_NtpList{
            State: &sonicpb.NocsysNtp_Ntp_NtpList_State{},
        }
        infos := strings.Split(data, "\n")
        for i := 2; i < len(infos) - 1; i++ {
            fields := strings.Split(infos[i], " ")
            if fields[0] == adpt.ipaddr {
                if i, err := strconv.ParseUint(fields[5], 10, 64); err != nil {
                    return nil, err
                } else {
                    retval.State.Poll = &ywrapper.UintValue{Value: i}
                }
                if err, now := utils.Utils_execute_cmd("date", "+\"%Y-%m-%dT%H:%M:%SZ%:z\""); err != nil {
                    return nil, err
                } else {
                    retval.State.Current = &ywrapper.StringValue{Value: now}
                }
            }
        }
        return retval, nil
    }
}