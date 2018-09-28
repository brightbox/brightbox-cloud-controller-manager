module github.com/brightbox/brightbox-cloud-controller-manager

require (
	bitbucket.org/ww/goautoneg v0.0.0-20120707110453-75cd24fc2f2c
	github.com/NYTimes/gziphandler v1.0.1
	github.com/PuerkitoBio/purell v1.1.0
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578
	github.com/aws/aws-sdk-go-v2 v2.0.0-preview.4+incompatible
	github.com/beorn7/perks v0.0.0-20180321164747-3a771d992973
	github.com/brightbox/gobrightbox v0.0.1
	github.com/coreos/etcd v3.3.9+incompatible
	github.com/coreos/go-semver v0.2.0
	github.com/coreos/go-systemd v0.0.0-20180511133405-39ca1b05acc7
	github.com/davecgh/go-spew v1.1.1
	github.com/docker/distribution v0.0.0-20170726174610-edc3ab29cdff
	github.com/docker/docker v1.13.1 // indirect
	github.com/elazarl/go-bindata-assetfs v1.0.0
	github.com/emicklei/go-restful v2.8.0+incompatible
	github.com/emicklei/go-restful-swagger12 v0.0.0-20170208215640-dcef7f557305
	github.com/evanphx/json-patch v3.0.0+incompatible
	github.com/ghodss/yaml v1.0.0
	github.com/go-ini/ini v1.38.1
	github.com/go-openapi/jsonpointer v0.0.0-20180322222829-3a0015ad55fa
	github.com/go-openapi/jsonreference v0.0.0-20180322222742-3fb327e6747d
	github.com/go-openapi/spec v0.0.0-20180801175345-384415f06ee2
	github.com/go-openapi/swag v0.0.0-20180715190254-becd2f08beaf
	github.com/go-test/deep v1.0.1
	github.com/gogo/protobuf v1.1.1
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/groupcache v0.0.0-20180513044358-24b0969c4cb7
	github.com/golang/protobuf v1.2.0
	github.com/google/btree v0.0.0-20180813153112-4030bb1f1f0c
	github.com/google/gofuzz v0.0.0-20170612174753-24818f796faf
	github.com/googleapis/gnostic v0.2.0
	github.com/gregjones/httpcache v0.0.0-20180305231024-9cad4c3443a7
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/hashicorp/golang-lru v0.0.0-20180201235237-0fb14efe8c47
	github.com/imdario/mergo v0.3.6
	github.com/inconshreveable/mousetrap v1.0.0
	github.com/jmespath/go-jmespath v0.0.0-20160202185014-0b12d6b521d8
	github.com/json-iterator/go v0.0.0-20180612202835-f2b4162afba3
	github.com/mailru/easyjson v0.0.0-20180730094502-03f2033d19d5
	github.com/matttproud/golang_protobuf_extensions v1.0.1
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd
	github.com/modern-go/reflect2 v0.0.0-20180701023420-4b7aa43c6742
	github.com/opencontainers/go-digest v1.0.0-rc1
	github.com/pborman/uuid v0.0.0-20170612153648-e790cca94e6c
	github.com/petar/GoLLRB v0.0.0-20130427215148-53be0d36a84c
	github.com/peterbourgon/diskv v2.0.1+incompatible
	github.com/prometheus/client_golang v0.8.0
	github.com/prometheus/client_model v0.0.0-20180712105110-5c3871d89910
	github.com/prometheus/common v0.0.0-20180801064454-c7de2306084e
	github.com/prometheus/procfs v0.0.0-20180725123919-05ee40e3a273
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v0.0.0-20180821114517-d929dcbb1086
	github.com/tomnomnom/linkheader v0.0.0-20160328204959-6953a30d4443
	github.com/ugorji/go v1.1.1
	golang.org/x/crypto v0.0.0-20180820150726-614d502a4dac
	golang.org/x/net v0.0.0-20180821023952-922f4815f713
	golang.org/x/oauth2 v0.0.0-20180821212333-d2e6202438be
	golang.org/x/sys v0.0.0-20180821140842-3b58ed4ad339
	golang.org/x/text v0.3.0
	golang.org/x/time v0.0.0-20180412165947-fbb02b2291d2
	google.golang.org/appengine v1.1.0
	google.golang.org/genproto v0.0.0-20180817151627-c66870c02cf8
	google.golang.org/grpc v1.14.0
	gopkg.in/inf.v0 v0.9.1
	gopkg.in/natefinch/lumberjack.v2 v2.0.0-20170531160350-a96e63847dc3
	gopkg.in/square/go-jose.v2 v2.1.8
	gopkg.in/yaml.v2 v2.2.1
	k8s.io/api v0.0.0-20180925152912-a191abe0b71e
	k8s.io/apiextensions-apiserver v0.0.0-20180808065829-408db4a50408
	k8s.io/apimachinery v0.0.0-20180927151612-c6dd271be006
	k8s.io/apiserver v0.0.0-20180928074203-ec2b99f30258
	k8s.io/client-go v9.0.0+incompatible
	k8s.io/csi-api v0.0.0-20180925155520-31ae05d8096d // indirect
	k8s.io/kube-controller-manager v0.0.0-20180909204853-481b4013cbb6 // indirect
	k8s.io/kube-openapi v0.0.0-20180731170545-e3762e86a74c
	k8s.io/kubernetes v1.12.0
	k8s.io/utils v0.0.0-20180918230422-cd34563cd63c // indirect
)
