package lru

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type optionCapacity struct{ capacity int }
type optionTTL struct{ ttl time.Duration }
type optionMetrics struct {
	namespace   string
	subsystem   string
	constLabels prometheus.Labels
}
type optionSync struct{}
type optionDiscreteClock struct{ updateInterval time.Duration }
type optionEvictCallback struct{ cb func(string) }
type optionExpireCallback struct{ cb func(string) }
