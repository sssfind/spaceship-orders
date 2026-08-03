package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Счётчик созданных заказов
	OrdersTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "orders_total",
		Help: "Total number of created orders",
	})

	// Счётчик выручки
	OrdersRevenueTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "orders_revenue_total",
		Help: "Total revenue from created orders",
	})
)
