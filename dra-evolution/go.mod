module github.com/kubernetes-sigs/wg-device-management/dra-evolution

go 1.22.1

replace github.com/kubernetes-sigs/wg-device-management/nv-partitionable-resources => ../nv-partitionable-resources

require (
	github.com/NVIDIA/go-nvml v0.12.0-5
	github.com/google/cel-go v0.20.1
	github.com/kubernetes-sigs/wg-device-management/nv-partitionable-resources v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.9.0
	k8s.io/api v0.30.0
	k8s.io/apimachinery v0.30.0
	sigs.k8s.io/kubebuilder-declarative-pattern/mockkubeapiserver v0.0.0-20240404191132-83bd9c05741b
	sigs.k8s.io/yaml v1.4.0
)

require (
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/NVIDIA/go-nvlib v0.3.0 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stoewer/go-strcase v1.2.0 // indirect
	golang.org/x/exp v0.0.0-20231110203233-9a3e6036ecaa // indirect
	golang.org/x/net v0.23.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20230803162519-f966b187b2e5 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230803162519-f966b187b2e5 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/klog/v2 v2.120.1 // indirect
	k8s.io/utils v0.0.0-20240423183400-0849a56e8f22 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.4.1 // indirect
)
