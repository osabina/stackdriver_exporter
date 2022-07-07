// Copyright 2020 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package collectors

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type CacheEntry struct {
	val float64
	ts  time.Time
}

var (
	counterCache           = make(map[uint64]*CacheEntry)
	counterSentinelStrings = []string{""}
)

func SetCounterSentinelStrings(strs []string) {
	counterSentinelStrings = strs
}

// DeltaIsActuallyCounter returns true if a delta metric should actually be a counter
// based off a match from the monitoring.counter-substring flag.
func DeltaIsActuallyCounter(fq_name string) bool {
	// For now we split on _ to get the _ separated words as we have found the word
	// "count" in the metric name to be indicative of this delta-as-counter situation, but
	// we may eventually want to simply look for a substring or RE for flexibility/specificity.
	words := strings.Split(fq_name, "_")
	for _, sentinel := range counterSentinelStrings {
		for _, value := range words {
			if value == sentinel {
				return true
			}
		}
	}
	return false
}

// GetCounterValue retrieves the previously stored value for a metric.
func GetCounterValue(key uint64) float64 {
	if entry, ok := counterCache[key]; ok {
		return entry.val
	}
	newEntry := CacheEntry{ts: time.Time{}, val: 0}
	counterCache[key] = &newEntry
	return 0
}

// SetCounterValue sets the current value for a metric.
func SetCounterValue(key uint64, val float64, ts time.Time) bool {
	if entry, ok := counterCache[key]; ok {
		if ts.After(entry.ts) {
			newEntry := CacheEntry{ts: ts, val: val}
			counterCache[key] = &newEntry
			return true
		}
	}

	return false
}

// SerializeLabels returns a sorted delimiter-separated list of label key-values.
func SerializeLabels(keys []string, values []string) string {
	labels := make(map[string]string)
	labelstrs := []string{}
	keysCopy := append([]string{}, keys...)
	for i := range keys {
		labels[keys[i]] = values[i]
	}
	sort.Strings(keysCopy)
	for _, k := range keysCopy {
		labelstrs = append(labelstrs, fmt.Sprintf("%s:%s", k, labels[k]))
	}
	return strings.Join(labelstrs, "|")
}

// GetCacheKey returns the hash value of the cache key, comprised of the metric name and associated labels.
func GetCacheKey(fq_name string, keys []string, values []string) uint64 {
	ptextkey := fq_name + "-" + SerializeLabels(keys, values)
	h := hashNew()
	h = hashAdd(h, ptextkey)
	return h
}
