package main

import "fmt"

//go:generate go run generator.go -path=public/release/models/acl,public/release/models/types,public/release/models,public/third_party/ietf -output_file=../internal/pkg/openconfig/acl/acl.go -package_name=openconfig -shorten_enum_leaf_names=true -compress_paths=true -exclude_modules=ietf-interfaces public/release/models/acl/openconfig-acl.yang

func main() {
	fmt.Println("Done.")
}
