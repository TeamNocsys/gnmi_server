package swsssdk

type DBConfig struct {
    id        int
    separator string
    hostname  string
    port      int
}

type SonicConfig struct {
    ver    string
    dbcfgs map[string]DBConfig
}

func (sc *SonicConfig) GetList() []string {
    keys := make([]string, 0, len(sc.dbcfgs))
    for key := range sc.dbcfgs {
        keys = append(keys, key)
    }
    return keys
}

func (sc *SonicConfig) GetDBId(db_name string) int {
    if dbcfg, ok := sc.dbcfgs[db_name]; ok {
        return dbcfg.id
    }
    return -1
}

func (sc *SonicConfig) GetDBSeparator(db_name string) string {
    if dbcfg, ok := sc.dbcfgs[db_name]; ok {
        return dbcfg.separator
    }
    return "|"
}

func (sc *SonicConfig) GetDBHostname(db_name string) string {
    if dbcfg, ok := sc.dbcfgs[db_name]; ok {
        return dbcfg.hostname
    }
    return REDIS_HOST
}

func (sc *SonicConfig) GetDBPort(db_name string) int {
    if dbcfg, ok := sc.dbcfgs[db_name]; ok {
        return dbcfg.port
    }
    return REDIS_PORT
}
