package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RegistrationsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "iam_registrations_total",
		Help: "Total number of successful user registrations",
	})

	LoginsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "iam_logins_total",
		Help: "Total number of login attempts",
	}, []string{"result"})

	LogoutsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "iam_logouts_total",
		Help: "Total number of successful logouts",
	})
)
