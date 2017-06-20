package config

import (
	"fmt"

	"log"
	"time"

	client "github.com/influxdata/influxdb/client/v2"
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
	retryCount     int
	chTearDown     chan bool // just for testing, don't use it in production env
	chTearDownDone chan bool // just for testing, don't use it in production env
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

// Activity represents the operation that specific module performs
type Activity struct {
	module  string
	status  string // 状态， 成功：success, 失败: fail
	action  string // 操作
	message string
}

// DataNode 数据点
type DataNode struct {
	Tags   map[string]string
	Fields map[string]interface{}
}

// Hook represents the hook to InfluxDB
type Hook struct {
	config       InfluxdbInfo
	chActivities chan *Activity
	chData       chan DataNode
}

// NewInfluxDBHook generate a new InfluxDB hook based on the given configuration
func NewInfluxDBHook() (*Hook, error) {
	hook := &Hook{
		config: Cfg.InfluxDB,
	}

	hook.config.setDefaults()

	// Make a buffered channel so that senders will not block.
	hook.chActivities = make(chan *Activity, hook.config.BatchSize)
	hook.chData = make(chan DataNode, hook.config.BatchSize)
	hook.config.chTearDown = make(chan bool)
	hook.config.chTearDownDone = make(chan bool)
	go hook.startBatchHandler()

	return hook, nil
}

// Address 补全 URL
func (config InfluxdbInfo) Address() string {
	return fmt.Sprintf("http://%s:%d", config.Hostname, config.Port)
}

func (config *InfluxdbInfo) setDefaults() {
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

func (hook *Hook) startBatchHandler() {
	done := false
	// Make client
	for ; done || hook.config.retryCount < hook.config.MaxRetryCount; hook.config.retryCount++ {
		// wait for some seconds and retry
		time.Sleep(time.Duration(hook.config.retryCount) * time.Second)

		var err error
		var c client.Client
		if hook.config.UseUDP {
			c, err = client.NewUDPClient(client.UDPConfig{Addr: hook.config.Address()})
		} else {
			c, err = client.NewHTTPClient(client.HTTPConfig{
				Addr:     hook.config.Address(),
				Username: hook.config.Username,
				Password: hook.config.Password,
			})
		}
		if err != nil {
			log.Printf("[logops] Make client #%d, Error: %v", hook.config.retryCount, err)
			continue
		}
		defer c.Close()

		if hook.config.UseUDP == false {
			q := client.NewQuery(fmt.Sprintf("CREATE DATABASE %s", hook.config.Database), "", "")
			if _, err := c.Query(q); err != nil {
				log.Print("[logops] Failed to create db ", hook.config.Database, err)
				continue
			}
		}

		for true {
			select {
			case pointData := <-hook.chData:
				// Create a new point batch
				bp, err := client.NewBatchPoints(client.BatchPointsConfig{
					Database:  hook.config.Database,
					Precision: hook.config.Precision,
				})
				if err != nil {
					log.Print("[logops] NewBatchPoints Error: ", err)
					break
				}

				// Create a point and add to batch
				pt, err := client.NewPoint(hook.config.Measurement, pointData.Tags, pointData.Fields)
				if err != nil {
					log.Print("[logops] NewPoint Error: ", err)
					break
				}
				bp.AddPoint(pt)

				// Write the batch
				err = c.Write(bp)
				if err != nil {
					log.Print("[logops] Write Error: ", err)
					break
				}

			case a := <-hook.chActivities:
				// Create a new point batch
				bp, err := client.NewBatchPoints(client.BatchPointsConfig{
					Database:  hook.config.Database,
					Precision: hook.config.Precision,
				})
				if err != nil {
					log.Fatal("[logops] NewBatchPoints Error: ", err)
					break
				}

				// Create a point and add to batch
				tags := map[string]string{"module": a.module}
				fields := map[string]interface{}{
					"value":   1,
					"status":  a.status,
					"action":  a.action,
					"message": a.message,
				}
				pt, err := client.NewPoint(hook.config.Measurement, tags, fields)
				if err != nil {
					log.Print("[logops] NewPoint Error: ", err)
					break
				}
				bp.AddPoint(pt)

				// Write the batch
				err = c.Write(bp)
				if err != nil {
					log.Print("[logops] Write Error: ", err)
					break
				}
			// For testing
			case <-hook.config.chTearDown:
				if hook.config.UseUDP == false {
					q := client.NewQuery(fmt.Sprintf("DROP DATABASE %s", hook.config.Database), "", "")
					if response, err := c.Query(q); err != nil {
						log.Print("Failed to create db ", hook.config.Database, response.Error())
					}
				}
				done = true
				hook.config.chTearDownDone <- done
				break
			}
		}
	}

	panic("Filed to start batch handler")
}

func (hook *Hook) Write(module, status, action, message string) {
	a := &Activity{
		module:  module,
		status:  status,
		action:  action,
		message: message,
	}
	hook.chActivities <- a
}

// WriteNode ...
func (hook *Hook) WriteNode(data DataNode) {
	hook.chData <- data
}

func (hook *Hook) tearDown() {
	hook.config.chTearDown <- true
	<-hook.config.chTearDownDone
}
