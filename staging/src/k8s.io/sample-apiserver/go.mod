// This is a generated file. Do not edit directly.

module k8s.io/sample-apiserver

go 1.16

require (
	github.com/NYTimes/gziphandler v1.1.1 // indirect
	github.com/google/gofuzz v1.1.0
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spf13/cobra v1.2.1
	google.golang.org/genproto v0.0.0-20210831024726-fe130286e0e2 // indirect
	k8s.io/apimachinery v0.0.0
	k8s.io/apiserver v0.0.0
	k8s.io/client-go v0.0.0
	k8s.io/code-generator v0.0.0
	k8s.io/component-base v0.0.0
	k8s.io/kube-openapi v0.0.0-20211115234752-e816edb12b65
	k8s.io/utils v0.0.0-20211208161948-7d6a63dca704
)

replace (
	github.com/pquerna/cachecontrol => github.com/pquerna/cachecontrol v0.0.0-20171018203845-0dec1b30a021
	k8s.io/api => ../api
	k8s.io/apimachinery => ../apimachinery
	k8s.io/apiserver => ../apiserver
	k8s.io/client-go => ../client-go
	k8s.io/code-generator => ../code-generator
	k8s.io/component-base => ../component-base
	k8s.io/sample-apiserver => ../sample-apiserver
)
