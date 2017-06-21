package config

import (
	"fmt"
)

// InfluxdbInfo stores the hostname and port number
type InfluxdbInfo struct {
	Hostname    string `json:"Hostname"`
	Port        int    `json:"Port"`
	Database    string `json:"Database"`
	Measurement string `json:"Measurement"`
	Username    string `json:"Username"`
	Password    string `json:"Password"`

	// Tags that we will extract from the log fields and set them as Influx point tags.
	Precision      string // 精确度
	Tags           []string
	BatchInterval  int // seconds
	BatchSize      int
	MaxRetryCount  int
	UseUDP         bool
	RetryCount     int
	ChTearDown     chan bool // just for testing, don't use it in production env
	ChTearDownDone chan bool // just for testing, don't use it in production env
}

const (
	// PrecisionDefault represents the default precision used for the InfluxDB points.
	PrecisionDefault = "ns"
	// DatabaseDefault is the default database that we will write to, if not specified otherwise in the Config for the hook.
	DatabaseDefault = "Devices"
	// DefaultMeasurement is the default measurement that we will assign to each point, unless there is a field called "measurement".
	DefaultMeasurement = "Measurement"
	// BatchIntervalDefault represents the number of seconds that we wait for a batch to fill up.
	// After that we flush it to InfluxDB whatsoever.
	BatchIntervalDefault = 5
	// BatchSizeDefault represents the maximum size of a batch.
	BatchSizeDefault = 1000
	// MaxRetryCountDefault represents the maximum number of retrying to connect DB.
	MaxRetryCountDefault = 10
)

// Address 补全 URL
func (config InfluxdbInfo) Address() string {
	return fmt.Sprintf("http://%s:%d", config.Hostname, config.Port)
}

func (config *InfluxdbInfo) SetDefaults() {
	if config.Precision == "" {
		config.Precision = PrecisionDefault
	}
	if config.Database == "" {
		config.Database = DatabaseDefault
	}
	if config.Measurement == "" {
		config.Measurement = DefaultMeasurement
	}
	if config.BatchInterval <= 0 {
		config.BatchInterval = BatchIntervalDefault
	}
	if config.BatchSize <= 0 {
		config.BatchSize = BatchSizeDefault
	}

	if config.MaxRetryCount <= 0 {
		config.MaxRetryCount = MaxRetryCountDefault
	}
}
