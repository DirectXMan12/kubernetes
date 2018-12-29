/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// This file provides abstractions for setting the provider (e.g., prometheus)
// of metrics.

package cache

import (
	"k8s.io/client-go/tools/metrics"
)

var (
	listsTotal = metrics.NewCounterVec(metrics.Counter{
		Name:      "reflector_lists_total",
		Help:      "Total number of API lists done by the reflectors",
	}, "name")

	listsDuration = metrics.NewHistogramVec(metrics.Histogram{
		Name:      "reflector_list_duration_seconds",
		Help:      "How long an API list takes to return and decode for the reflectors",
	}, "name")

	itemsPerList = metrics.NewHistogramVec(metrics.Histogram{
		Name:      "reflector_items_per_list",
		Help:      "How many items an API list returns to the reflectors",
	}, "name")

	watchesTotal = metrics.NewCounterVec(metrics.Counter{
		Name:      "reflector_watches_total",
		Help:      "Total number of API watches done by the reflectors",
	}, "name")

	shortWatchesTotal = metrics.NewCounterVec(metrics.Counter{
		Name:      "reflector_short_watches_total",
		Help:      "Total number of short API watches done by the reflectors",
	}, "name")

	watchDuration = metrics.NewHistogramVec(metrics.Histogram{
		Name:      "reflector_watch_duration_seconds",
		Help:      "How long an API watch takes to return and decode for the reflectors",
	}, "name")

	itemsPerWatch = metrics.NewHistogramVec(metrics.Histogram{
		Name:      "reflector_items_per_watch",
		Help:      "How many items an API watch returns to the reflectors",
	}, "name")

	lastResourceVersion = metrics.NewGaugeVec(metrics.Gauge{
		Name:      "reflector_last_resource_version",
		Help:      "Last resource version seen for the reflectors",
	}, "name")
)

func init() {
	listsTotal.MustRegisterIn(metrics.DefaultRegistry())
	listsDuration.MustRegisterIn(metrics.DefaultRegistry())
	itemsPerList.MustRegisterIn(metrics.DefaultRegistry())
	watchesTotal.MustRegisterIn(metrics.DefaultRegistry())
	shortWatchesTotal.MustRegisterIn(metrics.DefaultRegistry())
	watchDuration.MustRegisterIn(metrics.DefaultRegistry())
	itemsPerWatch.MustRegisterIn(metrics.DefaultRegistry())
	lastResourceVersion.MustRegisterIn(metrics.DefaultRegistry())
}

type reflectorMetrics struct {
	numberOfLists       metrics.CounterImpl
	listDuration        metrics.HistogramImpl
	numberOfItemsInList metrics.HistogramImpl

	numberOfWatches      metrics.CounterImpl
	numberOfShortWatches metrics.CounterImpl
	watchDuration        metrics.HistogramImpl
	numberOfItemsInWatch metrics.HistogramImpl

	lastResourceVersion metrics.GaugeImpl
}

func newReflectorMetrics(name string) *reflectorMetrics {
	var ret *reflectorMetrics
	if len(name) == 0 {
		return ret
	}
	return &reflectorMetrics{
		numberOfLists:        listsTotal.WithLabelValues(name),
		listDuration:         listsDuration.WithLabelValues(name),
		numberOfItemsInList:  itemsPerList.WithLabelValues(name),
		numberOfWatches:      watchesTotal.WithLabelValues(name),
		numberOfShortWatches: shortWatchesTotal.WithLabelValues(name),
		watchDuration:        watchDuration.WithLabelValues(name),
		numberOfItemsInWatch: itemsPerWatch.WithLabelValues(name),
		lastResourceVersion:  lastResourceVersion.WithLabelValues(name),
	}
}
