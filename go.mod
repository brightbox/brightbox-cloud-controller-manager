module github.com/brightbox/brightbox-cloud-controller-manager

require (
	github.com/aws/aws-sdk-go-v2 v0.13.0
	github.com/brightbox/gobrightbox v0.4.2
	github.com/go-test/deep v1.0.1
	golang.org/x/oauth2 v0.0.0-20190402181905-9f3314589c9a
	k8s.io/api v0.15.8
	k8s.io/apimachinery v0.15.8
	k8s.io/cloud-provider v0.15.8
	k8s.io/component-base v0.15.8
	k8s.io/klog v0.3.1
	k8s.io/kubernetes v1.15.8
)

go 1.12

replace k8s.io/api => k8s.io/api v0.15.8

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.15.8

replace k8s.io/apimachinery => k8s.io/apimachinery v0.15.9-beta.0

replace k8s.io/apiserver => k8s.io/apiserver v0.15.8

replace k8s.io/cli-runtime => k8s.io/cli-runtime v0.15.8

replace k8s.io/client-go => k8s.io/client-go v0.15.8

replace k8s.io/cloud-provider => k8s.io/cloud-provider v0.15.8

replace k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.15.8

replace k8s.io/code-generator => k8s.io/code-generator v0.15.9-beta.0

replace k8s.io/component-base => k8s.io/component-base v0.15.8

replace k8s.io/cri-api => k8s.io/cri-api v0.15.9-beta.0

replace k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.15.8

replace k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.15.8

replace k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.15.8

replace k8s.io/kube-proxy => k8s.io/kube-proxy v0.15.8

replace k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.15.8

replace k8s.io/kubelet => k8s.io/kubelet v0.15.8

replace k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.15.8

replace k8s.io/metrics => k8s.io/metrics v0.15.8

replace k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.15.8
