package helper

import (
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/ygot/proto/ywrapper"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
)

type Neighbor struct {
    Keys []string
    Client command.Client
    Data *sonicpb.AcctonNeighor_Neighor_NeighorList
}

func (c *Neighbor) LoadFromDB() error {
    conn := c.Client.State()
    if conn == nil {
        return swsssdk.ErrConnNotExist
    }

    if c.Data == nil {
        c.Data = &sonicpb.AcctonNeighor_Neighor_NeighorList{}
    }
    if data, err := conn.GetAll(swsssdk.APPL_DB, append([]string{"NEIGH_TABLE"}, c.Keys...)); err != nil {
        return err
    } else {
        for k, v := range data {
            switch k {
            case "neigh":
                c.Data.Neigh = &ywrapper.StringValue{Value: v}
            }
        }
    }
    return nil
}

func (c *Neighbor) SaveToDB(replace bool) error {
    e := make(map[string]interface{})
    if c.Data.Neigh != nil {
        e["neigh"] = c.Data.Neigh.Value
    }

    conn := c.Client.Config()
    if conn == nil {
        return swsssdk.ErrConnNotExist
    }

    if replace {
        if _, err := conn.SetEntry("NEIGH", c.Keys, e); err != nil {
            return err
        }
    } else {
        if _, err := conn.ModEntry("NEIGH", c.Keys, e); err != nil {
            return err
        }
    }
    return nil
}

func (c *Neighbor) RemoveFromDB() error {
    conn := c.Client.Config()
    if conn == nil {
        return swsssdk.ErrConnNotExist
    }
    if _, err := conn.DeleteAllByPattern(swsssdk.CONFIG_DB, append([]string{"NEIGH"}, c.Keys...)); err != nil {
        return err
    }
    return nil
}