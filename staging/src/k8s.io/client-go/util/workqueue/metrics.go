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

package workqueue

import (
	"time"

	"k8s.io/apimachinery/pkg/util/clock"
	"k8s.io/client-go/tools/metrics"
)

var (
	depth = metrics.NewGaugeVec(metrics.Gauge{
		Name:      "depth",
		Help:      "Current depth of workqueue",
	}, "workqueue")

	adds = metrics.NewCounterVec(metrics.Counter{
		Name:      "adds",
		Help:      "Total number of adds handled by workqueue",
	}, "workqueue")

	latency = metrics.NewHistogramVec(metrics.Histogram{
		Name:      "queue_latency",
		Help:      "How long an item stays in workqueue before being requested.",
	}, "workqueue")

	workDuration = metrics.NewHistogramVec(metrics.Histogram{
		Name:      "work_duration",
		Help:      "How long processing an item from workqueue takes.",
	}, "workqueue")

	unfinished = metrics.NewGaugeVec(metrics.Gauge{
		Name:      "unfinished_work_seconds",
		Help: "How many seconds of work has done that " +
			"is in progress and hasn't been observed by work_duration. Large " +
			"values indicate stuck threads. One can deduce the number of stuck " +
			"threads by observing the rate at which this increases.",
	}, "workqueue")

	longestRunning = metrics.NewGaugeVec(metrics.Gauge{
		Name:      "longest_running_processor_microseconds",
		Help: "How many microseconds has the longest running " +
			"processor for been running.",
	}, "workqueue")

	retries = metrics.NewCounterVec(metrics.Counter{
		Name:      "retries",
		Help:      "Total number of retries handled by workqueue",
	}, "workqueue")
)

func init() {
	depth.MustRegisterIn(metrics.DefaultRegistry())
	adds.MustRegisterIn(metrics.DefaultRegistry())
	latency.MustRegisterIn(metrics.DefaultRegistry())
	workDuration.MustRegisterIn(metrics.DefaultRegistry())
	unfinished.MustRegisterIn(metrics.DefaultRegistry())
	longestRunning.MustRegisterIn(metrics.DefaultRegistry())
	retries.MustRegisterIn(metrics.DefaultRegistry())
}

// This file provides abstractions for setting the provider (e.g., prometheus)
// of metrics.

type queueMetrics interface {
	add(item t)
	get(item t)
	done(item t)
	updateUnfinishedWork()
}

// defaultQueueMetrics expects the caller to lock before setting any metrics.
type defaultQueueMetrics struct {
	clock clock.Clock

	// current depth of a workqueue
	depth metrics.GaugeImpl
	// total number of adds handled by a workqueue
	adds metrics.CounterImpl
	// how long an item stays in a workqueue
	latency metrics.HistogramImpl
	// how long processing an item from a workqueue takes
	workDuration         metrics.HistogramImpl
	addTimes             map[t]time.Time
	processingStartTimes map[t]time.Time

	// how long have current threads been working?
	unfinishedWorkSeconds   metrics.GaugeImpl
	longestRunningProcessor metrics.GaugeImpl
}

func (m *defaultQueueMetrics) add(item t) {
	if m == nil {
		return
	}

	m.adds.Inc()
	m.depth.Inc()
	if _, exists := m.addTimes[item]; !exists {
		m.addTimes[item] = m.clock.Now()
	}
}

func (m *defaultQueueMetrics) get(item t) {
	if m == nil {
		return
	}

	m.depth.Dec()
	m.processingStartTimes[item] = m.clock.Now()
	if startTime, exists := m.addTimes[item]; exists {
		m.latency.Observe(m.sinceInMicroseconds(startTime))
		delete(m.addTimes, item)
	}
}

func (m *defaultQueueMetrics) done(item t) {
	if m == nil {
		return
	}

	if startTime, exists := m.processingStartTimes[item]; exists {
		m.workDuration.Observe(m.sinceInMicroseconds(startTime))
		delete(m.processingStartTimes, item)
	}
}

func (m *defaultQueueMetrics) updateUnfinishedWork() {
	// Note that a summary metric would be better for this, but prometheus
	// doesn't seem to have non-hacky ways to reset the summary metrics.
	var total float64
	var oldest float64
	for _, t := range m.processingStartTimes {
		age := m.sinceInMicroseconds(t)
		total += age
		if age > oldest {
			oldest = age
		}
	}
	// Convert to seconds; microseconds is unhelpfully granular for this.
	total /= 1000000
	m.unfinishedWorkSeconds.Set(total)
	m.longestRunningProcessor.Set(oldest) // in microseconds.
}

// Gets the time since the specified start in microseconds.
func (m *defaultQueueMetrics) sinceInMicroseconds(start time.Time) float64 {
	return float64(m.clock.Since(start).Nanoseconds() / time.Microsecond.Nanoseconds())
}

type retryMetrics interface {
	retry()
}

type defaultRetryMetrics struct {
	retries metrics.CounterImpl
}

func (m *defaultRetryMetrics) retry() {
	if m == nil {
		return
	}

	m.retries.Inc()
}

func newQueueMetrics(name string, clock clock.Clock) queueMetrics {
	return &defaultQueueMetrics{
		clock:                   clock,
		depth:                   depth.WithLabelValues(name),
		adds:                    adds.WithLabelValues(name),
		latency:                 latency.WithLabelValues(name),
		workDuration:            workDuration.WithLabelValues(name),
		unfinishedWorkSeconds:   unfinished.WithLabelValues(name),
		longestRunningProcessor: longestRunning.WithLabelValues(name),
		addTimes:                map[t]time.Time{},
		processingStartTimes:    map[t]time.Time{},
	}
}

func newRetryMetrics(name string) retryMetrics {
	var ret *defaultRetryMetrics
	if len(name) == 0 {
		return ret
	}
	return &defaultRetryMetrics{
		retries: retries.WithLabelValues(name),
	}
}
