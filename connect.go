package dexp

import (
	_ "github.com/go-sql-driver/mysql"
	"os"
)

var DB *Database

func Mysql() error {
	var err error

	var (
		mysqlConnectionUrl = os.Getenv("MYSQL_URL")
		redisConnectionUrl = os.Getenv("REDIS_URL")
	)
	DB, err = NewDatabase(mysqlConnectionUrl, redisConnectionUrl)
	if err != nil {
		panic(err)
	}

	return err
}
