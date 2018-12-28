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

package prometheus

import (
	prom "github.com/prometheus/client_golang/prometheus"

	"k8s.io/client-go/tools/metrics"
)

// NB(directxman12): keep this in a subdirectory and don't depend on it from elsewhere
// in client-go, so that prometheus is an optional dependency.

// NewPrometheusRegistry returns a new metrics.Registry backed
// by the given prometheus.Registerer.
func NewPrometheusRegistry(reg prom.Registerer) metrics.Registry {
	return &prometheusRegistry{
		reg: reg,
	}
}

// ResolvreToDefaultPrometheus resolves the global registry promise to the default prometheus registerer.
func ResolveToDefaultPrometheus() {
	metrics.ResolveRegistry(NewPrometheusRegistry(prom.DefaultRegisterer))
}

type prometheusRegistry struct {
	reg prom.Registerer
}

func (p *prometheusRegistry) MustRegisterGauge(metric *metrics.Gauge) {
	actual := prom.NewGauge(prom.GaugeOpts{
		Name: metric.Name,
		Help: metric.Help,
	})
	p.reg.MustRegister(actual)
	metric.GaugeImpl = actual
}

func (p *prometheusRegistry) MustRegisterGaugeVec(metric *metrics.GaugeVec) {
	actual := prom.NewGaugeVec(prom.GaugeOpts{
		Name: metric.Name,
		Help: metric.Help,
	}, metric.LabelNames)
	p.reg.MustRegister(actual)
	metric.GaugeVecImpl = gaugeVec{actual}
}

func (p *prometheusRegistry) MustRegisterHistogram(metric *metrics.Histogram) {
	actual := prom.NewHistogram(prom.HistogramOpts{
		Name: metric.Name,
		Help: metric.Help,
		Buckets: metric.Buckets,
	})
	p.reg.MustRegister(actual)
	metric.HistogramImpl = actual
}

func (p *prometheusRegistry) MustRegisterHistogramVec(metric *metrics.HistogramVec) {
	actual := prom.NewHistogramVec(prom.HistogramOpts{
		Name: metric.Name,
		Help: metric.Help,
		Buckets: metric.Buckets,
	}, metric.LabelNames)
	p.reg.MustRegister(actual)
	metric.HistogramVecImpl = histogramVec{actual}
}

func (p *prometheusRegistry) MustRegisterCounter(metric *metrics.Counter) {
	actual := prom.NewCounter(prom.CounterOpts{
		Name: metric.Name,
		Help: metric.Help,
	})
	p.reg.MustRegister(actual)
	metric.CounterImpl = actual
}

func (p *prometheusRegistry) MustRegisterCounterVec(metric *metrics.CounterVec) {
	actual := prom.NewCounterVec(prom.CounterOpts{
		Name: metric.Name,
		Help: metric.Help,
	}, metric.LabelNames)
	p.reg.MustRegister(actual)
	metric.CounterVecImpl = counterVec{actual}
}

type gaugeVec struct {
	vec *prom.GaugeVec
}
func (v gaugeVec) WithLabelValues(vals ...string) metrics.GaugeImpl {
	return v.vec.WithLabelValues(vals...)
}

type histogramVec struct {
	vec *prom.HistogramVec
}
func (v histogramVec) WithLabelValues(vals ...string) metrics.HistogramImpl {
	return v.vec.WithLabelValues(vals...)
}

type counterVec struct {
	vec *prom.CounterVec
}
func (v counterVec) WithLabelValues(vals ...string) metrics.CounterImpl {
	return v.vec.WithLabelValues(vals...)
}
