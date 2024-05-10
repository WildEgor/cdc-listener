package monitor

import (
	"github.com/WildEgor/cdc-listener/internal/configs"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var _ IMonitor = (*PromMetrics)(nil)

const (
	labelApp  = "app"
	labelDb   = "db"
	labelColl = "coll"
	labelKind = "kind"
)

const (
	ProblemKindParse   = "parse"
	ProblemKindDecode  = "decode"
	ProblemKindPublish = "publish"
	ProblemWatch       = "watch"
)

type PromMetrics struct {
	appConfig                                               *configs.AppConfig
	metricsConfig                                           *configs.MetricsConfig
	filterSkippedEvents, publishedEvents, problematicEvents *prometheus.CounterVec
	regAdapter                                              *PromMetricsRegistry
}

func NewPromMetrics(reg *PromMetricsRegistry, appConfig *configs.AppConfig, metricsConfig *configs.MetricsConfig) *PromMetrics {
	return &PromMetrics{
		appConfig:     appConfig,
		metricsConfig: metricsConfig,
		regAdapter:    reg,
		publishedEvents: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "published_events_total",
			Help: "The total number of published events",
		},
			[]string{labelApp, labelDb, labelColl},
		),
		problematicEvents: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "problematic_events_total",
			Help: "The total number of skipped problematic events",
		},
			[]string{labelApp, labelKind},
		),
		filterSkippedEvents: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "filter_skipped_events_total",
			Help: "The total number of skipped events",
		},
			[]string{labelApp, labelDb, labelColl},
		),
	}
}

func (m *PromMetrics) IncPublishedEventsCounter(db, coll string) {
	if m.metricsConfig.IsEnabled() {
		m.publishedEvents.With(prometheus.Labels{labelApp: m.appConfig.Name, labelDb: db, labelColl: coll}).Inc()
	}
}

func (m *PromMetrics) IncProblematicEventsCounter(kind string) {
	if m.metricsConfig.IsEnabled() {
		m.problematicEvents.With(prometheus.Labels{labelApp: m.appConfig.Name, labelKind: kind}).Inc()
	}
}

func (m *PromMetrics) IncFilteredEventsCounter(db, coll string) {
	if m.metricsConfig.IsEnabled() {
		m.filterSkippedEvents.With(prometheus.Labels{labelApp: m.appConfig.Name, labelDb: db, labelColl: coll}).Inc()
	}
}
