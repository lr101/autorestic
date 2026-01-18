package monitoring

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cupcakearmy/autorestic/internal/metadata"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

// InfluxReporter handles connections and writes to InfluxDB
type InfluxReporter struct {
	client   influxdb2.Client
	writeAPI api.WriteAPIBlocking
}

func NewInfluxReporter(name string, monitor Monitor) (Reporter, error) {
	// Generate the merged environment map
	env := getEnv(name, monitor)

	serverURL := env["INFLUX_URL"]
	token := env["INFLUX_TOKEN"]
	org := env["INFLUX_ORG"]
	bucket := env["INFLUX_BUCKET"]


	// Return nil if configuration is incomplete
	if serverURL == "" || token == "" || org == "" || bucket == "" {
		return nil, fmt.Errorf("influx configuration is missing required fields (URL, Token, Org, or Bucket)")
	}

	client := influxdb2.NewClient(serverURL, token)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
	if _, err := client.Health(ctx); err != nil {
        client.Close() // Clean up
        return nil, fmt.Errorf("InfluxDB connection failed: %v", err)
    }
    bucketAPI := client.BucketsAPI()
    if _, err := bucketAPI.FindBucketByName(ctx, bucket); err != nil {
        client.Close()
        return nil, fmt.Errorf("InfluxDB bucket '%s' validation failed: %v", bucket, err)
    }

	writeAPI := client.WriteAPIBlocking(org, bucket)

	return &InfluxReporter{
		client:   client,
		writeAPI: writeAPI,
	}, nil
}

// Close ensures the client connection is closed gracefully
func (ir *InfluxReporter) Close() {
	if ir.client != nil {
		ir.client.Close()
	}
}


func (ir *InfluxReporter) Report(md *metadata.BackupLogMetadata, locationName string, tag string, backendName string) error {
	// 1. Prepare Tags
	tags := map[string]string{
		"location":   locationName,
		"backend":    backendName,
		"exit_code":  md.ExitCode,
		"snapshot_id": md.SnapshotID,
		"tag": 	  tag,
	}

	// 2. Prepare Fields
	fields := map[string]interface{}{
		"files_added":          parseInt(md.Files.Added),
		"files_changed":        parseInt(md.Files.Changed),
		"files_unmodified":     parseInt(md.Files.Unmodified),
		"dirs_added":           parseInt(md.Dirs.Added),
		"dirs_changed":         parseInt(md.Dirs.Changed),
		"dirs_unmodified":      parseInt(md.Dirs.Unmodified),
		"added_size_bytes":     parseBytes(md.AddedSize),
		"processed_files":      parseInt(md.Processed.Files),
		"processed_size_bytes": parseBytes(md.Processed.Size),
		"duration_seconds":     parseDuration(md.Processed.Duration),
	}

	// 3. Create Point
	p := influxdb2.NewPoint(
		"autorestic_backup",
		tags,
		fields,
		time.Now(),
	)

	return ir.writeAPI.WritePoint(context.Background(), p)
}

func parseInt(s string) int {
	clean := strings.ReplaceAll(s, ",", "")
	val, err := strconv.Atoi(clean)
	if err != nil {
		return 0
	}
	return val
}

func parseBytes(s string) int64 {
	re := regexp.MustCompile(`(?i)^([\d\.]+)\s*([kMGTP]?i?B?)$`)
	matches := re.FindStringSubmatch(strings.TrimSpace(s))
	
	if len(matches) < 3 {
		return 0
	}

	val, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0
	}

	unit := strings.ToUpper(matches[2])
	var multiplier float64 = 1

	switch unit {
	case "KB", "KIB": multiplier = 1024
	case "MB", "MIB": multiplier = 1024 * 1024
	case "GB", "GIB": multiplier = 1024 * 1024 * 1024
	case "TB", "TIB": multiplier = 1024 * 1024 * 1024 * 1024
	}

	return int64(val * multiplier)
}

func parseDuration(s string) float64 {
	if strings.Contains(s, ":") {
		parts := strings.Split(s, ":")
		if len(parts) == 2 {
			min, _ := strconv.ParseFloat(parts[0], 64)
			sec, _ := strconv.ParseFloat(parts[1], 64)
			return (min * 60) + sec
		}
	}

	d, err := time.ParseDuration(s)
	if err != nil {
		return 0
	}
	return d.Seconds()
}