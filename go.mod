module github.com/brightbox/brightbox-cloud-controller-manager

go 1.13

require (
	github.com/beorn7/perks v1.0.0 // indirect
	github.com/brightbox/gobrightbox v0.4.2
	github.com/brightbox/k8ssdk v0.2.0
	github.com/coreos/go-systemd v0.0.0-20190321100706-95778dfbb74e // indirect
	github.com/go-openapi/jsonreference v0.19.3 // indirect
	github.com/go-openapi/spec v0.19.3 // indirect
	github.com/go-test/deep v1.0.4
	github.com/golang/protobuf v1.3.2 // indirect
	github.com/google/btree v1.0.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.1-0.20190118093823-f849b5445de4 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.9.5 // indirect
	github.com/json-iterator/go v1.1.8 // indirect
	github.com/mailru/easyjson v0.7.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/onsi/ginkgo v1.10.1 // indirect
	github.com/onsi/gomega v1.7.0 // indirect
	github.com/soheilhy/cmux v0.1.4 // indirect
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	go.uber.org/atomic v1.3.2 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.10.0 // indirect
	golang.org/x/sys v0.0.0-20190826190057-c7b8b68b1456 // indirect
	google.golang.org/grpc v1.23.1 // indirect
	gopkg.in/yaml.v2 v2.2.7 // indirect
	k8s.io/api v0.17.0
	k8s.io/apimachinery v0.17.0
	k8s.io/cloud-provider v0.17.0
	k8s.io/component-base v0.17.0
	k8s.io/klog v1.0.0
	k8s.io/kube-openapi v0.0.0-20191107075043-30be4d16710a // indirect
	k8s.io/kubernetes v1.16.5
	k8s.io/utils v0.0.0-20191114184206-e782cd3c129f // indirect
)

replace k8s.io/api => k8s.io/api v0.16.5

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.16.5

replace k8s.io/apimachinery => k8s.io/apimachinery v0.16.6-beta.0

replace k8s.io/apiserver => k8s.io/apiserver v0.16.5

replace k8s.io/cli-runtime => k8s.io/cli-runtime v0.16.5

replace k8s.io/client-go => k8s.io/client-go v0.16.5

replace k8s.io/cloud-provider => k8s.io/cloud-provider v0.16.5

replace k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.16.5

replace k8s.io/code-generator => k8s.io/code-generator v0.16.6-beta.0

replace k8s.io/component-base => k8s.io/component-base v0.16.5

replace k8s.io/cri-api => k8s.io/cri-api v0.16.6-beta.0

replace k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.16.5

replace k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.16.5

replace k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.16.5

replace k8s.io/kube-proxy => k8s.io/kube-proxy v0.16.5

replace k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.16.5

replace k8s.io/kubectl => k8s.io/kubectl v0.16.5

replace k8s.io/kubelet => k8s.io/kubelet v0.16.5

replace k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.16.5

replace k8s.io/metrics => k8s.io/metrics v0.16.5

replace k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.16.5
