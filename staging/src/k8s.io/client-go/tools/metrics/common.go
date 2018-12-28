/*
Copyright 2018 The Kubernetes Authors.

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

// Package metrics provides abstractions for registering which metrics
// to record.
package metrics

// Metric represents some metric that can be registered in a registry.
// A given metric may only ever be registered in a single registry,
// and does nothing until registered.
type Metric interface {
	// MustRegisterIn tries to register this metric in the given registry,
	// and panics if it fails.
	MustRegisterIn(Registry)
}

// GeneralMetricDesc contains common descriptions across all metrics.
type GeneralMetricDesc struct {
	// Name is the name of the metric.  Metrics with the same name must have
	// the same set of label names.  It must be a valid Prometheus metric name.
	Name string
	// Help is the help string for this metric.
	Help string
}

// GeneralVecDesc contains common descriptions across all metric vectors.
type GeneralVecDesc struct {
	// LabelNames is the names of all common variable labels for this vector.
	LabelNames []string
}

// Registry knows how to actually implement and expose metrics.
type Registry interface {
	// MustRegisterGauge tries to register the given gauge, panicing if it cannot.
	MustRegisterGauge(*Gauge)
	// MustRegisterGaugeVec tries to register the given gauge series, panicing if it cannot.
	MustRegisterGaugeVec(*GaugeVec)
	// MustRegisterHistogram tries to register the given histogram, panicing if it cannot.
	MustRegisterHistogram(*Histogram)
	// MustRegisterHistogramVec tries to register the given histogram series, panicing if it cannot.
	MustRegisterHistogramVec(*HistogramVec)
	// MustRegisterCounter tries to register the given counter, panicing if it cannot.
	MustRegisterCounter(*Counter)
	// MustRegisterCounterVec tries to register the given counter series, panicing if it cannot.
	MustRegisterCounterVec(*CounterVec)
}

// noopMetric is a metric implementation that does nothing.
type noopMetric struct{}
func (m noopMetric) Inc() {}
func (m noopMetric) Dec() {}
func (m noopMetric) Set(_ float64) {}
func (m noopMetric) Observe(_ float64) {}

