package monitoring

import (
	"fmt"
	"os"
	"strings"
	"regexp"
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

var nonAlphaRegex = regexp.MustCompile(`[^a-zA-Z0-9]+`)

func getEnv(name string, m Monitor) map[string]string {
	env := make(map[string]string)

	// 1. Load from YAML config (convert keys to UPPERCASE)
	for key, value := range m.Env {
		env[strings.ToUpper(key)] = value
	}

	// 2. Load from System Env / .autorestic.env
	// Format: AUTORESTIC_<NAME>_<KEY>
	nameForEnv := strings.ToUpper(name)
	nameForEnv = nonAlphaRegex.ReplaceAllString(nameForEnv, "_")
	prefix := "AUTORESTIC_" + nameForEnv + "_"

	for _, variable := range os.Environ() {
		// Split "KEY=VALUE"
		parts := strings.SplitN(variable, "=", 2)
		if len(parts) != 2 {
			continue
		}

		// Check if this env var belongs to this monitor
		if strings.HasPrefix(parts[0], prefix) {
			// Remove prefix to get the config key (e.g. "INFLUX_URL")
			key := strings.TrimPrefix(parts[0], prefix)
			env[key] = parts[1]
		}
	}
	return env
}

// NewReporter is the "Factory". It takes a configuration struct and returns the correct Reporter.
func NewReporter(name string, conf Monitor) (Reporter, error) {
	switch conf.Type {
	case MonitorTypeInflux:
		reporter, err := NewInfluxReporter(name, conf)
		if err != nil {
			return nil, err
		}
		return reporter, nil

	default:
		return nil, fmt.Errorf("unknown monitor type: '%s'", conf.Type)
	}
}

