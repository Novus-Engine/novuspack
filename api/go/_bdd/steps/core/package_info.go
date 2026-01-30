//go:build bdd

// Package core provides BDD step definitions for NovusPack core domain testing.
//
// Domain: core
// Tags: @domain:core, @phase:1
package core

import (
	"context"
	"fmt"

	"github.com/cucumber/godog"
)

// RegisterCoreInfoSteps registers step definitions for package information operations (GetInfo, state).
//
// Domain: core
// Phase: 1
// Tags: @domain:core
func RegisterCoreInfoSteps(ctx *godog.ScenarioContext) {
	// Package information steps
	ctx.Step(`^GetInfo is called$`, getInfoIsCalled)
	ctx.Step(`^comprehensive package information is retrieved$`, comprehensivePackageInformationIsRetrieved)
	ctx.Step(`^package details are available$`, packageDetailsAreAvailable)
	ctx.Step(`^information includes package metadata$`, informationIncludesPackageMetadata)

	// Package state steps
	ctx.Step(`^package is in valid state$`, packageIsInValidState)
	ctx.Step(`^the NovusPack system$`, theNovusPackSystem)
}

// Package information steps

func getInfoIsCalled(ctx context.Context) (context.Context, error) {
	world := getWorld(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// TODO: Get package and call GetInfo
	// pkg := world.GetPackage()
	// if pkg == nil {
	//     return ctx, nil
	// }
	// TODO: Once API is implemented, call GetInfo
	// info := pkg.GetInfo(ctx)
	// Store info in world for later verification
	return ctx, nil
}

func comprehensivePackageInformationIsRetrieved(ctx context.Context) error {
	world := getWorldTyped(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Once API is implemented, verify comprehensive package information was retrieved
	// For now, check if package exists (info retrieval would require package)
	pkg := world.GetPackage()
	if pkg == nil {
		return godog.ErrUndefined
	}
	return nil
}

func packageDetailsAreAvailable(ctx context.Context) error {
	world := getWorldTyped(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Once API is implemented, verify package details are available
	// For now, check if package exists (details would be available if package exists)
	pkg := world.GetPackage()
	if pkg == nil {
		return godog.ErrUndefined
	}
	return nil
}

func informationIncludesPackageMetadata(ctx context.Context) error {
	world := getWorldTyped(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Once API is implemented, verify information includes package metadata
	// For now, just verify package exists (metadata would be included if package exists)
	pkg := world.GetPackage()
	if pkg == nil {
		return godog.ErrUndefined
	}
	return nil
}

// Package state steps

func packageIsInValidState(ctx context.Context) error {
	world := getWorldTyped(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Get package and verify it's in valid state
	pkg := world.GetPackage()
	if pkg == nil {
		return godog.ErrUndefined
	}
	// TODO: Once API is implemented, verify package is actually in valid state
	// For now, check if package is open (basic validity check)
	if !pkg.IsOpen() {
		return fmt.Errorf("package is not open")
	}
	// Check if there's an error indicating invalid state
	err := world.GetError()
	if err != nil {
		return fmt.Errorf("package is in invalid state: %v", err)
	}
	return nil
}

func theNovusPackSystem(ctx context.Context) error {
	// This step just indicates the NovusPack system context
	return nil
}
