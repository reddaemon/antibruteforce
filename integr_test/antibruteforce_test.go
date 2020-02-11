package tests

import (
	"context"
	"fmt"
	"time"

	"os"

	"github.com/cucumber/godog"
	api "github.com/reddaemon/antibruteforce/protofiles"
	"google.golang.org/grpc"
)

var grpcService = os.Getenv("GRPC_SERVICE")

func (a *apiTest) iCallGrpcMethod(method string) error {
	cc, err := grpc.Dial(grpcService, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("unable to connect: %v", err)
	}
	defer cc.Close()

	c := api.NewAntiBruteforceClient(cc)
	ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
	defer cancel()

	switch method {
	case "Auth":
		_, err = c.Auth(ctx,
			&api.AuthRequest{
				Login:    "a.login",
				Password: "a.password",
				Ip:       "a.ip",
			})
		a.responseError = err
	case "Drop":
		_, err = c.Drop(ctx,
			&api.DropRequest{
				Login: a.login,
				Ip:    a.ip,
			})
		a.responseError = err
	case "AddToBlacklist":
		_, err = c.AddToBlacklist(ctx,
			&api.AddToBlacklistRequest{
				Subnet: a.subnet,
			})
		a.responseError = err
	case "RemoveFromBlacklist":
		_, err = c.RemoveFromBlacklist(ctx,
			&api.RemoveFromBlacklistRequest{
				Subnet: a.subnet,
			})
	case "AddToWhitelist":
		_, err = c.AddToWhitelist(ctx, &api.AddToWhitelistRequest{
			Subnet: a.subnet})
		a.responseError = err

	case "RemoveFromWhitelist":
		_, err = c.RemoveFromWhitelist(ctx, &api.RemoveFromWhitelistRequest{
			Subnet: a.subnet})
		a.responseError = err

	}
	return nil
}

func (a *apiTest) responseErrorShouldBe(error string) error {

	if error != "nil" {
		error = "rpc error: code = Unknown desc = " + error
	}

	if error == "nil" && a.responseError != nil {
		return fmt.Errorf("unexpected error, expected %s, got %v", error, a.responseError)
	}

	if error != "nil" && a.responseError == nil {
		return fmt.Errorf("unexpected error, expected %s, got %v", error, nil)
	}

	if a.responseError != nil && error != a.responseError.Error() {
		return fmt.Errorf("unexpected error, expected %s, got %v", error, a.responseError.Error())
	}

	return nil

}

func ipIs(arg1 string) error {
	return godog.ErrPending
}

func FeatureContext(s *godog.Suite) {
	var t apiTest
	s.BeforeScenario(func(i interface{}) {
		t.login = ""
		t.password = ""
		t.ip = ""
		t.responseError = nil
	})

	s.Step(`^login is "([^"]*)"$`, t.loginIs)
	s.Step(`^password is "([^"]*)"$`, t.passwordIs)
	s.Step(`^ip is <"([^"]*)"$`, t.ipIs)
	s.Step(`^subnet is "([^"]*)"$`, t.subnetIs)

	s.Step(`^I call grpc method "([^"]*)"$`, t.iCallGrpcMethod)
	s.Step(`^response error should be "([^"]*)"$`, t.responseErrorShouldBe)

}
