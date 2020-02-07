package main

import (
	_ "github.com/jackc/pgx/stdlib"
	_ "github.com/lib/pq"
	"github.com/reddaemon/antibrutforce/cmd"
)

func main() {
	cmd.Execute()
}
