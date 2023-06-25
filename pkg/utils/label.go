package utils

import "k8s.io/apimachinery/pkg/labels"

func MergeLabels(sets ...map[string]string) labels.Set {
	merged := labels.Set{}
	for _, set := range sets {
		merged = labels.Merge(merged, set)
	}
	return merged
}
