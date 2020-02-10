package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/reddaemon/antibruteforce/config"
)

func GetDb(c *config.Config) (*sqlx.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		c.Db["host"], c.Db["port"], c.Db["user"], c.Db["pass"], c.Db["name"])

	return sqlx.Connect("postgres", psqlInfo)
}
