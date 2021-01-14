module github.com/brightbox/brightbox-cloud-controller-manager

go 1.15

require (
	github.com/brightbox/gobrightbox v0.5.7
	github.com/brightbox/k8ssdk v0.6.1
	github.com/go-test/deep v1.0.7
	github.com/spf13/cobra v1.1.1
	k8s.io/api v0.20.2
	k8s.io/apimachinery v0.20.2
	k8s.io/cloud-provider v0.20.2
	k8s.io/component-base v0.20.2
	k8s.io/controller-manager v0.20.2
	k8s.io/klog/v2 v2.4.0
	k8s.io/kubernetes v1.20.2
)

replace vbom.ml/util => github.com/fvbommel/util v0.0.0-20160121211510-db5cfe13f5cc

replace k8s.io/api => k8s.io/api v0.20.2

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.20.2

replace k8s.io/apimachinery => k8s.io/apimachinery v0.21.0-alpha.0

replace k8s.io/apiserver => k8s.io/apiserver v0.20.2

replace k8s.io/cli-runtime => k8s.io/cli-runtime v0.20.2

replace k8s.io/client-go => k8s.io/client-go v0.20.2

replace k8s.io/cloud-provider => k8s.io/cloud-provider v0.20.2

replace k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.20.2

replace k8s.io/code-generator => k8s.io/code-generator v0.20.3-rc.0

replace k8s.io/component-base => k8s.io/component-base v0.20.2

replace k8s.io/component-helpers => k8s.io/component-helpers v0.20.2

replace k8s.io/controller-manager => k8s.io/controller-manager v0.20.2

replace k8s.io/cri-api => k8s.io/cri-api v0.20.3-rc.0

replace k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.20.2

replace k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.20.2

replace k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.20.2

replace k8s.io/kube-proxy => k8s.io/kube-proxy v0.20.2

replace k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.20.2

replace k8s.io/kubectl => k8s.io/kubectl v0.20.2

replace k8s.io/kubelet => k8s.io/kubelet v0.20.2

replace k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.20.2

replace k8s.io/metrics => k8s.io/metrics v0.20.2

replace k8s.io/mount-utils => k8s.io/mount-utils v0.20.2-rc.0

replace k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.20.2
