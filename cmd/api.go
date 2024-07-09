package cmd

import (
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/reddaemon/antibruteforce/config"
	dbInstance "github.com/reddaemon/antibruteforce/db"
	"github.com/reddaemon/antibruteforce/internal/database/postgres"
	"github.com/reddaemon/antibruteforce/internal/service/api/bucket"
	"github.com/reddaemon/antibruteforce/internal/service/api/server"
	"github.com/reddaemon/antibruteforce/internal/service/api/usage"
	"github.com/reddaemon/antibruteforce/logger"
	grpcapi "github.com/reddaemon/antibruteforce/protofiles/protofiles/api"
	"github.com/spf13/cobra"
)

// apiCmd represents the api command
var apiCmd = &cobra.Command{ //nolint
	Use:   "api",
	Short: "start api",
	Long:  `start api`,
	Run: func(cmd *cobra.Command, args []string) {
		configPath := flag.String("config", "config.yml", "path to config file")
		flag.Parse()
		c, err := config.GetConfig(*configPath)
		if err != nil {
			log.Fatalf("unabler to get config: %v", err)

		}
		l, err := logger.GetLogger(c)
		if err != nil {
			log.Fatalf("unable to get logger: %v", err)
		}

		db, err := dbInstance.GetDb(c)
		if err != nil {
			log.Fatalf("unable to get db: %v", err)
		}
		lis, err := net.Listen("tcp", c.URL)
		if err != nil {
			l.Fatal(fmt.Sprintf("failed to listen %v", err))
		}
		l.Info("server started at " + c.URL)
		grpcServer := grpc.NewServer()

		if c.IsDev() {
			reflection.Register(grpcServer)
		}

		r := postgres.NewPsqlRepository(db, l)
		br := bucket.NewMemRepo(l)
		u := usage.NewUsage(r, br, l, c)

		grpcapi.RegisterAntiBruteforceServer(grpcServer, server.NewServer(u, l))

		err = grpcServer.Serve(lis)
		if err != nil {
			l.Fatal(err.Error())
		}
	},
}

func init() { //nolint
	rootCmd.AddCommand(apiCmd)
}
