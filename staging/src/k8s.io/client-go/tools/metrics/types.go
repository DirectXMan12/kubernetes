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

// GaugeImpl knows how to manage the actual mechanics of a gauge metric,
// which can be incremented, decremented, or set to an arbitrary value.
type GaugeImpl interface {
	// Inc increments the gauge by 1.0.
	Inc()
	// Dec decrements the gauge by 1.0.
	Dec()
	// Set sets the gauge to an arbitrary value.
	Set(float64)
}

// Gauge is a kind of metric that varies over time, and can be set to arbitrary values.
type Gauge struct {
	// Name is the name of the metric.  Metrics with the same name must have
	// the same set of label names.  It must be a valid Prometheus metric name.
	Name string
	// Help is the help string for this metric.
	Help string

	GaugeImpl
}

func (g *Gauge) MustRegisterIn(reg Registry) {
	reg.MustRegisterGauge(g)
}

// NewGauge creates a new guage with the given name.
func NewGauge(base Gauge) *Gauge {
	base.GaugeImpl = noopMetric{}
	return &base
}

// HistogramImpl knows how to manage the actual mechanics of a histogram metric,
// which can be observed and aggregated
type HistogramImpl interface {
	// Observe notes that a particular value occurred.
	Observe(float64)
}

// Histogram is a kind of metric that records an aggregation of observed values.
type Histogram struct {
	// Name is the name of the metric.  Metrics with the same name must have
	// the same set of label names.  It must be a valid Prometheus metric name.
	Name string
	// Help is the help string for this metric.
	Help string

	Buckets []float64
	HistogramImpl
}

func (g *Histogram) MustRegisterIn(reg Registry) {
	reg.MustRegisterHistogram(g)
}

// NewHistogram creates a new histogram with the given name.
func NewHistogram(base Histogram) *Histogram {
	base.HistogramImpl = noopMetric{}
	return &base
}

// CounterImpl knows how to manage the actual mechanics of a counter metric,
// which can be incremented over time
type CounterImpl interface {
	// Inc increments the metric by 1.
	Inc()
}

// Counter is a kind of metric that only ever increments.
type Counter struct {
	// Name is the name of the metric.  Metrics with the same name must have
	// the same set of label names.  It must be a valid Prometheus metric name.
	Name string
	// Help is the help string for this metric.
	Help string

	CounterImpl
}

func (g *Counter) MustRegisterIn(reg Registry) {
	reg.MustRegisterCounter(g)
}

// NewCounter creates a new counter with the given name.
func NewCounter(base Counter) *Counter {
	base.CounterImpl = noopMetric{}
	return &base
}

// GaugeVecImpl knows how to manage the actual mechanics of a series of gauges with common label names.
type GaugeVecImpl interface {
	// WithLabelValues returns a gauge with labels named by this vector set to the given values.
	WithLabelValues(...string) GaugeImpl
}

// GaugeVec is a series of gauges with the same label names but different values.
type GaugeVec struct {
	Gauge
	// LabelNames is the names of all common variable labels for this vector.
	LabelNames []string

	GaugeVecImpl
}

func (g *GaugeVec) MustRegisterIn(reg Registry) {
	reg.MustRegisterGaugeVec(g)
}

// NewGaugeVec creates a new guage series with the given name and label names.
func NewGaugeVec(base Gauge, labelNames ...string) *GaugeVec {
	return &GaugeVec{
		Gauge: base,
		LabelNames: labelNames,
		GaugeVecImpl: noopGaugeVec{},
	}
}

type noopGaugeVec struct{}
func (n noopGaugeVec) WithLabelValues(_ ...string) GaugeImpl { return noopMetric{} }

// HistogramVecImpl knows how to manage the actual mechanics of a series of histograms with common label names.
type HistogramVecImpl interface {
	// WithLabelValues returns a histogram with labels named by this vector set to the given values.
	WithLabelValues(...string) HistogramImpl
}

// HistogramVec is a series of histograms with the same label names but different values.
type HistogramVec struct {
	Histogram
	// LabelNames is the names of all common variable labels for this vector.
	LabelNames []string

	HistogramVecImpl
}

func (g *HistogramVec) MustRegisterIn(reg Registry) {
	reg.MustRegisterHistogramVec(g)
}

// NewHistogramVec creates a new guage series with the given name and label names.
func NewHistogramVec(base Histogram, labelNames ...string) *HistogramVec {
	return &HistogramVec{
		Histogram: base,
		LabelNames: labelNames,
		HistogramVecImpl: noopHistogramVec{},
	}
}

type noopHistogramVec struct{}
func (n noopHistogramVec) WithLabelValues(_ ...string) HistogramImpl { return noopMetric{} }

// CounterVecImpl knows how to manage the actual mechanics of a series of counters with common label names.
type CounterVecImpl interface {
	// WithLabelValues returns a counter with labels named by this vector set to the given values.
	WithLabelValues(...string) CounterImpl
}

// CounterVec is a series of counters with the same label names but different values.
type CounterVec struct {
	Counter
	// LabelNames is the names of all common variable labels for this vector.
	LabelNames []string

	CounterVecImpl
}

func (g *CounterVec) MustRegisterIn(reg Registry) {
	reg.MustRegisterCounterVec(g)
}

// NewCounterVec creates a new guage series with the given name and label names.
func NewCounterVec(base Counter, labelNames ...string) *CounterVec {
	return &CounterVec{
		Counter: base,
		LabelNames: labelNames,
		CounterVecImpl: noopCounterVec{},
	}
}

type noopCounterVec struct{}
func (n noopCounterVec) WithLabelValues(_ ...string) CounterImpl { return noopMetric{} }
