// +build ec

package cmd

import (
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
)

func (adpt *NtpAdapter) Config(data *sonicpb.NocsysNtp_Ntp_NtpList, oper OperType) error {
    var cmdstr string
    if oper == ADD || oper == UPDATE {
        cmdstr = "config ntp add " + adpt.ipaddr
    } else if oper == DEL {
        cmdstr = "config ntp del " + adpt.ipaddr
    }
    return adpt.exec(cmdstr)
}