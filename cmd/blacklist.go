package cmd

import (
	"log"
	"time"

	"context"

	api "github.com/reddaemon/antibruteforce/protofiles"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// blacklistCmd represents the blacklist command
var blacklistCmd = &cobra.Command{ // nolint
	Use:   "blacklist",
	Short: "blacklist",
	Long:  `blacklist`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() { // nolint
	blacklistAdd.PersistentFlags().StringVarP(&address, "address", "a", "localhost:8080", "server address")
	blacklistAdd.PersistentFlags().StringVarP(&subnet, "subnet", "s", "", "subnet")

	blacklistRemove.PersistentFlags().StringVarP(&address, "address", "a", "localhost:8080", "server address")
	blacklistRemove.PersistentFlags().StringVarP(&subnet, "subnet", "s", "", "subnet")
	rootCmd.AddCommand(blacklistCmd)
	blacklistCmd.AddCommand(blacklistAdd, blacklistRemove)
}

var blacklistAdd = &cobra.Command{ //nolint
	Use:   "add",
	Short: "add to blacklist",
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("unable to connect: %v", err)
		}
		defer conn.Close()
		c := api.NewAntiBruteforceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		r, err := c.AddToBlacklist(ctx, &api.AddToBlacklistRequest{Subnet: subnet})
		if err != nil {
			log.Fatalf("unable to add to blacklist: %v", err)
		}
		log.Printf("Done: %t", r.Ok)
	},
}

var blacklistRemove = &cobra.Command{ //nolint
	Use:   "remove",
	Short: "remove from blacklist",
	Long:  `remove from blacklist`,
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("unable to connect: %v", err)
		}
		defer conn.Close()
		c := api.NewAntiBruteforceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := c.RemoveFromBlacklist(ctx, &api.RemoveFromBlacklistRequest{Subnet: subnet})
		if err != nil {
			log.Fatalf("unable to remove from blacklist")
		}
		log.Printf("Done: %t", r.Ok)

	},
}
