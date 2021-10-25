package lru

import (
	"github.com/pkg/errors"
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

type MetricsConfig struct {
	Enabled   bool              `mapstructure:"enabled" json:"enabled" yaml:"enabled"`
	Namespace string            `mapstructure:"namespace" json:"namespace" yaml:"namespace"`
	Subsystem string            `mapstructure:"subsystem" json:"subsystem" yaml:"subsystem"`
	Labels    map[string]string `mapstructure:"labels" json:"labels" yaml:"labels"`
}

type ClockConfig struct {
	Precise  *ClockConfigPrecise  `mapstructure:"precise" json:"precise" yaml:"precise"`
	Discrete *ClockConfigDiscrete `mapstructure:"discrete" json:"discrete" yaml:"discrete"`
}

type ClockConfigPrecise struct{}
type ClockConfigDiscrete struct {
	UpdateInterval time.Duration `mapstructure:"update_interval" json:"update_interval" yaml:"update_interval"`
}

type Config struct {
	Capacity   int            `mapstructure:"capacity" json:"capacity" yaml:"capacity"`
	TTL        time.Duration  `mapstructure:"ttl" json:"ttl" yaml:"ttl"`
	Concurrent bool           `mapstructure:"concurrent" json:"concurrent" yaml:"concurrent"`
	Metrics    *MetricsConfig `mapstructure:"metrics" json:"metrics" yaml:"metrics"`
	Clock      *ClockConfig   `mapstructure:"clock" json:"clock" yaml:"clock"`
}

func (c *Config) Validate() error {
	// we have mandatory non-default parameters like size and ttl,
	// so empty config is not okay here
	if c == nil {
		return errors.New("empty config")
	}

	if c.Capacity <= 0 {
		return errors.New("capacity must be greater than zero")
	}

	if err := c.Metrics.Validate(); err != nil {
		return err
	}

	if err := c.Clock.Validate(); err != nil {
		return err
	}

	return nil
}

func (c *MetricsConfig) Validate() error {
	// empty config is okay
	if c == nil {
		return nil
	}

	if c.Enabled == false {
		return nil
	}

	if c.Namespace == "" {
		return errors.New("metrics namespace is empty")
	}

	if c.Subsystem == "" {
		return errors.New("metrics subsystem is empty")
	}

	return nil
}

func (c *ClockConfig) Validate() error {
	// empty config is okay
	if c == nil {
		return nil
	}

	clockConfigsFound := 0

	if c.Precise != nil {
		clockConfigsFound++
		if err := c.Precise.Validate(); err != nil {
			return errors.Wrap(err, "precise")
		}
	}

	if c.Discrete != nil {
		clockConfigsFound++
		if err := c.Discrete.Validate(); err != nil {
			return errors.Wrap(err, "discrete")
		}
	}

	if clockConfigsFound != 1 {
		return errors.New("exactly one clock config expected")
	}

	return nil
}

func (c *ClockConfigPrecise) Validate() error {
	return nil
}

func (c *ClockConfigDiscrete) Validate() error {
	return nil
}

func (c *Config) withDefaults() *Config {
	var ret = *c
	if ret.Metrics == nil {
		ret.Metrics = &MetricsConfig{
			Enabled: false,
		}
	}

	if ret.Clock == nil {
		ret.Clock = &ClockConfig{
			Precise: &ClockConfigPrecise{},
		}
	}

	return &ret
}
