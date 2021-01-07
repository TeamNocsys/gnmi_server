package swsssdk

import (
    "context"
    "github.com/go-redis/redis/v8"
    "strings"
)

type ConfigDBConnector struct {
    Connector
}

func NewConfigDBConnector() *ConfigDBConnector {
    return &ConfigDBConnector{Connector{
        &manager{
            make(map[int]*redis.Client),
            make(map[int]*redis.PubSub),
            context.Background(),
        }}}
}

func (cc *ConfigDBConnector) Close() {
    cc.mgmt.close()
}

func (cc *ConfigDBConnector) Connect() bool {
    if id := gscfg.GetDBId(CONFIG_DB); id > 0 {
        return cc.mgmt.connect(id, gscfg.GetDBHostname(CONFIG_DB), gscfg.GetDBPort(CONFIG_DB))
    }
    return false
}

func (cc *ConfigDBConnector) Disconnect() {
    cc.Connector.Disconnect(CONFIG_DB)
}

func (cc *ConfigDBConnector) SetEntry(table string, keys interface{}, values map[string]interface{}) (bool, error) {
    if key := cc.serialize_key(table, keys); key == "" {
        return false, ErrInvalidParameters
    } else {
        if id := gscfg.GetDBId(CONFIG_DB); id > 0 {
            if values == nil {
                num, err := cc.mgmt.delete(id, key)
                return num > 0, err
            } else {
                if entries, err := cc.GetEntry(table, keys); err != nil {
                    return false, err
                } else if num, err := cc.mgmt.hset(id, key, cc.typed_to_raw(values)); err != nil {
                    return num > 0, err
                } else {
                    // 删除旧的无效条目
                    for k, _ := range entries {
                        if _, ok := values[k]; !ok {
                            if _, err := cc.mgmt.delete(id, k); err != nil {
                                return false, err
                            }
                        }
                    }
                    return true, nil
                }
            }
        }
        return false, ErrDatabaseNotExist
    }
}

func (cc *ConfigDBConnector) ModEntry(table string, keys interface{}, values map[string]interface{}) (bool, error) {
    if key := cc.serialize_key(table, keys); key == "" {
        return false, ErrInvalidParameters
    } else {
        if id := gscfg.GetDBId(CONFIG_DB); id > 0 {
            if values == nil {
                num, err := cc.mgmt.delete(id, key)
                return num > 0, err
            } else {
                num, err := cc.mgmt.hset(id, key, cc.typed_to_raw(values))
                return num > 0, err
            }
        }
        return false, ErrDatabaseNotExist
    }
}

func (cc *ConfigDBConnector) GetEntry(table string, keys interface{}) (map[string]interface{}, error) {
    if key := cc.serialize_key(table, keys); key == "" {
        return map[string]interface{}{}, ErrInvalidParameters
    } else {
        if id := gscfg.GetDBId(CONFIG_DB); id > 0 {
            if values, err := cc.mgmt.get_all(id, key); err != nil {
                return map[string]interface{}{}, err
            } else {
                return cc.raw_to_typed(values), nil
            }
        }
        return map[string]interface{}{}, ErrDatabaseNotExist
    }

}

func (cc *ConfigDBConnector) GetKeys(table string) ([]string, error) {
    if id := gscfg.GetDBId(CONFIG_DB); id > 0 {
        if pattern := cc.serialize_key(table, "*"); pattern != "" {
            return cc.mgmt.keys(id, pattern)
        }
    }
    return []string{}, ErrDatabaseNotExist
}

func (cc *ConfigDBConnector) GetTable(table string) (map[string]map[string]interface{}, error) {
    content := make(map[string]map[string]interface{})
    if id := gscfg.GetDBId(CONFIG_DB); id > 0 {
        if pattern := cc.serialize_key(table, "*"); pattern != "" {
            if keys, err := cc.mgmt.keys(id, pattern); err != nil {
                return content, err
            } else {
                s := gscfg.GetDBSeparator(CONFIG_DB)
                for _, key := range keys {
                    if entry, err := cc.mgmt.get_all(id, key); err != nil {
                        return content, err
                    } else {
                        mkeys := strings.SplitN(key, s, 2)
                        content[mkeys[len(mkeys)-1]] = cc.raw_to_typed(entry)
                    }
                }
                return content, nil
            }
        }
    }
    return content, ErrDatabaseNotExist
}

func (cc *ConfigDBConnector) DeleteTable(table string) (bool, error) {
    if id := gscfg.GetDBId(CONFIG_DB); id > 0 {
        if pattern := cc.serialize_key(table, "*"); pattern != "" {
            num, err := cc.mgmt.delete_all_by_pattern(id, pattern)
            return num > 0, err
        }
    }
    return false, ErrDatabaseNotExist
}

func (cc *ConfigDBConnector) HasEntry(table string, keys interface{}) (bool, error) {
    if id := gscfg.GetDBId(CONFIG_DB); id > 0 {
        if pattern := cc.serialize_key(table, keys); pattern != "" {
            num, err := cc.mgmt.exists(id, pattern)
            return num > 0, err
        }
    }
    return false, ErrDatabaseNotExist
}

func (conn *ConfigDBConnector) serialize_key(table string, keys interface{}) string {
    switch keys.(type) {
    case string:
        return strings.Join([]string{table, keys.(string)}, gscfg.GetDBSeparator(CONFIG_DB))
    case []string:
        var merge []string
        merge = append(merge, table)
        merge = append(merge,  keys.([]string)...)
        return strings.Join(merge, gscfg.GetDBSeparator(CONFIG_DB))
    default:
        return ""
    }
}
