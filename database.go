package dexp

import (
	"database/sql"
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"log"
	"strconv"
	"time"
)

type Database struct {
	DB          *sql.DB
	RedisPool   *redis.Pool
	Transaction *sql.Tx
}

func NewDatabase(mysqlUrl, redisUrl string) (*Database, error) {

	db, err := sql.Open("mysql", mysqlUrl)

	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(2)
	db.SetMaxIdleConns(1)
	return &Database{DB: db, RedisPool: NewPool(redisUrl)}, nil
}

func NewPool(add string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		// Dial or DialContext must be set. When both are set, DialContext takes precedence over Dial.
		Dial: func() (redis.Conn, error) { return redis.Dial("tcp", add) },
	}
}

func (d *Database) Begin() (*Database, error) {
	tx, err := d.DB.Begin()

	if err != nil {
		return nil, err
	}

	return &Database{Transaction: tx}, nil
}

func (d *Database) Rollback() error {
	if d.Transaction != nil {
		return d.Transaction.Rollback()
	}

	return nil
}
func (d *Database) Commit() error {

	if d.Transaction != nil {
		return d.Transaction.Commit()
	}

	return nil
}
func (d *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	if d.Transaction != nil {
		return d.Transaction.Exec(query, args...)
	}
	return d.DB.Exec(query, args...)
}

func (d *Database) QueryRow(query string, args ...interface{}) *sql.Row {
	if d.Transaction != nil {

		return d.Transaction.QueryRow(query, args...)
	}

	return d.DB.QueryRow(query, args...)
}

func (d *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {

	if d.Transaction != nil {
		return d.Transaction.Query(query, args...)
	}

	return d.DB.Query(query, args...)
}

func (d *Database) Select(table, alias string) *SelectQuery {

	return &SelectQuery{
		db:    d,
		table: &SelectQueryTable{Name: table, Alias: alias},
	}

}

func (d *Database) Update(table string) *UpdateQuery {
	return &UpdateQuery{
		db:     d,
		_table: table,
	}
}

func (d *Database) Insert(table string) *InsertQuery {

	return &InsertQuery{
		db:     d,
		_table: table,
	}
}

func getEntityCacheID(ID int64) string {
	return "entity:" + strconv.FormatInt(ID, 10)
}

func getUserCacheID(ID int64) string {
	return "user:" + strconv.FormatInt(ID, 10)
}

func (d *Database) SetUserCache(ID int64, entity interface{}) error {

	c := DB.RedisPool.Get()
	defer c.Close()

	b, err := json.Marshal(entity)

	if err != nil {
		return err
	}
	if _, err := c.Do("SET", getUserCacheID(ID), b); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (d *Database) GetUserCache(ID int64, des interface{}) error {
	c := DB.RedisPool.Get()
	defer c.Close()

	s, err := redis.String(c.Do("GET", getUserCacheID(ID)))

	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(s), &des)

	return nil

}

func (d *Database) SetEntityCache(ID int64, entity interface{}) error {

	c := DB.RedisPool.Get()
	defer c.Close()

	b, err := json.Marshal(entity)

	if err != nil {
		return err
	}
	if _, err := c.Do("SET", getEntityCacheID(ID), b); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (d *Database) GetEntityCache(ID int64, des interface{}) error {
	c := DB.RedisPool.Get()
	defer c.Close()

	s, err := redis.String(c.Do("GET", getEntityCacheID(ID)))

	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(s), &des)

	return nil

}

func (d *Database) publish(channel string, message string) error {

	c := DB.RedisPool.Get()
	defer c.Close()
	_, err := c.Do("publish", channel, message)
	return err

}
