package main

import (
	_ "github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
	"github.com/reddaemon/antibruteforce/cmd"
)

func main() {
	cmd.Execute()
}
