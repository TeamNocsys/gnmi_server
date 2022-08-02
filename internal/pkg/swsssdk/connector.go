package swsssdk

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/go-redis/redis/v8"
    "github.com/sirupsen/logrus"
    "strings"
)

type Connector struct {
    mgmt *manager
}

func NewConnector() *Connector {
    return &Connector{&manager{
        make(map[int]*redis.Client),
        make(map[int]*redis.PubSub),
        context.Background(),
    }}
}

func (conn *Connector) Close() {
    conn.mgmt.close()
}

func (conn *Connector) Connect(db_name string) bool {
    if id := gscfg.GetDBId(db_name); id >= 0 {
        logrus.Info(fmt.Sprintf("Connect|%s|id:%d|host:%s|port:%d",
            db_name,
            id,
            gscfg.GetDBHostname(db_name),
            gscfg.GetDBPort(db_name)))
        return conn.mgmt.connect(id, gscfg.GetDBHostname(db_name), gscfg.GetDBPort(db_name))
    }
    return false
}

func (conn *Connector) Disconnect(db_name string) {
    if id := gscfg.GetDBId(db_name); id >= 0 {
        conn.mgmt.disconnect(id)
    }
}

func (conn *Connector) Set(db_name string, keys interface{}, value interface{}) (bool, error) {
    if id := gscfg.GetDBId(db_name); id >= 0 {
        if key := conn.serialize_key(db_name, keys); key != "" {
            num, err := conn.mgmt.hset(id, key, value)
            return num > 0, err
        } else {
            return false, ErrInvalidParameters
        }
    } else {
        return false, ErrDatabaseNotExist
    }
}

func (conn *Connector) Get(db_name string, keys interface{}, field string) (string, error) {
    if id := gscfg.GetDBId(db_name); id >= 0 {
        if key := conn.serialize_key(db_name, keys); key != "" {
            return conn.mgmt.get(id, key, field)
        } else {
            return "", ErrInvalidParameters
        }
    } else {
        return "", ErrDatabaseNotExist
    }
}

func (conn *Connector) GetAll(db_name string, keys interface{}) (map[string]string, error) {
    if id := gscfg.GetDBId(db_name); id >= 0 {
        if key := conn.serialize_key(db_name, keys); key != "" {
            return conn.mgmt.get_all(id, key)
        } else {
            return map[string]string{}, ErrInvalidParameters
        }
    } else {
        return map[string]string{}, ErrDatabaseNotExist
    }
}

func (conn *Connector) GetAllWithTrace(db_name string, keys interface{}) (map[string]string, error) {
    if id := gscfg.GetDBId(db_name); id >= 0 {
        if key := conn.serialize_key(db_name, keys); key != "" {
            logrus.Tracef("KEY|%s", key)
            v, err := conn.mgmt.get_all(id, key)
            s, _ := json.Marshal(v)
            logrus.Tracef("KEY|VALUE|%s", s)
            return v, err
        } else {
            return map[string]string{}, ErrInvalidParameters
        }
    } else {
        return map[string]string{}, ErrDatabaseNotExist
    }
}

func (conn *Connector) GetAllByPattern(db_name string, patterns interface{}) (map[string]map[string]string, error) {
    content := make(map[string]map[string]string)
    if id := gscfg.GetDBId(db_name); id >= 0 {
        if pattern := conn.serialize_key(db_name, patterns); pattern != "" {
            if keys, err := conn.mgmt.keys(id, pattern); err != nil {
                return content, err
            } else {
                for _, key := range keys {
                    if entry, err := conn.mgmt.get_all(id, key); err != nil {
                        return content, err
                    } else {
                        content[key] = entry
                    }
                }
                return content, nil
            }
        } else {
            return content, ErrInvalidParameters
        }
    } else {
        return content, ErrDatabaseNotExist
    }
}

func (conn *Connector) Delete(db_name string, keys interface{}) (bool, error) {
    if id := gscfg.GetDBId(db_name); id >= 0 {
        if key := conn.serialize_key(db_name, keys); key != "" {
            num, err := conn.mgmt.delete(id, key)
            return num > 0, err
        } else {
            return false, ErrInvalidParameters
        }
    } else {
        return false, ErrDatabaseNotExist
    }
}

func (conn *Connector) DeleteAllByPattern(db_name string, patterns interface{}) (bool, error) {
    if id := gscfg.GetDBId(db_name); id >= 0 {
        if pattern := conn.serialize_key(db_name, patterns); pattern != "" {
            num, err := conn.mgmt.delete_all_by_pattern(id, pattern)
            return num > 0, err
        } else {
            return false, ErrInvalidParameters
        }
    } else {
        return false, ErrDatabaseNotExist
    }
}

