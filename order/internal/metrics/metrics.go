package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	OrdersCreatedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "orders_created_total",
		Help: "Total number of created orders",
	})

	OrdersPaidTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "orders_paid_total",
		Help: "Total number of paid orders",
	}, []string{"payment_method"})

	OrdersCancelledTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "orders_cancelled_total",
		Help: "Total number of cancelled orders",
	})

	OrdersRevenueTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "orders_revenue_total",
		Help: "Total revenue from paid orders",
	}, []string{"payment_method"})

	OrderTotalPrice = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "order_total_price",
		Help:    "Distribution of order total price at creation",
		Buckets: prometheus.ExponentialBuckets(100, 2, 12),
	})

	OrdersCreateFailedTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "orders_create_failed_total",
		Help: "Total number of failed order create attempts",
	}, []string{"reason"})
)
