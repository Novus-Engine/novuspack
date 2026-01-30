//go:build bdd

package support

import (
	"context"

	"github.com/cucumber/godog"
	"github.com/novus-engine/novuspack/api/go/_bdd/contextkeys"
	compressionsteps "github.com/novus-engine/novuspack/api/go/_bdd/steps/compression"
	coresteps "github.com/novus-engine/novuspack/api/go/_bdd/steps/core"
	fileformatsteps "github.com/novus-engine/novuspack/api/go/_bdd/steps/file_format"
	filemgmtsteps "github.com/novus-engine/novuspack/api/go/_bdd/steps/file_mgmt"
	metadatasteps "github.com/novus-engine/novuspack/api/go/_bdd/steps/metadata"
	writingsteps "github.com/novus-engine/novuspack/api/go/_bdd/steps/writing"
	// stepspkg "github.com/novus-engine/novuspack/api/go/_bdd/steps" // Commented out - all step files backed up, will be reorganized
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
		return context.WithValue(ctx, contextkeys.WorldContextKey, world), nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		if w != nil {
			_ = w.Cleanup()
		}
		return ctx, nil
	})

	// Register common/shared steps first
	RegisterCommonSteps(ctx)

	// Register per-domain steps BEFORE generic patterns in RegisterCoreSteps
	// This ensures specific handlers match before generic consolidated patterns
	// TODO: Reorganize and restore these domains:
	// stepspkg.RegisterBasicOpsSteps(ctx)                    // Backed up - to be reorganized
	// File format domain - split registrations (COMPLETED)
	fileformatsteps.RegisterFileFormatHeaderSteps(ctx)
	fileformatsteps.RegisterFileFormatEntrySteps(ctx)
	fileformatsteps.RegisterFileFormatIndexSteps(ctx)
	fileformatsteps.RegisterFileFormatSignatureSteps(ctx)
	fileformatsteps.RegisterFileFormatParsingSteps(ctx)
	// stepspkg.RegisterFileTypesSteps(ctx)                   // Backed up - to be reorganized
	// File management - split registrations (COMPLETED)
	filemgmtsteps.RegisterFileMgmtAdditionSteps(ctx)
	filemgmtsteps.RegisterFileMgmtRemovalSteps(ctx)
	filemgmtsteps.RegisterFileMgmtExtractionSteps(ctx)
	filemgmtsteps.RegisterFileMgmtQuerySteps(ctx)
	filemgmtsteps.RegisterFileMgmtTagSteps(ctx)
	filemgmtsteps.RegisterFileMgmtSourceSteps(ctx)
	filemgmtsteps.RegisterFileMgmtPatterns(ctx)
	// Writing domain - SafeWrite operations (COMPLETED)
	writingsteps.RegisterSafeWriteSteps(ctx)
	// Compression domain - split registrations (COMPLETED)
	compressionsteps.RegisterCompressionOperationsSteps(ctx)
	compressionsteps.RegisterCompressionTypesSteps(ctx)
	compressionsteps.RegisterCompressionConfigurationSteps(ctx)
	compressionsteps.RegisterCompressionStreamingSteps(ctx)
	compressionsteps.RegisterCompressionPatternsSteps(ctx)
	// Metadata domain - DirectoryEntry tag operations (COMPLETED)
	metadatasteps.RegisterMetadataDirectoryEntryTagSteps(ctx)
	// Metadata domain - Package-level metadata operations (COMPLETED)
	metadatasteps.RegisterPackageMetadataSteps(ctx)
	// Metadata domain - Path metadata operations (COMPLETED)
	metadatasteps.RegisterPathMetadataOperationsSteps(ctx)
	// TODO: Reorganize and restore these domains:
	// stepspkg.RegisterSignaturesSteps(ctx)                   // Backed up - to be reorganized
	// stepspkg.RegisterStreamingSteps(ctx)                    // Backed up - to be reorganized
	// stepspkg.RegisterDedupSteps(ctx)                        // Backed up - to be reorganized
	// stepspkg.RegisterSecurityValidationSteps(ctx)           // Backed up - to be reorganized
	// stepspkg.RegisterGenericsSteps(ctx)                     // Backed up - to be reorganized
	// stepspkg.RegisterMetadataSteps(ctx)                     // Backed up - to be reorganized
	// stepspkg.RegisterSecurityEncryptionSteps(ctx)           // Backed up - to be reorganized
	// stepspkg.RegisterValidationSteps(ctx)                   // Backed up - to be reorganized
	// stepspkg.RegisterTestingSteps(ctx)                      // Backed up - to be reorganized
	// stepspkg.RegisterMetadataSystemSteps(ctx)               // Backed up - to be reorganized

	// Register core domain - split registrations (specific before generic)
	coresteps.RegisterCoreLifecycleSteps(ctx)
	coresteps.RegisterCoreOperationsSteps(ctx)
	coresteps.RegisterCoreInfoSteps(ctx)
	coresteps.RegisterCorePropertiesSteps(ctx)
	coresteps.RegisterFileInfoSteps(ctx)
	// Register core/generic patterns LAST so specific handlers take precedence
	coresteps.RegisterCoreGenericPatterns(ctx)
}
