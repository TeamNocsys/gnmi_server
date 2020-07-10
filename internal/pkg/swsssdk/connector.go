package swsssdk

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strings"
)

// Errors raised by package swsssdk.
var (
	ErrInvalidParameters = fmt.Errorf("invalid parameters")
	ErrTableNotExist = fmt.Errorf("table does not exist")
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

func (conn *Connector) Connect(db_name string) bool {
	if id := gscfg.GetDBId(db_name); id > 0 {
		return conn.mgmt.connect(id, gscfg.GetDBHostname(db_name), gscfg.GetDBPort(db_name))
	}
	return false
}

func (conn *Connector) Close(db_name string) {
	if id := gscfg.GetDBId(db_name); id > 0 {
		conn.mgmt.close(id)
	}
}

func (conn *Connector) Set(db_name string, key, field string, value interface{}) (bool, error) {
	if id := gscfg.GetDBId(db_name); id > 0 {
		num, err := conn.mgmt.hset(id, key, field, value)
		return num > 0, err
	} else {
		return false, ErrTableNotExist
	}
}

func (conn *Connector) Get(db_name string, key, field string) (string, error) {
	if id := gscfg.GetDBId(db_name); id > 0 {
		return conn.mgmt.get(id, key, field)
	} else {
		return "", ErrTableNotExist
	}
}

func (conn *Connector) GetAll(db_name string, key string) (map[string]string, error) {
	if id := gscfg.GetDBId(db_name); id > 0 {
		return conn.mgmt.get_all(id, key)
	} else {
		return map[string]string{}, ErrTableNotExist
	}
}

func (conn *Connector) Delete(db_name string, key string) (bool, error) {
	if id := gscfg.GetDBId(db_name); id > 0 {
		num, err := conn.mgmt.delete(id, key)
		return num > 0, err
	} else {
		return false, ErrTableNotExist
	}
}

func (conn *Connector) DeleteAllByPattern(db_name string, pattern string) (bool, error) {
	if id := gscfg.GetDBId(db_name); id > 0 {
		num, err := conn.mgmt.delete_all_by_pattern(id, pattern)
		return num > 0, err
	} else {
		return false, ErrTableNotExist
	}
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

func (conn *Connector) typed_to_raw(typed_data map[string]interface{}) map[string]string {
	if typed_data == nil {
		return map[string]string{}
	} else if len(typed_data) == 0 {
		return map[string]string{"NULL": "NULL"}
	} else {
		raw_data := make(map[string]string)
		for key, value := range typed_data {
			switch value.(type) {
			case []interface{}:
				strs := []string{}
				for _, v := range value.([]interface{}) {
					strs = append(strs, fmt.Sprintf("%v", v))
				}
				raw_data[fmt.Sprintf("%s@", key)] = strings.Join(strs,",")
			default:
				raw_data[key] = fmt.Sprintf("%v", value)
			}
		}
		return raw_data
	}
}