package swsssdk

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

const (
	default_config_filename = "/var/run/redis/sonic-db/database_config.json"
	default_config_content = `
	{
		"INSTANCES": {
			"redis":{
				"hostname" : "127.0.0.1",
				"port" : 6379,
				"unix_socket_path" : "/var/run/redis/redis.sock"
			}
		},
		"DATABASES" : {
			"APPL_DB" : {
				"id" : 0,
				"separator": ":",
				"instance" : "redis"
			},
			"ASIC_DB" : {
				"id" : 1,
				"separator": ":",
				"instance" : "redis"
			},
			"COUNTERS_DB" : {
				"id" : 2,
				"separator": ":",
				"instance" : "redis"
			},
			"LOGLEVEL_DB" : {
				"id" : 3,
				"separator": ":",
				"instance" : "redis"
			},
			"CONFIG_DB" : {
				"id" : 4,
				"separator": "|",
				"instance" : "redis"
			},
			"PFC_WD_DB" : {
				"id" : 5,
				"separator": ":",
				"instance" : "redis"
			},
			"FLEX_COUNTER_DB" : {
				"id" : 5,
				"separator": ":",
				"instance" : "redis"
			},
			"STATE_DB" : {
				"id" : 6,
				"separator": "|",
				"instance" : "redis"
			},
			"SNMP_OVERLAY_DB" : {
				"id" : 7,
				"separator": "|",
				"instance" : "redis"
			}
		},
		"VERSION" : "1.0"
	}`
)

var gscfg SonicConfig

func init() {
	content := []byte(default_config_content)
	if _, err := os.Stat(default_config_filename); err == nil {
		if content, err = ioutil.ReadFile(default_config_filename); err != nil {
			log.Fatalln(err)
		}
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
		log.Fatalln(err)
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
	}
}
