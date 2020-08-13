package swsssdk

import (
    "fmt"
    "strings"
)

const (
    TABLE_NAME_SEPARATOR = "|"
    KEY_SEPARATOR        = "|"
)

type HashKey struct {
    value string
}

func NewHashKey(table string, keys interface{}) (HashKey, error) {
    var hkey HashKey
    switch v := keys.(type) {
    case string:
        hkey = HashKey{fmt.Sprintf("%s%s%s", strings.ToUpper(table), TABLE_NAME_SEPARATOR, v)}
    case []string:
        hkey = HashKey{fmt.Sprintf("%s%s%s", strings.ToUpper(table), TABLE_NAME_SEPARATOR, strings.Join(v, KEY_SEPARATOR))}
    default:
        return HashKey{}, ErrInvalidParameters
    }
    return hkey, nil
}

func NewHashKeyByRaw(key string) HashKey {
    return HashKey{key}
}

func (hk *HashKey) Get() string {
    return hk.value
}

func (hk *HashKey) GetTable() string {
    return strings.SplitN(hk.value, TABLE_NAME_SEPARATOR, 1)[0]
}

func (hk *HashKey) GetKey() string {
    return strings.SplitN(hk.value, TABLE_NAME_SEPARATOR, 1)[1]
}

func (hk *HashKey) GetKeys() []string {
    return strings.Split(hk.GetKey(), KEY_SEPARATOR)
}