func (conn *Connector) SetEntry(db_name string, keys interface{}, values map[string]interface{}) (bool, error) {
    if key := conn.serialize_key(db_name, keys); key == "" {
        return false, ErrInvalidParameters
    } else {
        if id := gscfg.GetDBId(db_name); id >= 0 {
            if values == nil {
                num, err := conn.mgmt.delete(id, key)
                return num > 0, err
            } else {
                if entries, err := conn.GetEntry(db_name, keys); err != nil {
                    return false, err
                } else if num, err := conn.mgmt.hset(id, key, conn.typed_to_raw(values)); err != nil {
                    return num > 0, err
                } else {
                    // 删除旧的无效条目
                    for k, _ := range entries {
                        if _, ok := values[k]; !ok {
                            if _, err := conn.mgmt.delete(id, k); err != nil {
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

func (conn *Connector) ModEntry(db_name string, keys interface{}, values map[string]interface{}) (bool, error) {
    if key := conn.serialize_key(db_name, keys); key == "" {
        return false, ErrInvalidParameters
    } else {
        if id := gscfg.GetDBId(db_name); id >= 0 {
            if values == nil {
                num, err := conn.mgmt.delete(id, key)
                return num > 0, err
            } else {
                num, err := conn.mgmt.hset(id, key, conn.typed_to_raw(values))
                return num > 0, err
            }
        }
        return false, ErrDatabaseNotExist
    }
}

func (conn *Connector) GetEntry(db_name string, keys interface{}) (map[string]interface{}, error) {
    if key := conn.serialize_key(db_name, keys); key == "" {
        return map[string]interface{}{}, ErrInvalidParameters
    } else {
        if id := gscfg.GetDBId(db_name); id >= 0 {
            if values, err := conn.mgmt.get_all(id, key); err != nil {
                return map[string]interface{}{}, err
            } else {
                return conn.raw_to_typed(values), nil
            }
        }
        return map[string]interface{}{}, ErrDatabaseNotExist
    }

}

func (conn *Connector) HasEntry(db_name string, keys interface{}) (bool, error) {
    if id := gscfg.GetDBId(db_name); id > 0 {
        if pattern := conn.serialize_key(db_name, keys); pattern != "" {
            num, err := conn.mgmt.exists(id, pattern)
            return num > 0, err
        }
    }
    return false, ErrDatabaseNotExist
}

func (conn *Connector) GetKeys(db_name string, keys interface{}) ([]string, error) {
    if id := gscfg.GetDBId(db_name); id >= 0 {
        if pattern := conn.serialize_key(db_name, keys); pattern != "" {
            if ss, err := conn.mgmt.keys(id, pattern); err != nil {
                return []string{}, err
            } else {
                var hkeys []string
                for _, s := range ss {
                    hkeys = append(hkeys, s)
                }
                return hkeys, nil
            }
        }
    }
    return []string{}, ErrDatabaseNotExist
}

func (conn *Connector) SplitKeys(db_name string, keys string) []string {
    s := gscfg.GetDBSeparator(db_name)
    return strings.Split(keys, s)[1:]
}

func (conn *Connector) SplitKeysN(db_name string, keys string, m int) []string {
    s := gscfg.GetDBSeparator(db_name)
    return strings.SplitN(keys, s, m)[1:]
}

func (conn *Connector) serialize_key(db_name string, keys interface{}) string {
    key := ""
    switch keys.(type) {
    case string:
        key = keys.(string)
    case []string:
        key = strings.Join(keys.([]string), gscfg.GetDBSeparator(db_name))
    }

    return key
}

func (conn *Connector) raw_to_typed(raw_data map[string]string) map[string]interface{} {
    if raw_data == nil {
        return map[string]interface{}{}
    } else {
        typed_data := make(map[string]interface{})
        for key, value := range raw_data {
            if key == "NULL" {
                continue
            } else if strings.HasSuffix(key, "@") {
                typed_data[strings.TrimRight(key, "@")] = strings.Split(value, ",")
            } else {
                typed_data[key] = value
            }
        }
        return typed_data
    }
}

func (conn *Connector) typed_to_raw(typed_data map[string]interface{}) map[string]interface{} {
    if typed_data == nil {
        return map[string]interface{}{}
    } else if len(typed_data) == 0 {
        return map[string]interface{}{"NULL": "NULL"}
    } else {
        raw_data := make(map[string]interface{})
        for key, value := range typed_data {
            switch value.(type) {
            case []interface{}:
                strs := []string{}
                for _, v := range value.([]interface{}) {
                    strs = append(strs, fmt.Sprintf("%v", v))
                }
                raw_data[fmt.Sprintf("%s@", key)] = strings.Join(strs, ",")
            default:
                raw_data[key] = value
            }
        }
        return raw_data
    }
}
