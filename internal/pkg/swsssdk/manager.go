package swsssdk

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

const (
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
)

var (
	ErrConnNotExist = fmt.Errorf("database connection does not exist")
)

// 多客户端连接管理
type manager struct {
	clients map[int]*redis.Client
	subs map[int]*redis.PubSub
	ctx context.Context
}

func (mgr *manager) connect(id int, hostname string, port int) bool {
	if _, ok := mgr.clients[id]; ok {
		return true
	}
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", hostname, port),
		Password: "",
		DB: id,
		MaxRetries: MAXIMUM_NUMBER_CONNECT_RETRY,
	})
	if _, err := client.Ping(mgr.ctx).Result(); err != nil {
		return false
	} else {
		mgr.clients[id] = client
		return true
	}
	return false
}

func (db *manager) close(id int) {
	if client, ok := db.clients[id]; ok {
		client.Close()
		delete(db.clients, id)
	}
}

func (mgr *manager) subscribe(id int) {
	if client, ok := mgr.clients[id]; ok {
		if _, ok := mgr.subs[id]; !ok {
			mgr.subs[id] = client.Subscribe(mgr.ctx)
			mgr.subs[id].PSubscribe(mgr.ctx, KEYSPACE_PATTERN)
		}
	}
}

func (mgr *manager) unsubscribe(id int) {
	if sub, ok := mgr.subs[id]; ok {
		sub.PUnsubscribe(mgr.ctx, KEYSPACE_PATTERN)
		if err := sub.Close(); err == nil {
			delete(mgr.subs, id)
		}
	}
}

func (mgr *manager) expire(id int, key string, timeout_sec int) (bool, error) {
	if client, ok := mgr.clients[id]; ok {
		return client.Expire(mgr.ctx, key, time.Duration(timeout_sec) * time.Second).Result()
	}
	return true, ErrConnNotExist
}

func (mgr *manager) exists(id int, key string) (int64, error) {
	if client, ok := mgr.clients[id]; ok {
		return client.Exists(mgr.ctx, key).Result()
	}
	return 0, ErrConnNotExist
}

func (mgr *manager) keys(id int, pattern string) ([]string, error) {
	if client, ok := mgr.clients[id]; ok {
		return client.Keys(mgr.ctx, pattern).Result()
	}
	return []string{}, ErrConnNotExist
}

func (mgr *manager) get(id int, key, field string) (string, error) {
	if client, ok := mgr.clients[id]; ok {
		return client.HGet(mgr.ctx, key, field).Result()
	}
	return "", ErrConnNotExist
}

func (mgr *manager) get_all(id int, key string) (map[string]string, error) {
	if client, ok := mgr.clients[id]; ok {
		return client.HGetAll(mgr.ctx, key).Result()
	}
	return map[string]string{}, ErrConnNotExist
}

func (mgr *manager) hset(id int, key, field string, value interface{}) (int64, error) {
	if client, ok := mgr.clients[id]; ok {
		return client.HSet(mgr.ctx, key, field, value).Result()
	}
	return 0, ErrConnNotExist
}

func (mgr *manager) hmset(id int, key string, value interface{}) (bool, error) {
	if client, ok := mgr.clients[id]; ok {
		return client.HMSet(mgr.ctx, key, value).Result()
	}
	return false, ErrConnNotExist
}

func (mgr *manager) delete(id int, key string) (int64, error) {
	if client, ok := mgr.clients[id]; ok {
		return client.Del(mgr.ctx, key).Result()
	}
	return 0, ErrConnNotExist
}

func (mgr *manager) delete_all_by_pattern(id int, pattern string) (int64, error) {
	if client, ok := mgr.clients[id]; ok {
		if keys, err := client.Keys(mgr.ctx, pattern).Result(); err != nil {
			return 0, err
		} else {
			for _, key := range keys {
				if _, err := client.Del(mgr.ctx, key).Result(); err != nil {
					return 0, err
				}
			}
			return int64(len(keys)), nil
		}
	}
	return 0, ErrConnNotExist
}