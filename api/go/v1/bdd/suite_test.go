package bdd

import (
	"flag"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	support "github.com/novus-engine/novuspack/api/go/v1/bdd/support"
)

var opt = godog.Options{
	Format: "pretty",
	Paths:  []string{"../../../../features"},
}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opt)
}

func TestFeatures(t *testing.T) {
	opt.Output = colors.Colored(os.Stdout)

	suite := godog.TestSuite{
		Name:                "novuspack",
		ScenarioInitializer: support.InitializeScenario,
		Options:             &opt,
	}

	if suite.Run() != 0 {
		t.Fail()
	}
}
