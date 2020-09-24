package eagle

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	frameworkruntime "k8s.io/kubernetes/pkg/scheduler/framework/runtime"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
)

const (
	// Name of eagle plugin
	Name              = "eagle"
	preFilterStateKey = "PreFilter" + Name
)

// Eagle is our custom plugin.
type Eagle struct {
	args   *Args
	handle framework.FrameworkHandle
	resourceAllocationScorer
}

var (
	_ framework.PreFilterPlugin = &Eagle{}
	_ framework.FilterPlugin    = &Eagle{}
	_ framework.ScorePlugin     = &Eagle{}
	_ framework.ScoreExtensions = &Eagle{}
)

// Args maintains basic args for running a scheduler.
type Args struct {
	KubeConfig string `json:"kubeconfig,omitempty"`
	Master     string `json:"master,omitempty"`
}

// Name returns the name of Eagle plugin.
func (e Eagle) Name() string {
	return Name
}

// New initializes a new plugin and return it.
func New(configuration runtime.Object, f framework.FrameworkHandle) (framework.Plugin, error) {
	args := &Args{}
	if err := frameworkruntime.DecodeInto(configuration, args); err != nil {
		return nil, err
	}

	resToWeightMap := make(resourceToWeightMap)
	resToWeightMap["cpu"] = 1
	resToWeightMap["memory"] = 1

	klog.V(3).Infof("get plugin config args: %+v", args)
	return &Eagle{
		args:   args,
		handle: f,
		resourceAllocationScorer: resourceAllocationScorer{
			Name:                "NodeResourcesEagleAllocated",
			scorer:              eagleResourceScorer(),
			resourceToWeightMap: resToWeightMap,
		},
	}, nil
}

// PreFilter invoked at the prefilter extension point.
// It returns the summary of containers requests.
func (e *Eagle) PreFilter(ctx context.Context, cycleState *framework.CycleState, pod *v1.Pod) *framework.Status {
	cycleState.Write(preFilterStateKey, computePodResourceRequest(pod))
	return nil
}

// PreFilterExtensions returns prefilter extensions, pod add and remove.
func (e *Eagle) PreFilterExtensions() framework.PreFilterExtensions {
	return nil
}

// Filter invoked at the filter extension point.
// Checks if a node has sufficient resources, such as cpu, memory, gpu, opaque int resources etc to run a pod.
// It returns a list of insufficient resources, if empty, then the node has all the resources requested by the pod.
func (e *Eagle) Filter(ctx context.Context, cycleState *framework.CycleState, pod *v1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
	s, err := getPreFilterState(preFilterStateKey, cycleState)
	if err != nil {
		return framework.NewStatus(framework.Error, err.Error())
	}
	insufficietnResources := fitsRequest(s, nodeInfo)

	if len(insufficietnResources) != 0 {
		// keep all failure reasons.
		failureReasons := make([]string, 0, len(insufficietnResources))
		for _, r := range insufficietnResources {
			failureReasons = append(failureReasons, r.Reason)
		}
		return framework.NewStatus(framework.Unschedulable, failureReasons...)
	}
	return framework.NewStatus(framework.Success, "")
}

// Score rank nodes that passed the filtering phase, and it is invoked at the Score extension point.
// it call the function in allocationScorer which is implemented below.
func (e *Eagle) Score(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) (int64, *framework.Status) {
	nodeInfo, err := e.handle.SnapshotSharedLister().NodeInfos().Get(nodeName)
	if err != nil || nodeInfo.Node() == nil {
		return 0, framework.NewStatus(framework.Error, fmt.Sprintf("getting node %q from Snapshot: %v, node is nil: %v", nodeName, err, nodeInfo.Node() == nil))
	}

	return e.score(pod, nodeInfo)
}

// ScoreExtensions is defined in interface ScorePlugin and
// return a ScoreExtensions interafce if it has been implemented.
func (e *Eagle) ScoreExtensions() framework.ScoreExtensions {
	return e
}

// NormalizeScore TODO
func (e *Eagle) NormalizeScore(ctx context.Context, state *framework.CycleState, pod *v1.Pod, score framework.NodeScoreList) *framework.Status {
	return nil
}