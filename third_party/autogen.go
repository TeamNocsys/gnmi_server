package main

import "fmt"

//go:generate go run generator.go -path=public/release/models/acl,public/release/models/types,public/release/models,public/third_party/ietf -output_file=../internal/pkg/openconfig/acl/acl.go -package_name=acl -shorten_enum_leaf_names=true -compress_paths=true -exclude_modules=ietf-interfaces public/release/models/acl/openconfig-acl.yang
//go:generate go run generator.go -path=public/release/models,public/third_party/ietf -output_file=../internal/pkg/openconfig/interfaces/interfaces.go -package_name=interfaces -shorten_enum_leaf_names=true -compress_paths=true -exclude_modules=ietf-interfaces public/release/models/interfaces/openconfig-interfaces.yang
//go:generate go run generator.go -path=public/release/models,public/third_party/ietf -output_file=../internal/pkg/openconfig/lldp/lldp.go -package_name=lldp -shorten_enum_leaf_names=true -compress_paths=true -exclude_modules=ietf-interfaces public/release/models/lldp/openconfig-lldp.yang
//go:generate go run generator.go -path=public/release/models,public/third_party/ietf -output_file=../internal/pkg/openconfig/local-routing/local_routing.go -package_name=local_routing -shorten_enum_leaf_names=true -compress_paths=true -exclude_modules=ietf-interfaces public/release/models/local-routing/openconfig-local-routing.yang
//go:generate go run generator.go -path=public/release/models,public/third_party/ietf -output_file=../internal/pkg/openconfig/network-instance/network_instance.go -package_name=network_instance -shorten_enum_leaf_names=true -compress_paths=true -exclude_modules=ietf-interfaces public/release/models/network-instance/openconfig-network-instance.yang
//go:generate go run generator.go -path=public/release/models,public/third_party/ietf -output_file=../internal/pkg/openconfig/platform/platform.go -package_name=platform -shorten_enum_leaf_names=true -compress_paths=true -exclude_modules=ietf-interfaces public/release/models/platform/openconfig-platform.yang public/release/models/platform/openconfig-platform-fan.yang
//go:generate go run generator.go -path=public/release/models,public/third_party/ietf -output_file=../internal/pkg/openconfig/qos/qos.go -package_name=qos -shorten_enum_leaf_names=true -compress_paths=true -exclude_modules=ietf-interfaces public/release/models/qos/openconfig-qos.yang
//go:generate go run generator.go -path=public/release/models,public/third_party/ietf -output_file=../internal/pkg/openconfig/system/system.go -package_name=system -shorten_enum_leaf_names=true -compress_paths=true -exclude_modules=ietf-interfaces public/release/models/system/openconfig-system.yang
//go:generate go run generator.go -path=public/release/models,public/third_party/ietf -output_file=../internal/pkg/openconfig/vlan/vlan.go -package_name=vlan -shorten_enum_leaf_names=true -compress_paths=true -exclude_modules=ietf-interfaces public/release/models/vlan/openconfig-vlan.yang

func main() {
	fmt.Println("Done.")
}
