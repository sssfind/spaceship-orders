package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Гистограмма длительности сборки
	AssemblyDurationSeconds = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "assembly_duration_seconds",
		Help:    "Duration of ship assembly processing in seconds",
		Buckets: prometheus.DefBuckets,
	})
)
