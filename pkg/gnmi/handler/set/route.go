package set

import my_gnmi "gnmi_server/pkg/gnmi"

func SetServeMux() *my_gnmi.SetServeMux {
    mux := my_gnmi.NewSetServeMux()
    route(mux)
    return mux
}

func route(mux *my_gnmi.SetServeMux) {

}