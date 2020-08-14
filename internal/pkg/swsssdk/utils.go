package swsssdk

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
)

const (
    REDIS_HOST = "127.0.0.1"
    /*
    	SONiC does not use a password-protected database. By default, Redis will only allow connections to unprotected
    	DBs over the loopback ip.
    */

    REDIS_PORT = 6379
    /*
    	SONiC uses the default port.
    */

    MAXIMUM_NUMBER_CONNECT_RETRY = 2
    /*
    	Maximum number of retries before giving up.
    */

    DATA_RETRIEVAL_WAIT_TIME = 3
    /*
    	Wait period in seconds to wait before attempting to retrieve missing data.
    */

    PUB_SUB_NOTIFICATION_TIMEOUT = 10.0
    /*
       Time to wait for any given message to arrive via pub-sub.
    */

    PUB_SUB_MAXIMUM_DATA_WAIT = 60.0
    /*
    	Maximum allowable time to wait on a specific pub-sub notification.
    */

    KEYSPACE_PATTERN = "__key*__:*"
    /*
    	Pub-sub keyspace pattern
    */

    KEYSPACE_EVENTS = "KEA"
    /*
    	In Redis, by default keyspace events notifications are disabled because while not
    	very sensible the feature uses some CPU power. Notifications are enabled using
    	the notify-keyspace-events of redis.conf or via the CONFIG SET.
    	In order to enable the feature a non-empty string is used, composed of multiple characters,
    	where every character has a special meaning according to the following table:
    	K - Keyspace events, published with __keyspace@<db>__ prefix.
    	E - Keyevent events, published with __keyevent@<db>__ prefix.
    	g - Generic commands (non-type specific) like DEL, EXPIRE, RENAME, ...
    	$ - String commands
    	l - List commands
    	s - Set commands
    	h - Hash commands
    	z - Sorted set commands
    	x - Expired events (events generated every time a key expires)
    	e - Evicted events (events generated when a key is evicted for maxmemory)
    	A - Alias for g$lshzxe, so that the "AKE" string means all the events.
    	ACS Redis db mainly uses hash, therefore h is selected.
    */

    CONFIG_DB       = "CONFIG_DB"
    APPL_DB         = "APPL_DB"
    ASIC_DB         = "ASIC_DB"
    COUNTERS_DB     = "COUNTERS_DB"
    LOGLEVEL_DB     = "LOGLEVEL_DB"
    PFC_WD_DB       = "PFC_WD_DB"
    FLEX_COUNTER_DB = "FLEX_COUNTER_DB"
    STATE_DB        = "STATE_DB"
    SNMP_OVERLAY_DB = "SNMP_OVERLAY_DB"
)

// Errors raised by package swsssdk.
var (
    ErrInvalidParameters = fmt.Errorf("invalid parameters")
    ErrDatabaseNotExist  = fmt.Errorf("database does not exist")
    ErrConnNotExist      = fmt.Errorf("database connection does not exist")
)

func Config() *SonicConfig {
    return &gscfg
}

func LoadConfig(path string) error {
    content := []byte(default_config_content)
    if _, err := os.Stat(path); err == nil {
        if content, err = ioutil.ReadFile(path); err != nil {
            return err
        }
    } else {
        return err
    }
    var cfg struct {
        Instances map[string]struct {
            Hostname         string
            Port             int
            Unix_socket_path string
        } `json:"INSTANCES"`
        Databases map[string]struct {
            Id        int
            Separator string
            Instance  string
        } `json:"DATABASES"`
        Version string `json:VERSION`
    }
    if err := json.Unmarshal(content, &cfg); err != nil {
        return err
    } else {
        dbcfgs := make(map[string]DBConfig)
        for k, v := range cfg.Databases {
            if instance, ok := cfg.Instances[v.Instance]; ok {
                dbcfgs[k] = DBConfig{
                    v.Id,
                    v.Separator,
                    instance.Hostname,
                    instance.Port,
                }
            }
        }
        gscfg = SonicConfig{cfg.Version, dbcfgs}
        return nil
    }
}
