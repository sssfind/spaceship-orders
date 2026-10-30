package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	OrdersCreatedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "orders_created_total",
		Help: "Total number of successfully created orders",
	})

	OrdersCreateFailedTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "orders_create_failed_total",
		Help: "Total number of failed order create attempts",
	}, []string{"reason"})

	OrdersPaidTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "orders_paid_total",
		Help: "Total number of successfully paid orders",
	}, []string{"payment_method"})

	OrdersPayFailedTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "orders_pay_failed_total",
		Help: "Total number of failed pay attempts",
	}, []string{"reason"})

	OrdersCancelledTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "orders_cancelled_total",
		Help: "Total number of cancelled orders",
	})

	OrdersRevenueTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "orders_revenue_total",
		Help: "Total revenue from created orders",
	})
)

func init() {
	// CounterVec не попадает в /metrics, пока не создан хотя бы один лейбл.
	for _, reason := range []string{"injected", "empty_parts", "not_found", "stock", "repo"} {
		OrdersCreateFailedTotal.WithLabelValues(reason)
	}
	for _, reason := range []string{"repo", "not_found", "status"} {
		OrdersPayFailedTotal.WithLabelValues(reason)
	}
	for _, method := range []string{"CARD", "SBP", "CREDIT_CARD", "INVESTOR_MONEY", "UNKNOWN"} {
		OrdersPaidTotal.WithLabelValues(method)
	}
}
