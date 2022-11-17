package observer

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

const DefaultNamespace = "timetrack"

type Observer interface {
	Count(ctx string, category string)
	DurationOf(ctx, category string, startedAt time.Time)
}

type observerImpl struct {
	counter  *prometheus.CounterVec
	duration *prometheus.SummaryVec
}

func New(subSystem string) Observer {
	var (
		counter = prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: DefaultNamespace,
			Subsystem: subSystem,
			Name:      "counter",
			Help:      "Counter by context",
		}, []string{"context", "category"})

		duration = prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Namespace: DefaultNamespace,
			Subsystem: subSystem,
			Name:      "duration",
			Help:      "Duration by context",
			Objectives: map[float64]float64{
				0.5:  0.05,
				0.9:  0.01,
				0.99: 0.001,
			},
		}, []string{"context", "category"})
	)

	prometheus.MustRegister(counter)
	prometheus.MustRegister(duration)

	return &observerImpl{
		counter:  counter,
		duration: duration,
	}
}

func (o *observerImpl) Count(ctx, category string) {
	go func() {
		o.counter.WithLabelValues(ctx, category).Inc()
	}()
}

func (o *observerImpl) DurationOf(ctx, category string, startedAt time.Time) {
	go func() {
		o.duration.
			WithLabelValues(ctx, fmt.Sprint(category)).
			Observe(float64(time.Since(startedAt).Nanoseconds()) / 1e6)
	}()
}
