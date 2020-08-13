package swsssdk

import (
    "context"
    "fmt"
    "github.com/go-redis/redis/v8"
    "time"
)

// 多客户端连接管理
type manager struct {
    clients map[int]*redis.Client
    subs    map[int]*redis.PubSub
    ctx     context.Context
}

func (mgr *manager) close() {
    for id, client := range mgr.clients {
        client.Close()
        delete(mgr.clients, id)
    }
}

func (mgr *manager) connect(id int, hostname string, port int) bool {
    if _, ok := mgr.clients[id]; ok {
        return true
    }
    client := redis.NewClient(&redis.Options{
        Addr:       fmt.Sprintf("%s:%d", hostname, port),
        Password:   "",
        DB:         id,
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

func (mgr *manager) disconnect(id int) {
    if client, ok := mgr.clients[id]; ok {
        client.Close()
        delete(mgr.clients, id)
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
        return client.Expire(mgr.ctx, key, time.Duration(timeout_sec)*time.Second).Result()
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

func (mgr *manager) hset(id int, key string, value interface{}) (int64, error) {
    if client, ok := mgr.clients[id]; ok {
        return client.HSet(mgr.ctx, key, value).Result()
    }
    return 0, ErrConnNotExist
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
