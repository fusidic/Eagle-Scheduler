package eagle

import (
	"math"
)

func eagleResourceScorer() func(requested, allocable resourceToValueMap, includeVolumes bool, requestedVolumes int, allocatableVolumes int) int64 {
	return func(requested, allocable resourceToValueMap, includeVolumes bool, requestedVolumes int, allocatableVolumes int) int64 {
		var nodeScore, weightSum int64
		return 0
	}
}

func bias(a, b float32) float32 {
	return math.Abs(b - a)
}
