package app

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/Djarvur/allcups-itrally-2020-task/pkg/def"
)

//nolint:gochecknoglobals // Metrics are global anyway.
var (
	Metric def.Metrics // Common metrics used by all packages.
	metric struct {
		autosaveDuration prometheus.Histogram
	}
)

// InitMetrics must be called once before using this package.
// It registers and initializes metrics used by this package.
func InitMetrics(reg *prometheus.Registry) {
	const subsystem = "app"

	Metric = def.NewMetrics(reg)
	metric.autosaveDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Subsystem: subsystem,
			Name:      "autosave_duration_seconds",
			Help:      "Autosave latency distributions.",
		},
	)
	reg.MustRegister(metric.autosaveDuration)
}
