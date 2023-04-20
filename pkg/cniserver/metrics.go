package cniserver

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	PortCountFunc func() float64
	registry      *prometheus.Registry

	cniRequestCount        prometheus.Counter
	cniRequestInvalidCount prometheus.Counter
	cniAddSuccessCount     prometheus.Counter
	cniAddFailureCount     prometheus.Counter
	cniDelSuccessCount     prometheus.Counter
	cniDelFailureCount     prometheus.Counter
	cniCheckSuccessCount   prometheus.Counter
	cniCheckFailureCount   prometheus.Counter
	reapSuccessCount       prometheus.Counter
	reapFailureCount       prometheus.Counter
	portTotal              prometheus.GaugeFunc
}

func NewMetrics(registry *prometheus.Registry, portCountFunc func() float64) *Metrics {
	// Request
	metrics := &Metrics{registry: registry, PortCountFunc: portCountFunc}
	metrics.cniRequestCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cni_request_count",
			Help: "total count of /cni requests",
		},
	)
	metrics.registry.MustRegister(metrics.cniRequestCount)
	metrics.cniRequestInvalidCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cni_request_invalid_count",
			Help: "total count of invalid /cni requests",
		},
	)
	metrics.registry.MustRegister(metrics.cniRequestInvalidCount)

	// ADD
	metrics.cniAddSuccessCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cni_add_success_count",
			Help: "total count of successfully CNI ADD commands",
		},
	)
	metrics.registry.MustRegister(metrics.cniAddSuccessCount)

	metrics.cniAddFailureCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cni_add_failure_count",
			Help: "total count of failed CNI ADD commands",
		},
	)
	metrics.registry.MustRegister(metrics.cniAddFailureCount)

	// DEL
	metrics.cniDelSuccessCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cni_del_success_count",
			Help: "total count of successfully CNI DEL commands",
		},
	)
	metrics.registry.MustRegister(metrics.cniDelSuccessCount)

	metrics.cniDelFailureCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cni_del_failure_count",
			Help: "total count of failed CNI DEL commands",
		},
	)
	metrics.registry.MustRegister(metrics.cniDelFailureCount)

	// CHECK
	metrics.cniCheckSuccessCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cni_check_success_count",
			Help: "total count of successfully CNI CHECK commands",
		},
	)
	metrics.registry.MustRegister(metrics.cniCheckSuccessCount)

	metrics.cniCheckFailureCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cni_check_failure_count",
			Help: "total count of failed CNI CHECK commands",
		},
	)
	metrics.registry.MustRegister(metrics.cniCheckFailureCount)

	// Reaper
	metrics.reapSuccessCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cni_reaped_port_success_count",
			Help: "total count of successfully reaped ports",
		},
	)
	metrics.registry.MustRegister(metrics.reapSuccessCount)
	metrics.reapFailureCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cni_reaped_port_failure_count",
			Help: "total count of failed port reapings",
		},
	)
	metrics.registry.MustRegister(metrics.reapFailureCount)

	// Ports
	metrics.portTotal = prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "cni_port_total",
			Help: "number of active ports",
		},
		metrics.PortCountFunc,
	)
	metrics.registry.MustRegister(metrics.portTotal)

	return metrics
}

func (me *Metrics) Registry() *prometheus.Registry {
	return me.registry
}
