package eagle

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/kubernetes/pkg/features"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
)

const (
	// R0 ...
	// r0 ...
	R0 = 0.8
	r0 = 0.9
)

// preFilterState is computed at PreFilter and used at Filter.
type preFilterState struct {
	framework.Resource
}

// InsufficientResource describes what kind of resource limits is hit and caused the pod to not fit the node.
type InsufficientResource struct {
	ResourceName v1.ResourceName
	// explicitly pass a parameter for reason to avoid formatting messages.
	Reason    string
	Requested int64
	Used      int64
	Capacity  int64
}

func computePodResourceRequest(pod *v1.Pod) *preFilterState {
	result := &preFilterState{}
	for _, container := range pod.Spec.Containers {
		result.Add(container.Resources.Requests)
	}

	// take max_resource(sum_pod, any_init_container)
	for _, container := range pod.Spec.InitContainers {
		result.SetMaxResource(container.Resources.Requests)
	}

	// if Overhead is being utilized, add to the total requests for the pod
	if pod.Spec.Overhead != nil && utilfeature.DefaultFeatureGate.Enabled(features.PodOverhead) {
		result.Add(pod.Spec.Overhead)
	}
	return result
}

// Clone the prefilter state.
func (s *preFilterState) Clone() framework.StateData {
	return s
}

func getPreFilterState(preFilterStateKey framework.StateKey, cycleState *framework.CycleState) (*preFilterState, error) {
	c, err := cycleState.Read(preFilterStateKey)
	if err != nil {
		// preFilterState doesn't exist, likely PreFilter wasn't invoked.
		return nil, fmt.Errorf("error reading %q from cycleState: %v", preFilterStateKey, err)
	}
	// The dynamic convert here required an implement of method preFilterState.Clone
	s, ok := c.(*preFilterState)
	if !ok {
		return nil, fmt.Errorf("%+v convert to NodeResourcesFit.preFilterState error", c)
	}
	return s, nil
}

func fitsRequest(podRequest *preFilterState, nodeInfo *framework.NodeInfo) []InsufficientResource {
	insufficientResources := make([]InsufficientResource, 0, 4)

	allowedPodNumber := nodeInfo.Allocatable.AllowedPodNumber
	if len(nodeInfo.Pods)+1 > allowedPodNumber {
		insufficientResources = append(insufficientResources, InsufficientResource{
			v1.ResourcePods,
			"Too many pods",
			1,
			int64(len(nodeInfo.Pods)),
			int64(allowedPodNumber),
		})
	}

	if podRequest.MilliCPU == 0 &&
		podRequest.Memory == 0 &&
		podRequest.EphemeralStorage == 0 &&
		len(podRequest.ScalarResources) == 0 {
		return insufficientResources
	}

	if nodeInfo.Allocatable.MilliCPU < podRequest.MilliCPU+nodeInfo.Requested.MilliCPU {
		insufficientResources = append(insufficientResources, InsufficientResource{
			v1.ResourceCPU,
			"Insufficient cpu",
			podRequest.MilliCPU,
			nodeInfo.Requested.MilliCPU,
			nodeInfo.Allocatable.MilliCPU,
		})
	}

	if nodeInfo.Allocatable.Memory < podRequest.Memory+nodeInfo.Requested.Memory {
		insufficientResources = append(insufficientResources, InsufficientResource{
			v1.ResourceMemory,
			"Insufficient memory",
			podRequest.Memory,
			nodeInfo.Requested.Memory,
			nodeInfo.Allocatable.Memory,
		})
	}

	if reason, ok := fitEagle(podRequest, nodeInfo); !ok {
		insufficientResources = append(insufficientResources, InsufficientResource{
			v1.ResourceCPU,
			reason,
			podRequest.MilliCPU,
			nodeInfo.Requested.MilliCPU,
			nodeInfo.Allocatable.MilliCPU,
		})
	}

	// extension resources check can be implement here
	return insufficientResources
}

func fitEagle(podRequest *preFilterState, nodeInfo *framework.NodeInfo) (string, bool) {
	var cpuRatio, memRatio, y, x float64
	cpuRatio = float64((podRequest.MilliCPU + nodeInfo.Requested.MilliCPU) / nodeInfo.Allocatable.MilliCPU)
	memRatio = float64((podRequest.Memory + nodeInfo.Requested.Memory) / nodeInfo.Allocatable.Memory)
	if cpuRatio > 1 || memRatio > 1 {
		reason := "resource out of limit"
		return reason, false
	}
	if cpuRatio > memRatio {
		y = cpuRatio
		x = memRatio
	} else if cpuRatio < memRatio {
		y = memRatio
		x = cpuRatio
	} else {
		return "ok", true
	}

	if y <= (1-R0) || x >= R0 {
		return "ok", true
	} else if distanceToR0(x, y) {
		return "ok", true
	}

	return "Out of EAGLE bound", false
}

func getMinMax(a, b float64) (float64, float64) {
	if a > b {
		return b, a
	} else if a < b {
		return a, b
	}
	return a, b
}

func distanceToR0(x, y float64) bool {
	distance := (x-R0)*(x-R0) + (y-1+R0)*(y-1+R0)
	if distance <= R0*R0 {
		return true
	}
	return false
}
