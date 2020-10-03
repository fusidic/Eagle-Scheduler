module github.com/fusidic/Eagle-Scheduler

go 1.15

require (
	github.com/spf13/cobra v1.0.0
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/apiserver v0.19.2
	k8s.io/component-base v0.19.2
	k8s.io/klog/v2 v2.3.0
	k8s.io/kubernetes v1.19.2
)

replace (
	k8s.io/api => /root/workspace/kubernetes/staging/src/k8s.io/api
	k8s.io/apiextensions-apiserver => /root/workspace/kubernetes/staging/src/k8s.io/apiextensions-apiserver
	k8s.io/apimachinery => /root/workspace/kubernetes/staging/src/k8s.io/apimachinery
	k8s.io/apiserver => /root/workspace/kubernetes/staging/src/k8s.io/apiserver
	k8s.io/cli-runtime => /root/workspace/kubernetes/staging/src/k8s.io/cli-runtime
	k8s.io/client-go => /root/workspace/kubernetes/staging/src/k8s.io/client-go
	k8s.io/cloud-provider => /root/workspace/kubernetes/staging/src/k8s.io/cloud-provider
	k8s.io/cluster-bootstrap => /root/workspace/kubernetes/staging/src/k8s.io/cluster-bootstrap
	k8s.io/code-generator => /root/workspace/kubernetes/staging/src/k8s.io/code-generator
	k8s.io/component-base => /root/workspace/kubernetes/staging/src/k8s.io/component-base
	k8s.io/cri-api => /root/workspace/kubernetes/staging/src/k8s.io/cri-api
	k8s.io/csi-translation-lib => /root/workspace/kubernetes/staging/src/k8s.io/csi-translation-lib
	k8s.io/kube-aggregator => /root/workspace/kubernetes/staging/src/k8s.io/kube-aggregator
	k8s.io/kube-controller-manager => /root/workspace/kubernetes/staging/src/k8s.io/kube-controller-manager
	k8s.io/kube-proxy => /root/workspace/kubernetes/staging/src/k8s.io/kube-proxy
	k8s.io/kube-scheduler => /root/workspace/kubernetes/staging/src/k8s.io/kube-scheduler
	k8s.io/kubectl => /root/workspace/kubernetes/staging/src/k8s.io/kubectl
	k8s.io/kubelet => /root/workspace/kubernetes/staging/src/k8s.io/kubelet
	k8s.io/kubernetes => /root/workspace/kubernetes
	k8s.io/legacy-cloud-providers => /root/workspace/kubernetes/staging/src/k8s.io/legacy-cloud-providers
	k8s.io/metrics => /root/workspace/kubernetes/staging/src/k8s.io/metrics
	k8s.io/mount-utils => /root/workspace/kubernetes/staging/src/k8s.io/mount-utils
	k8s.io/sample-apiserver => /root/workspace/kubernetes/staging/src/k8s.io/sample-apiserver
)
