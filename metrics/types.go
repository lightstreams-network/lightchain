package metrics

import (
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/go-kit/kit/metrics/discard"
	"github.com/go-kit/kit/metrics"
)

type Counter metrics.Counter
type Gauge metrics.Gauge

type CounterSet struct {
	itLvl []string
	c     Counter
}

type GaugeSet struct {
	itLvl []string
	g     Gauge
}

func (cs *CounterSet) Add(delta float64, labelValues ...string) {
	var cc = cs.withLabelValues(labelValues...)
	cc.(Counter).Add(delta)
}

func (cs *CounterSet) withLabelValues(labelValues ...string) Counter {
	cc := cs.c
	for k, v := range cs.itLvl {
		if len(labelValues) > k {
			cc = cc.With(v, labelValues[k])
		} else {
			cc = cc.With(v, "")
		}
	}

	return cc
}

func (gs *GaugeSet) Add(delta float64, labelValues ...string) {
	var cc = gs.withLabelValues(labelValues...)
	cc.(Gauge).Add(delta)
}

func (gs *GaugeSet) Set(delta float64, labelValues ...string) {
	var cc = gs.withLabelValues(labelValues...)
	cc.(Gauge).Set(delta)
}

func (cs *GaugeSet) withLabelValues(labelValues ...string) Gauge {
	cc := cs.g
	for k, v := range cs.itLvl {
		if len(labelValues) > k {
			cc = cc.With(v, labelValues[k])
		} else {
			cc = cc.With(v, "")
		}
	}

	return cc
}


func NewGaugeMetric(registry *prometheus.Registry, opts prometheus.GaugeOpts, labelsAndValues ...string) Gauge {
	labels := []string{}
	for i := 0; i < len(labelsAndValues); i += 2 {
		labels = append(labels, labelsAndValues[i])
	}

	collection := prometheus.NewGaugeVec(opts, labels)
	registry.MustRegister(collection)
	return kitprometheus.NewGauge(collection).With(labelsAndValues...)
}

func NewCounterMetric(registry *prometheus.Registry, opts prometheus.CounterOpts, labelsAndValues ...string) Counter {
	labels := []string{}
	for i := 0; i < len(labelsAndValues); i += 2 {
		labels = append(labels, labelsAndValues[i])
	}

	collection := prometheus.NewCounterVec(opts, labels)
	registry.Register(collection)
	return kitprometheus.NewCounter(collection).With(labelsAndValues...)
}

func NewGaugeSetMetric(registry *prometheus.Registry,
	opts prometheus.GaugeOpts,
	itLvl []string,
	labelsAndValues ...string,
) GaugeSet {
	labels := []string{}
	for i := 0; i < len(labelsAndValues); i += 2 {
		labels = append(labels, labelsAndValues[i])
	}

	collection := prometheus.NewGaugeVec(opts, append(labels, itLvl...))
	registry.Register(collection)

	return GaugeSet{
		itLvl: itLvl,
		g: kitprometheus.NewGauge(collection).With(labelsAndValues...),
	}
}

func NewCounterSetMetric(registry *prometheus.Registry,
	opts prometheus.CounterOpts,
	itLvl []string,
	labelsAndValues ...string,
) CounterSet {
	labels := []string{}
	for i := 0; i < len(labelsAndValues); i += 2 {
		labels = append(labels, labelsAndValues[i])
	}

	collection := prometheus.NewCounterVec(opts, append(labels, itLvl...))
	registry.Register(collection)
	return CounterSet{
		itLvl: itLvl,
		c: kitprometheus.NewCounter(collection).With(labelsAndValues...),
	}
}

func NewGaugeDiscard() Gauge{
	return discard.NewGauge()
}

func NewCounterDiscard() Counter{
	return discard.NewCounter()
}

func NewCounterSetDiscard() CounterSet{
	return CounterSet{
		itLvl: []string{},
		c: discard.NewCounter(),
	}
}

func NewGaugeSetDiscard() GaugeSet{
	return GaugeSet{
		itLvl: []string{},
		g: discard.NewGauge(),
	}
}
