package handler

import "github.com/openconfig/gnmi/proto/gnmi"

func FetchPathKey(r *gnmi.GetRequest) map[string]string {
    kvs := map[string]string{}
    for _, path := range r.Path {
        for _, elem := range path.Elem {
            if elem.Key != nil {
                for k, v := range elem.Key {
                    kvs[k] = v
                }
            }
        }
    }
    return kvs
}