// This is a submodule to isolate k8s.io/code-generator from k8s.io/{api,apimachinery,client-go} dependencies in generated code

module k8s.io/code-generator/examples

go 1.16

require (
	k8s.io/api v0.0.0
	k8s.io/apimachinery v0.0.0
	k8s.io/client-go v0.0.0
	k8s.io/kube-openapi v0.0.8-gnostic.0.20220425145412-55d630839892
)

replace (
	k8s.io/api => ../../api
	k8s.io/apimachinery => ../../apimachinery
	k8s.io/client-go => ../../client-go
)

replace sigs.k8s.io/json => github.com/liggitt/json v0.0.0-20211020163728-48258682683b

replace k8s.io/kube-openapi => github.com/jefftree/kube-openapi v0.0.8-gnostic.0.20220425145412-55d630839892
