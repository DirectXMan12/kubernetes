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

import (
	"sync"
)

// registryPromise is a promise of an eventual registry.  All registrations
// in this registry are delayed until the promise is resolved.
type registryPromise struct {
	metrics []Metric
}

func (p *registryPromise) MustRegisterGauge(metric *Gauge) {
	p.metrics = append(p.metrics, metric)
}
func (p *registryPromise) MustRegisterGaugeVec(metric *GaugeVec) {
	p.metrics = append(p.metrics, metric)
}
func (p *registryPromise) MustRegisterHistogram(metric *Histogram) {
	p.metrics = append(p.metrics, metric)
}
func (p *registryPromise) MustRegisterHistogramVec(metric *HistogramVec) {
	p.metrics = append(p.metrics, metric)
}
func (p *registryPromise) MustRegisterCounter(metric *Counter) {
	p.metrics = append(p.metrics, metric)
}
func (p *registryPromise) MustRegisterCounterVec(metric *CounterVec) {
	p.metrics = append(p.metrics, metric)
}

var (
	regPromise = &registryPromise{}
	currentRegistry Registry = regPromise
	registryMu sync.RWMutex
)

// DefaultRegistry returns the current default registry.  It will be a
// "promise" until ResolveRegistry is called.
func DefaultRegistry() Registry {
	registryMu.RLock()
	defer registryMu.RUnlock()
	return currentRegistry
}

// ResolveRegistry turns the default registry promise into an actual registry,
// registering all the metrics registered on the promise, and swapping out
// the default registry for future registrations.
func ResolveRegistry(actual Registry) {
	registryMu.RLock()
	defer registryMu.RUnlock()
	if regPromise == nil {
		panic("unable to resolve already-resolved metrics registry promise")
	}

	for _, metric := range regPromise.metrics {
		metric.MustRegisterIn(actual)
	}

	regPromise = nil
	currentRegistry = actual
}
