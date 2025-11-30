package support

import (
	"context"

	"github.com/cucumber/godog"
	stepspkg "github.com/novus-engine/novuspack/api/go/v1/bdd/steps"
)

func InitializeScenario(ctx *godog.ScenarioContext) {
	var w *World

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		world, err := NewWorld()
		if err != nil {
			return ctx, err
		}
		w = world
		// Store world in context for step functions to access
		return context.WithValue(ctx, "world", world), nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		if w != nil {
			_ = w.Cleanup()
		}
		return ctx, nil
	})

	// Register common/shared steps first
	RegisterCommonSteps(ctx)

	// Register per-domain steps
	stepspkg.RegisterBasicOpsSteps(ctx)
	stepspkg.RegisterCoreSteps(ctx)
	stepspkg.RegisterFileFormatSteps(ctx)
	stepspkg.RegisterFileTypesSteps(ctx)
	stepspkg.RegisterFileMgmtSteps(ctx)
	stepspkg.RegisterWritingSteps(ctx)
	stepspkg.RegisterCompressionSteps(ctx)
	stepspkg.RegisterSignaturesSteps(ctx)
	stepspkg.RegisterStreamingSteps(ctx)
	stepspkg.RegisterDedupSteps(ctx)
	stepspkg.RegisterSecurityValidationSteps(ctx)
	stepspkg.RegisterGenericsSteps(ctx)
	stepspkg.RegisterMetadataSteps(ctx)
	stepspkg.RegisterSecurityEncryptionSteps(ctx)
	stepspkg.RegisterValidationSteps(ctx)
	stepspkg.RegisterTestingSteps(ctx)
	stepspkg.RegisterMetadataSystemSteps(ctx)
}
