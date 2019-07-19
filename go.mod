module github.com/brightbox/brightbox-cloud-controller-manager

require (
	github.com/aws/aws-sdk-go-v2 v2.0.0-preview.4+incompatible
	github.com/brightbox/gobrightbox v0.4.2
	github.com/go-ini/ini v1.42.0 // indirect
	github.com/go-test/deep v1.0.1
	golang.org/x/oauth2 v0.0.0-20190402181905-9f3314589c9a
	gopkg.in/ini.v1 v1.42.0 // indirect
	k8s.io/api v0.0.0
	k8s.io/apimachinery v0.0.0
	k8s.io/cloud-provider v0.0.0
	k8s.io/component-base v0.0.0
	k8s.io/klog v0.3.1
	k8s.io/kubernetes v1.15.1
)

replace k8s.io/api => k8s.io/api v0.0.0-20190718183219-b59d8169aab5

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190718185103-d1ef975d28ce

replace k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190612205821-1799e75a0719

replace k8s.io/apiserver => k8s.io/apiserver v0.0.0-20190718184206-a1aa83af71a7

replace k8s.io/cli-runtime => k8s.io/cli-runtime v0.0.0-20190718185405-0ce9869d0015

replace k8s.io/client-go => k8s.io/client-go v0.0.0-20190718183610-8e956561bbf5

replace k8s.io/cloud-provider => k8s.io/cloud-provider v0.0.0-20190718190308-f8e43aa19282

replace k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.0.0-20190718190146-f7b0473036f9

replace k8s.io/code-generator => k8s.io/code-generator v0.0.0-20190612205613-18da4a14b22b

replace k8s.io/component-base => k8s.io/component-base v0.0.0-20190718183727-0ececfbe9772

replace k8s.io/cri-api => k8s.io/cri-api v0.0.0-20190531030430-6117653b35f1

replace k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.0.0-20190718190424-bef8d46b95de

replace k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.0.0-20190718184434-a064d4d1ed7a

replace k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.0.0-20190718190030-ea930fedc880

replace k8s.io/kube-proxy => k8s.io/kube-proxy v0.0.0-20190718185641-5233cb7cb41e

replace k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.0.0-20190718185913-d5429d807831

replace k8s.io/kubelet => k8s.io/kubelet v0.0.0-20190718185757-9b45f80d5747

replace k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.0.0-20190718190548-039b99e58dbd

replace k8s.io/metrics => k8s.io/metrics v0.0.0-20190718185242-1e1642704fe6

replace k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.0.0-20190718184639-baafa86838c0
