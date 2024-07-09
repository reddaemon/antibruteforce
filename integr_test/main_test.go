package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/cucumber/godog"
)

func TestMain(m *testing.M) {
	suite := godog.TestSuite{
		ScenarioInitializer: FeatureContext,
		Options: &godog.Options{
			Format:    "progress",
			Paths:     []string{"features"},
			Randomize: 0,
		},
	}

	if suite.Run() != 0 {
		fmt.Println("Failed")
		os.Exit(1)
	}
}
