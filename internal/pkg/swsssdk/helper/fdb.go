package helper

import (
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/ygot/proto/ywrapper"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
    "strings"
)

type Fdb struct {
    Key string
    Client command.Client
    Data *sonicpb.AcctonFdb_Fdb_FdbList
}

func (c *Fdb) LoadFromDB() error {
    conn := c.Client.State()
    if conn == nil {
        return swsssdk.ErrConnNotExist
    }

    // 获取配置信息
    if c.Data == nil {
        c.Data = &sonicpb.AcctonFdb_Fdb_FdbList{}
    }
    if data, err := conn.GetAll(swsssdk.STATE_DB, []string{"FDB_TABLE", c.Key}); err != nil {
        return err
    } else {
        for k, v := range data {
            switch k {
            case "type":
                switch strings.ToUpper(v) {
                case "STATIC":
                    c.Data.Type = sonicpb.AcctonFdb_Fdb_FdbList_TYPE_STATIC
                case "DYNAMIC":
                    c.Data.Type = sonicpb.AcctonFdb_Fdb_FdbList_TYPE_DYNAMIC
                }

            case "port":
                c.Data.Port = &ywrapper.StringValue{Value: v}
            }
        }
    }

    return nil
}

func (c *Fdb) SaveToDB(replace bool) error {
    e := make(map[string]interface{})
    if c.Data.Type != sonicpb.AcctonFdb_Fdb_FdbList_TYPE_UNSET {
        switch c.Data.Type {
        case sonicpb.AcctonFdb_Fdb_FdbList_TYPE_STATIC:
            e["type"] = "static"
        case sonicpb.AcctonFdb_Fdb_FdbList_TYPE_DYNAMIC:
            e["type"] = "dynamic"
        }
    }
    if c.Data.Port != nil {
        e["port"] = c.Data.Port.Value
    }

    conn := c.Client.Config()
    if conn == nil {
        return swsssdk.ErrConnNotExist
    }

    if replace {
        if _, err := conn.SetEntry(swsssdk.CONFIG_DB, []string{"FDB", c.Key}, e); err != nil {
            return err
        }
    } else {
        if _, err := conn.ModEntry(swsssdk.CONFIG_DB, []string{"FDB", c.Key}, e); err != nil {
            return err
        }
    }
    return nil
}

func (c *Fdb) RemoveFromDB() error {
    conn := c.Client.State()
    if conn == nil {
        return swsssdk.ErrConnNotExist
    }
    if _, err := conn.DeleteAllByPattern(swsssdk.CONFIG_DB, []string{"FDB", c.Key}); err != nil {
        return err
    }
    return nil
}
