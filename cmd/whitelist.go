package cmd

import (
	"log"
	"time"

	api "github.com/reddaemon/antibruteforce/protofiles"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// whitelistCmd represents the whitelist command
var whitelistCmd = &cobra.Command{ //nolint
	Use:   "whitelist",
	Short: "whitelist",
	Long:  `whitelist`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() { //nolint
	whitelistAdd.PersistentFlags().StringVarP(&address, "address", "a", "localhost:8080", "server address")
	whitelistAdd.PersistentFlags().StringVarP(&subnet, "subnet", "s", "", "subnet")

	whitelistRemove.PersistentFlags().StringVarP(&address, "address", "a", "localhost:8080", "server address")
	whitelistRemove.PersistentFlags().StringVarP(&subnet, "subnet", "s", "", "subnet")
	rootCmd.AddCommand(blacklistCmd)
	whitelistCmd.AddCommand(whitelistAdd, whitelistRemove)
}

var whitelistAdd = &cobra.Command{ //nolint
	Use:   "add",
	Short: "add to whitelist",
	Long:  "add to whitelist",
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
		r, err := c.AddToWhitelist(ctx, &api.AddToWhitelistRequest{Subnet: subnet})
		if err != nil {
			log.Fatalf("unable add to whitelist: %v", err)
		}
		log.Printf("Done: %t", r.Ok)
	},
}

var whitelistRemove = &cobra.Command{ //nolint
	Use:   "remove",
	Short: "remove",
	Long:  "remove",
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("unable to connect: %v", err)
		}
		c := api.NewAntiBruteforceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := c.RemoveFromWhitelist(ctx, &api.RemoveFromWhitelistRequest{Subnet: subnet})
		if err != nil {
			log.Fatalf("unable to remove from whitelist: %v", err)
		}
		log.Printf("Done: %t", r.Ok)
	},
}
