package set

import (
    "gnmi_server/pkg/gnmi"
)

func SetServeMux() *gnmi.SetServeMux {
    mux := gnmi.NewSetServeMux()
    route(mux)
    return mux
}

func route(mux *gnmi.SetServeMux) {

}
