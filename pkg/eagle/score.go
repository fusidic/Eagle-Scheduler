package eagle

import (
	"math"

	v1 "k8s.io/api/core/v1"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
)

func eagleResourceScorer() func(requested, allocable resourceToValueMap, includeVolumes bool, requestedVolumes int, allocatableVolumes int) int64 {
	return func(requested, allocable resourceToValueMap, includeVolumes bool, requestedVolumes int, allocatableVolumes int) int64 {
		cpuFraction := fractionOfCapacity(requested[v1.ResourceCPU], allocable[v1.ResourceCPU])
		memFraction := fractionOfCapacity(requested[v1.ResourceMemory], allocable[v1.ResourceMemory])

		// // didn't take r0 into consideration
		// if ((cpuFraction - 1) * (cpuFraction - 1) + (memFraction - 1) * (memFraction - 1)) < 0.01 {
		// 	return framework.MaxNodeScore
		// }

		// We take two models into account to evaluate the score.
		// Bias can describe the bias with with equal CPU/MEM resource requested.
		// Potential indicates the difference between nodes with same bias value.

		var x, y, potentialValue float64
		biasValue := bias(cpuFraction, memFraction)

		if cpuFraction > memFraction {
			y = cpuFraction
			x = memFraction
			potentialValue = potential(x, y)
		} else if cpuFraction < memFraction {
			y = memFraction
			x = cpuFraction
			potentialValue = potential(x, y)
		} else {
			potentialValue = 1
		}

		finalScore := normalization(biasValue, potentialValue)
		if finalScore > 1 {
			finalScore = 1
		}

		if finalScore < 0 {
			finalScore = 0
		}

		return int64(finalScore * float64(framework.MaxNodeScore))
	}
}

func fractionOfCapacity(requested, capacity int64) float64 {
	if capacity == 0 {
		return 0
	}
	return float64(requested) / float64(capacity)
}

func bias(a, b float64) float64 {
	return 1 - math.Abs(b-a)
}

// By default, we think a less than b
func potential(a, b float64) float64 {
	if a == 1 && b == 1 {
		return 1
	}
	return (1 - b) / (1 - a)
}

func normalization(biasValue, potentialValue float64) float64 {
	return ((biasValue * 10) + potentialValue) / 11
}
