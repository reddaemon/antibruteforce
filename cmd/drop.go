package cmd

import (
	"log"
	"time"

	api "github.com/reddaemon/antibruteforce/protofiles"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func init() { //nolint
	drop.PersistentFlags().StringVarP(&address, "address", "a", "localhost:8080", "server address")
	drop.PersistentFlags().StringVarP(&login, "login", "l", "", "login to drop")
	drop.PersistentFlags().StringVarP(&ip, "ip", "i", "", "ip to drop")
	rootCmd.AddCommand(drop)
}

var drop = &cobra.Command{ //nolint
	Use:   "drop",
	Short: "Drop",
	Long:  "Drop",
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("unable to connect: %v", err)
		}
		defer func(conn *grpc.ClientConn) {
			err := conn.Close()
			if err != nil {
				log.Fatalf("%s", err.Error())
			}
		}(conn)
		c := api.NewAntiBruteforceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := c.Drop(ctx, &api.DropRequest{Login: login, Ip: ip})
		if err != nil {
			log.Fatalf("unable to drop: %v", err)
		}
		log.Printf("Done: %t", r.Ok)
	},
}
