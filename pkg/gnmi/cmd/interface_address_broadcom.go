// +build broadcom

package cmd

import (
    "gnmi_server/internal/pkg/swsssdk"
)

func (adpt *IfAddrAdapter) Config(data interface{}, oper OperType) error {
    conn := adpt.client.Config()
    if conn == nil {
        return swsssdk.ErrConnNotExist
    }

    var cmdstr string
    if oper == ADD {
        // 如果存在则跳过重复设置
        if ok, err := conn.HasEntry(IfType_table[int32(adpt.ifType)], []string{adpt.ifname, adpt.ipaddr}); err != nil {
            return err
        } else if ok {
            return nil
        }

        cmdstr = "config interface ip add " + adpt.ifname + " " + adpt.ipaddr
    } else if oper == DEL {
        // 如果不存在则跳过删除
        if ok, err := conn.HasEntry(IfType_table[int32(adpt.ifType)], []string{adpt.ifname, adpt.ipaddr}); err != nil {
            return err
        } else if !ok {
            return nil
        }

        cmdstr = "config interface ip remove " + adpt.ifname + " " + adpt.ipaddr
    } else {
        return ErrInvalidOperType
    }
    return adpt.exec(cmdstr)
}