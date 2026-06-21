package buffer

import (
    "github.com/prometheus/client_golang/prometheus"
)

var (
    fillRatioGauge = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "telemetry_buffer_fill_ratio",
        Help: "Current fill ratio of the telemetry buffer",
    })
    eventsDroppedCounter = prometheus.NewCounter(prometheus.CounterOpts{
        Name: "telemetry_events_dropped_total",
        Help: "Total number of telemetry events dropped due to buffer full",
    })
    publishLatencyHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
        Name:    "telemetry_publish_latency_seconds",
        Help:    "Latency of publishing telemetry batches",
        Buckets: prometheus.DefBuckets,
    })
    publishErrorsCounter = prometheus.NewCounter(prometheus.CounterOpts{
        Name: "telemetry_publish_errors_total",
        Help: "Total number of publish errors",
    })
    enqueuedTotalCounter = prometheus.NewCounter(prometheus.CounterOpts{
        Name: "telemetry_events_enqueued_total",
        Help: "Total number of telemetry events enqueued",
    })
    consumerAliveGauge = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "telemetry_consumer_alive",
        Help: "Gauge indicating if the consumer goroutine is running (1) or stopped (0)",
    })
)

func init() {
    prometheus.MustRegister(
        fillRatioGauge,
        eventsDroppedCounter,
        publishLatencyHistogram,
        publishErrorsCounter,
        enqueuedTotalCounter,
        consumerAliveGauge,
    )
}
