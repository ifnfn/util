package config

import (
	"fmt"

	mgo "gopkg.in/mgo.v2"
)

// MgoSession ...
type MgoSession *mgo.Session

// MongodbInfo ...
type MongodbInfo struct {
	Hostname string `json:"Hostname"` // Server name
	Port     int    `json:"Port"`     // HTTP port
	Username string `json:"Username"`
	Password string `json:"Password"`
	Database string `json:"Database"`
}

// URL 得到 MongoDB 的URL
// [mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]
func (s MongodbInfo) URL() string {
	return fmt.Sprintf("mongodb://%s%s", s.userName(), s.address())
}

func (s MongodbInfo) address() string {
	if s.Port == 0 {
		return s.Hostname
	}

	return fmt.Sprintf("%s:%d", s.Hostname, s.Port)
}

func (s MongodbInfo) userName() string {
	if s.Username != "" && s.Password != "" {
		return s.Username + ":" + s.Password + "@"
	}

	return ""
}

// NewMongoClient ...
func NewMongoClient() (MgoSession, error) {
	return mgo.Dial(Cfg.MongoDB.URL())
}
