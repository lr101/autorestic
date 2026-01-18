package monitoring

import (
	"fmt"
	"github.com/cupcakearmy/autorestic/internal/metadata"
)


type MonitorType string

const (
	MonitorTypeInflux MonitorType = "influx"
)

type Monitor struct {
	Type MonitorType       `mapstructure:"type,omitempty" yaml:"type,omitempty"`
	Env  map[string]string `mapstructure:"env,omitempty" yaml:"env,omitempty"`
}

type Reporter interface {
	Report(md *metadata.BackupLogMetadata, locationName string, tag string, backendName string) error
	Close()
}

// NewReporter is the "Factory". It takes a configuration struct and returns the correct Reporter.
func NewReporter(conf Monitor) (Reporter, error) {
	switch conf.Type {
	case MonitorTypeInflux:
		reporter := NewInfluxReporter(conf)
		if reporter == nil {
			return nil, fmt.Errorf("influx configuration is missing required fields (URL, Token, Org, or Bucket)")
		}
		return reporter, nil

	default:
		return nil, fmt.Errorf("unknown monitor type: '%s'", conf.Type)
	}
}

