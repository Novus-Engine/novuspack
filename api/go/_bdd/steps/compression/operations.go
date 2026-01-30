//go:build bdd

// Package compression provides BDD step definitions for NovusPack compression domain testing.
//
// Domain: compression
// Tags: @domain:compression, @phase:3
package compression

import (
	"context"

	"github.com/cucumber/godog"
	"github.com/novus-engine/novuspack/api/go/_bdd/contextkeys"
)

// getWorld extracts the World from the context
func getWorld(ctx context.Context) interface{} {
	return ctx.Value(contextkeys.WorldContextKey)
}

// RegisterCompressionOperationsSteps registers step definitions for compression operations.
func RegisterCompressionOperationsSteps(ctx *godog.ScenarioContext) {
	// Compression operations
	ctx.Step(`^package is compressed$`, packageIsCompressed)
	ctx.Step(`^package is decompressed$`, packageIsDecompressed)
	ctx.Step(`^compression is applied$`, compressionIsApplied)
	ctx.Step(`^content matches original$`, contentMatchesOriginal)
}

func packageIsCompressed(ctx context.Context) (context.Context, error) {
	world := getWorld(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// TODO: Once API is implemented, compress package
	return ctx, nil
}

func packageIsDecompressed(ctx context.Context) (context.Context, error) {
	world := getWorld(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// TODO: Once API is implemented, decompress package
	return ctx, nil
}

func compressionIsApplied(ctx context.Context) (context.Context, error) {
	// TODO: Apply compression
	return ctx, nil
}

func contentMatchesOriginal(ctx context.Context) error {
	// TODO: Verify content matches original
	return nil
}
