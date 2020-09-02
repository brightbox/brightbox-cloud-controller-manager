module github.com/brightbox/brightbox-cloud-controller-manager

go 1.13

require (
	github.com/brightbox/gobrightbox v0.4.2
	github.com/brightbox/k8ssdk v0.2.1
	github.com/evanphx/json-patch v4.2.0+incompatible // indirect
	github.com/go-test/deep v1.0.4
	k8s.io/api v0.16.14
	k8s.io/apimachinery v0.16.14
	k8s.io/cloud-provider v0.16.14
	k8s.io/component-base v0.16.14
	k8s.io/klog v1.0.0
	k8s.io/kubernetes v1.16.14
)

replace k8s.io/api => k8s.io/api v0.16.14

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.16.14

replace k8s.io/apimachinery => k8s.io/apimachinery v0.16.15-rc.0

replace k8s.io/apiserver => k8s.io/apiserver v0.16.14

replace k8s.io/cli-runtime => k8s.io/cli-runtime v0.16.14

replace k8s.io/client-go => k8s.io/client-go v0.16.14

replace k8s.io/cloud-provider => k8s.io/cloud-provider v0.16.14

replace k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.16.14

replace k8s.io/code-generator => k8s.io/code-generator v0.16.14-rc.0

replace k8s.io/component-base => k8s.io/component-base v0.16.14

replace k8s.io/cri-api => k8s.io/cri-api v0.16.14-rc.0

replace k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.16.14

replace k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.16.14

replace k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.16.14

replace k8s.io/kube-proxy => k8s.io/kube-proxy v0.16.14

replace k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.16.14

replace k8s.io/kubectl => k8s.io/kubectl v0.16.14

replace k8s.io/kubelet => k8s.io/kubelet v0.16.14

replace k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.16.14

replace k8s.io/metrics => k8s.io/metrics v0.16.14

replace k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.16.14

replace vbom.ml/util => github.com/fvbommel/util v0.0.0-20160121211510-db5cfe13f5cc
