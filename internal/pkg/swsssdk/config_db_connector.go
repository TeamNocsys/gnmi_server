package swsssdk

import (
    "context"
    "github.com/go-redis/redis/v8"
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
    if id := gscfg.GetDBId(CONFIG_DB_NAME); id > 0 {
        return cc.mgmt.connect(id, gscfg.GetDBHostname(CONFIG_DB_NAME), gscfg.GetDBPort(CONFIG_DB_NAME))
    }
    return false
}

func (cc *ConfigDBConnector) Disconnect() {
    cc.Connector.Disconnect(CONFIG_DB_NAME)
}

func (cc *ConfigDBConnector) SetEntry(table string, key interface{}, values map[string]interface{}) (bool, error) {
    if hkey, err := NewHashKey(table, key); err != nil {
        return false, err
    } else {
        if id := gscfg.GetDBId(CONFIG_DB_NAME); id > 0 {
            if values == nil {
                num, err := cc.mgmt.delete(id, hkey.Get())
                return num > 0, err
            } else {
                if entries, err := cc.GetEntry(table, key); err != nil {
                    return false, err
                } else if num, err := cc.mgmt.hset(id, hkey.Get(), cc.typed_to_raw(values)); err != nil {
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

func (cc *ConfigDBConnector) ModEntry(table string, key interface{}, values map[string]interface{}) (bool, error) {
    if hkey, err := NewHashKey(table, key); err != nil {
        return false, err
    } else {
        if id := gscfg.GetDBId(CONFIG_DB_NAME); id > 0 {
            if values == nil {
                num, err := cc.mgmt.delete(id, hkey.Get())
                return num > 0, err
            } else {
                num, err := cc.mgmt.hset(id, hkey.Get(), cc.typed_to_raw(values))
                return num > 0, err
            }
        }
        return false, ErrDatabaseNotExist
    }
}

func (cc *ConfigDBConnector) GetEntry(table string, key interface{}) (map[string]interface{}, error) {
    if hkey, err := NewHashKey(table, key); err != nil {
        return map[string]interface{}{}, err
    } else {
        if id := gscfg.GetDBId(CONFIG_DB_NAME); id > 0 {
            if values, err := cc.mgmt.get_all(id, hkey.Get()); err != nil {
                return map[string]interface{}{}, err
            } else {
                return cc.raw_to_typed(values), nil
            }
        }
        return map[string]interface{}{}, ErrDatabaseNotExist
    }

}

func (cc *ConfigDBConnector) GetKeys(table string) ([]HashKey, error) {
    keys := []HashKey{}
    if id := gscfg.GetDBId(CONFIG_DB_NAME); id > 0 {
        if pattern, err := NewHashKey(table, "*"); err == nil {
            if hash_keys, err := cc.mgmt.keys(id, pattern.Get()); err == nil {
                for _, hash_key := range hash_keys {
                    keys = append(keys, NewHashKeyByRaw(hash_key))
                }
            }
            return keys, err
        }
    }
    return keys, ErrDatabaseNotExist
}

func (cc *ConfigDBConnector) GetTable(table string) (map[HashKey](map[string]interface{}), error) {
    content := make(map[HashKey](map[string]interface{}))
    if id := gscfg.GetDBId(CONFIG_DB_NAME); id > 0 {
        if pattern, err := NewHashKey(table, "*"); err == nil {
            if hash_keys, err := cc.mgmt.keys(id, pattern.Get()); err != nil {
                return content, err
            } else {
                for _, hash_key := range hash_keys {
                    if entry, err := cc.mgmt.get_all(id, hash_key); err != nil {
                        return content, err
                    } else {
                        content[NewHashKeyByRaw(hash_key)] = cc.raw_to_typed(entry)
                    }
                }
                return content, nil
            }
        }
    }
    return content, ErrDatabaseNotExist
}

func (cc *ConfigDBConnector) DeleteTable(table string) (bool, error) {
    if id := gscfg.GetDBId(CONFIG_DB_NAME); id > 0 {
        if pattern, err := NewHashKey(table, "*"); err == nil {
            num, err := cc.mgmt.delete_all_by_pattern(id, pattern.Get())
            return num > 0, err
        }
    }
    return false, ErrDatabaseNotExist
}
