package system

import (
	"fmt"

	"log"
	"time"

	"net/url"

	cc "github.com/influxdata/influxdb/client"
	client "github.com/influxdata/influxdb/client/v2"
	"github.com/ifnfn/util/config"
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
	config       config.InfluxdbInfo
	chActivities chan *Activity
	chData       chan DataNode
}

// NewInfluxDBHook generate a new InfluxDB hook based on the given configuration
func NewInfluxDBHook() (*Hook, error) {
	hook := &Hook{
		config: config.InfluxDB,
	}

	hook.config.SetDefaults()

	// Make a buffered channel so that senders will not block.
	hook.chActivities = make(chan *Activity, hook.config.BatchSize)
	hook.chData = make(chan DataNode, hook.config.BatchSize)
	hook.config.ChTearDown = make(chan bool)
	hook.config.ChTearDownDone = make(chan bool)
	go hook.startBatchHandler()

	return hook, nil
}

func (hook *Hook) startBatchHandler() {
	done := false
	// Make client
	for ; done || hook.config.RetryCount < hook.config.MaxRetryCount; hook.config.RetryCount++ {
		// wait for some seconds and retry
		time.Sleep(time.Duration(hook.config.RetryCount) * time.Second)

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
			log.Printf("[logops] Make client #%d, Error: %v", hook.config.RetryCount, err)
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
			case <-hook.config.ChTearDown:
				if hook.config.UseUDP == false {
					q := client.NewQuery(fmt.Sprintf("DROP DATABASE %s", hook.config.Database), "", "")
					if response, err := c.Query(q); err != nil {
						log.Print("Failed to create db ", hook.config.Database, response.Error())
					}
				}
				done = true
				hook.config.ChTearDownDone <- done
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
	hook.config.ChTearDown <- true
	<-hook.config.ChTearDownDone
}

// test read influxdb data
func (hook *Hook) QueryNode() {
	host, err := url.Parse(fmt.Sprintf("http://%s:%d", hook.config.Hostname, 8086))
	if err != nil {
		log.Fatal(err)
	}
	con, err := cc.NewClient(cc.Config{URL: *host})
	if err != nil {
		log.Fatal(err)
	}
	q := cc.Query{
		Command:  fmt.Sprintf("select * from %s limit 1000", hook.config.Measurement),
		Database: hook.config.Database,
	}
	if response, err := con.Query(q); err == nil && response.Error() == nil {
		log.Println(response.Results)
	}
}
