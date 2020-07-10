package swsssdk

import (
	"context"
	"github.com/go-redis/redis/v8"
	"gnmi_server/internal/pkg/swsssdk/utils"
)

const (
	CONFIG_DB_NAME = "CONFIG_DB"
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

func (cc *ConfigDBConnector) Connect() bool {
	if id := gscfg.GetDBId(CONFIG_DB_NAME); id > 0 {
		return cc.mgmt.connect(id, gscfg.GetDBHostname(CONFIG_DB_NAME), gscfg.GetDBPort(CONFIG_DB_NAME))
	}
	return false
}

func (cc *ConfigDBConnector) SetEntry(table string, key interface{}, values map[string]interface{}) (bool, error) {
	if hkey, err := utils.NewHashKey(table, key); err != nil {
		return false, err
	} else {
		if id := gscfg.GetDBId(CONFIG_DB_NAME); id > 0 {
			if values == nil {
				num, err := cc.mgmt.delete(id, hkey.Get())
				return num > 0, err
			} else {
				if entries, err := cc.GetEntry(table, key); err != nil {
					return false, err
				} else if success, err := cc.mgmt.hmset(id, hkey.Get(), cc.typed_to_raw(values)); err != nil || !success {
					return success, err
				} else {
					// 删除旧的无效条目
					for k, _ := range entries {
						if _, ok := values[k]; !ok {
							if _, err := cc.mgmt.delete(id, k); err != nil {
								return false, err
							}
						}
					}
				}
			}
		}
		return false, ErrTableNotExist
	}
}

func (cc *ConfigDBConnector) ModEntry(table string, key interface{}, values map[string]interface{}) (bool, error) {
	if hkey, err := utils.NewHashKey(table, key); err != nil {
		return false, err
	} else {
		if id := gscfg.GetDBId(CONFIG_DB_NAME); id > 0 {
			if values == nil {
				num, err := cc.mgmt.delete(id, hkey.Get())
				return num > 0, err
			} else {
				return cc.mgmt.hmset(id, hkey.Get(), cc.typed_to_raw(values))
			}
		}
		return false, ErrTableNotExist
	}
}

func (cc *ConfigDBConnector) GetEntry(table string, key interface{}) (map[string]interface{}, error) {
	if hkey, err := utils.NewHashKey(table, key); err != nil {
		return map[string]interface{}{}, err
	} else {
		if id := gscfg.GetDBId(CONFIG_DB_NAME); id > 0 {
			if values, err := cc.mgmt.get_all(id, hkey.Get()); err != nil {
				return map[string]interface{}{}, err
			} else {
				return cc.raw_to_typed(values), nil
			}
		}
		return map[string]interface{}{}, ErrTableNotExist
	}

}

func (cc *ConfigDBConnector) GetKeys(table string) ([]utils.HashKey, error) {
	keys := []utils.HashKey{}
	if id := gscfg.GetDBId(CONFIG_DB_NAME); id > 0 {
		if pattern, err := utils.NewHashKey(table, "*"); err == nil {
			if hash_keys, err := cc.mgmt.keys(id, pattern.Get()); err != nil {
				for _, hash_key := range hash_keys {
					keys = append(keys, utils.NewHashKeyByRaw(hash_key))
				}
			} else {
				return keys, err
			}
		}
	}
	return keys, ErrTableNotExist
}

func (cc *ConfigDBConnector) GetTable(table string) (map[utils.HashKey](map[string]interface{}), error) {
	content := make(map[utils.HashKey](map[string]interface{}))
	if id := gscfg.GetDBId(CONFIG_DB_NAME); id > 0 {
		if pattern, err := utils.NewHashKey(table, "*"); err == nil {
			if hash_keys, err := cc.mgmt.keys(id, pattern.Get()); err != nil {
				return content, err
			} else {
				for _, hash_key := range hash_keys {
					if entry, err := cc.mgmt.get_all(id, hash_key); err != nil {
						return content, err
					} else {
						content[utils.NewHashKeyByRaw(hash_key)] = cc.raw_to_typed(entry)
					}
				}
				return content, nil
			}
		}
	}
	return content, ErrTableNotExist
}

func (cc *ConfigDBConnector) DeleteTable(table string) (bool, error) {
	if id := gscfg.GetDBId(CONFIG_DB_NAME); id > 0 {
		if pattern, err := utils.NewHashKey(table, "*"); err == nil {
			num, err := cc.mgmt.delete_all_by_pattern(id, pattern.Get())
			return num > 0, err
		}
	}
	return false, ErrTableNotExist
}
