package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

type PromMetricsRegistry struct {
	Reg *prometheus.Registry
}

func NewPromMetricsRegistry() *PromMetricsRegistry {
	reg := prometheus.NewRegistry()
	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	return &PromMetricsRegistry{
		Reg: reg,
	}
}
