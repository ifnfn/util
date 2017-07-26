package system

import (
	"fmt"
	"log"
	"time"

	"roabay.com/util/config"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
	mgo "gopkg.in/mgo.v2"

	"gopkg.in/redis.v5"
)

var (
	SQL *sqlx.DB // SQL wrapper
)

type DBError struct {
	Text    string
	ErrorNo int
}

func (e *DBError) Error() string {
	return e.Text
}

// NewMySQL Connect to the database
func NewMySQL() (*sqlx.DB, error) {
	var err error
	var sql *sqlx.DB
	// Connect to MySQL
	if sql, err = sqlx.Connect("mysql", config.MySQL.DSN()); err != nil {
		// logger.Fatalf("SQL Driver Error: %s", err.Error())
	}
	sql = sql.Unsafe()

	// Check if is alive
	return sql, sql.Ping()
}

// NewMongoClient ...
func NewMongoClient() (*mgo.Session, error) {
	url := config.MongoDB.URL() + "/" + config.MongoDB.Database
	return mgo.Dial(url)
}

// NewRedisClient 初始化 Redis 服务器配置
func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         config.Redis.Address(),
		Password:     config.Redis.Password, // no password set
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolTimeout:  30 * time.Second,
		PoolSize:     10,
		DB:           0, // use default DB
	})
}

// DBInit MySQL 数据库初始化
func DBInit() {
	SQL, _ = NewMySQL()
}

// GetTableCount 根据过滤条件，查询指定表的记录数
func GetTableCount(tableName string, filter string) (int, error) {
	sql := "SELECT count(*) FROM " + tableName

	if filter != "" {
		sql += " WHERE " + filter
	}

	var count int
	if err := SQL.Get(&count, sql); err != nil {
		log.Fatal(err.Error())
		return -1, err
	}

	return count, nil
}

// GetSQL 合并 SQL 语句
func GetSQL(base, filter string, pageSize, pageID int) string {
	if filter != "" {
		base += " WHERE " + filter
	}
	if pageSize != 0 {
		base += fmt.Sprintf(" LIMIT %d,%d", pageID*pageSize, pageSize)
	}

	return base
}

// DeleteTableByID 根据 ID 删除记录
func DeleteTableByID(tableName string, id string) Error {
	return DeleteTable(tableName, fmt.Sprintf("id=\"%s\"", id))
}

// DeleteTableByKey ...
func DeleteTableByKey(tableName string, key, val string) Error {
	return DeleteTable(tableName, fmt.Sprintf("%s=\"%s\"", key, val))
}

// DeleteTable ...
func DeleteTable(tableName string, filter string) Error {
	sql := fmt.Sprintf("DELETE FROM %s WHERE %s;", tableName, filter)
	result, err := SQL.Exec(sql)
	if err != nil {
		// log.Warn(err.Error())
		return SqlQError(err)
	}

	affect, err := result.RowsAffected()
	if err != nil {
		// log.Warn(err.Error())
		return SqlQError(err)
	}

	if affect <= 0 {
		return NewQError(DelErrorID, fmt.Sprintf("No found record where = %s, from %s", filter, tableName))
	}

	return nil
}

// TruncateTable ...
func TruncateTable(tableName string) Error {
	sql := fmt.Sprintf("TRUNCATE TABLE %s;", tableName)
	_, err := SQL.Exec(sql)
	if err != nil {
		// log.Warn(err.Error())
		return SqlQError(err)
	}

	return nil
}
