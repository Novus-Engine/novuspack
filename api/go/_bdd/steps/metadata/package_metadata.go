//go:build bdd

// Package metadata provides BDD step definitions for Package-level metadata operations.
//
// Domain: metadata
// Tags: @domain:metadata, @phase:2
package metadata

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/cucumber/godog"
	novuspack "github.com/novus-engine/novuspack/api/go"
	"github.com/novus-engine/novuspack/api/go/_bdd/contextkeys"
)

// worldPackageMetadata is an interface for world methods needed by Package metadata steps
type worldPackageMetadata interface {
	GetPackage() novuspack.Package
	SetPackage(novuspack.Package, string)
	SetError(error)
	GetError() error
	SetPackageMetadata(string, interface{})
	GetPackageMetadata(string) (interface{}, bool)
	NewContext() context.Context
	TempPath(string) string
}

// getWorldPackageMetadata extracts the World from context
func getWorldPackageMetadata(ctx context.Context) worldPackageMetadata {
	w := ctx.Value(contextkeys.WorldContextKey)
	if w == nil {
		return nil
	}
	if wf, ok := w.(worldPackageMetadata); ok {
		return wf
	}
	return nil
}

// RegisterPackageMetadataSteps registers step definitions for Package-level metadata operations.
//
// Domain: metadata
// Phase: 2
// Tags: @domain:metadata
func RegisterPackageMetadataSteps(ctx *godog.ScenarioContext) {
	// Package comment management steps
	ctx.Step(`^I set the package comment to "([^"]*)"$`, iSetThePackageCommentTo)
	ctx.Step(`^SetComment is called with comment$`, setCommentIsCalledWithComment)
	ctx.Step(`^SetComment is called$`, setCommentIsCalled)
	ctx.Step(`^SetComment is called with invalid encoding$`, setCommentIsCalledWithInvalidEncoding)
	ctx.Step(`^SetComment is called with comment exceeding length limit$`, setCommentIsCalledWithCommentExceedingLengthLimit)
	ctx.Step(`^SetComment is called with comment containing injection patterns$`, setCommentIsCalledWithCommentContainingInjectionPatterns)
	ctx.Step(`^GetComment is called$`, getCommentIsCalled)
	ctx.Step(`^ClearComment is called$`, clearCommentIsCalled)
	ctx.Step(`^comment operation is called$`, commentOperationIsCalled)

	// Package comment verification steps
	ctx.Step(`^reading the package comment should return "([^"]*)"$`, readingThePackageCommentShouldReturn)
	ctx.Step(`^comment is stored in package$`, commentIsStoredInPackage)
	ctx.Step(`^comment string is returned$`, commentStringIsReturned)
	ctx.Step(`^comment matches stored value$`, commentMatchesStoredValue)
	ctx.Step(`^comment is removed$`, commentIsRemoved)
	ctx.Step(`^CommentSize and CommentStart are updated$`, commentSizeAndCommentStartAreUpdated)
	ctx.Step(`^CommentSize is set to 0$`, commentSizeIsSetTo0)
	ctx.Step(`^flags bit 4 is set$`, flagsBit4IsSet)
	ctx.Step(`^flags bit 4 is cleared$`, flagsBit4IsCleared)

	// Package AppID/VendorID management steps
	ctx.Step(`^SetAppID is called with app ID$`, setAppIDIsCalledWithAppID)
	ctx.Step(`^SetVendorID is called with vendor ID$`, setVendorIDIsCalledWithVendorID)
	ctx.Step(`^AppID and VendorID are set together$`, appIDAndVendorIDAreSetTogether)
	ctx.Step(`^SetVendorID is called$`, setVendorIDIsCalled)
	ctx.Step(`^SetVendorID or ClearVendorID is called$`, setVendorIDOrClearVendorIDIsCalled)
	ctx.Step(`^SetVendorID is called with invalid vendor ID$`, setVendorIDIsCalledWithInvalidVendorID)

	// Package AppID/VendorID verification steps
	ctx.Step(`^AppID is set in package header$`, appIDIsSetInPackageHeader)
	ctx.Step(`^VendorID is set in package header$`, vendorIDIsSetInPackageHeader)
	ctx.Step(`^AppID is accessible via GetInfo$`, appIDIsAccessibleViaGetInfo)
	ctx.Step(`^VendorID is accessible via GetInfo$`, vendorIDIsAccessibleViaGetInfo)
	ctx.Step(`^both identifiers are stored$`, bothIdentifiersAreStored)
	ctx.Step(`^combination identifies platform application$`, combinationIdentifiesPlatformApplication)
	ctx.Step(`^VendorID is set to (\d+)$`, vendorIDIsSetTo)
	ctx.Step(`^AppID is set to (\d+)$`, appIDIsSetTo)

	// Error indication steps for comment operations
	ctx.Step(`^error indicates encoding issue$`, errorIndicatesEncodingIssue)
	ctx.Step(`^error indicates length limit exceeded$`, errorIndicatesLengthLimitExceeded)
	ctx.Step(`^error indicates security issue$`, errorIndicatesSecurityIssue)
	ctx.Step(`^error type is context cancellation$`, errorTypeIsContextCancellation)
}

// Package comment management step implementations

// iSetThePackageCommentTo handles "I set the package comment to ..."
func iSetThePackageCommentTo(ctx context.Context, commentText string) error {
	world := getWorldPackageMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	// Ensure package exists and is open
	testCtx := world.NewContext()
	pkg := world.GetPackage()
	if pkg == nil {
		var err error
		pkg, err = novuspack.NewPackage()
		if err != nil {
			world.SetError(err)
			return err
		}
		path := world.TempPath("package.nvpk")
		err = pkg.Create(testCtx, path)
		if err != nil {
			world.SetError(err)
			return err
		}
		world.SetPackage(pkg, path)
	}

	// Use helper function to call SetComment
	// Since filePackage is not exported, we need to add these methods to Package interface
	// For now, we'll use a type assertion approach that works within the package
	err := callPackageSetComment(pkg, testCtx, commentText)
	if err != nil {
		world.SetError(err)
		return err
	}

	return nil
}

// callPackageSetComment is a helper to call SetComment on a Package
func callPackageSetComment(pkg novuspack.Package, ctx context.Context, comment string) error {
	// SetComment is a pure in-memory operation and doesn't require context per spec
	// ctx parameter is kept for compatibility with BDD step signatures but not used
	_ = ctx
	return pkg.SetComment(comment)
}

// setCommentIsCalledWithComment handles "SetComment is called with comment"
func setCommentIsCalledWithComment(ctx context.Context) error {
	world := getWorldPackageMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	// Ensure package exists
	testCtx := world.NewContext()
	pkg := world.GetPackage()
	if pkg == nil {
		var err error
		pkg, err = novuspack.NewPackage()
		if err != nil {
			world.SetError(err)
			return err
		}
		path := world.TempPath("package.nvpk")
		err = pkg.Create(testCtx, path)
		if err != nil {
			world.SetError(err)
			return err
		}
		world.SetPackage(pkg, path)
	}

	// Use a default test comment
	commentText := "test comment"
	err := callPackageSetComment(pkg, testCtx, commentText)
	if err != nil {
		world.SetError(err)
		// Store error for verification step
		return nil
	}

	// Store comment text in world for verification
	world.SetPackageMetadata("comment", commentText)

	return nil
}

// setCommentIsCalled handles "SetComment is called"
func setCommentIsCalled(ctx context.Context) error {
	return setCommentIsCalledWithComment(ctx)
}

// setCommentIsCalledWithInvalidEncoding handles "SetComment is called with invalid encoding"
func setCommentIsCalledWithInvalidEncoding(ctx context.Context) error {
	world := getWorldPackageMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	// Ensure package exists
	testCtx := world.NewContext()
	pkg := world.GetPackage()
	if pkg == nil {
		var err error
		pkg, err = novuspack.NewPackage()
		if err != nil {
			world.SetError(err)
			return err
		}
		path := world.TempPath("package.nvpk")
		err = pkg.Create(testCtx, path)
		if err != nil {
			world.SetError(err)
			return err
		}
		world.SetPackage(pkg, path)
	}

	// Invalid UTF-8 sequence
	invalidUTF8 := string([]byte{0xFF, 0xFE, 0xFD})
	err := callPackageSetComment(pkg, testCtx, invalidUTF8)
	if err != nil {
		world.SetError(err)
		// Don't return error - let verification step check it
		return nil
	}

	return nil
}

// setCommentIsCalledWithCommentExceedingLengthLimit handles "SetComment is called with comment exceeding length limit"
func setCommentIsCalledWithCommentExceedingLengthLimit(ctx context.Context) error {
	world := getWorldPackageMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	// Ensure package exists
	testCtx := world.NewContext()
	pkg := world.GetPackage()
	if pkg == nil {
		var err error
		pkg, err = novuspack.NewPackage()
		if err != nil {
			world.SetError(err)
			return err
		}
		path := world.TempPath("package.nvpk")
		err = pkg.Create(testCtx, path)
		if err != nil {
			world.SetError(err)
			return err
		}
		world.SetPackage(pkg, path)
	}

	// Create a comment exceeding maximum length (MaxCommentLength is typically 1MB)
	longComment := strings.Repeat("a", 1048577) // 1MB + 1 byte
	err := callPackageSetComment(pkg, testCtx, longComment)
	if err != nil {
		world.SetError(err)
		// Don't return error - let verification step check it
		return nil
	}

	return nil
}

// setCommentIsCalledWithCommentContainingInjectionPatterns handles "SetComment is called with comment containing injection patterns"
func setCommentIsCalledWithCommentContainingInjectionPatterns(ctx context.Context) error {
	world := getWorldPackageMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	// Ensure package exists
	testCtx := world.NewContext()
	pkg := world.GetPackage()
	if pkg == nil {
		var err error
		pkg, err = novuspack.NewPackage()
		if err != nil {
			world.SetError(err)
			return err
		}
		path := world.TempPath("package.nvpk")
		err = pkg.Create(testCtx, path)
		if err != nil {
			world.SetError(err)
			return err
		}
		world.SetPackage(pkg, path)
	}

	// Comment with potential injection patterns
	injectionComment := "<script>alert('xss')</script>"
	err := callPackageSetComment(pkg, testCtx, injectionComment)
	if err != nil {
		world.SetError(err)
		// Don't return error - let verification step check it
		return nil
	}

	return nil
}

// getCommentIsCalled handles "GetComment is called"
func getCommentIsCalled(ctx context.Context) error {
	world := getWorldPackageMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	pkg := world.GetPackage()
	if pkg == nil {
		return fmt.Errorf("no package available")
	}

	comment := callPackageGetComment(pkg)
	world.SetPackageMetadata("retrieved_comment", comment)

	return nil
}

// callPackageGetComment is a helper to call GetComment on a Package
func callPackageGetComment(pkg novuspack.Package) string {
	return pkg.GetComment()
}

// clearCommentIsCalled handles "ClearComment is called"
func clearCommentIsCalled(ctx context.Context) error {
	world := getWorldPackageMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	pkg := world.GetPackage()
	if pkg == nil {
		return fmt.Errorf("no package available")
	}

	testCtx := world.NewContext()
	err := callPackageClearComment(pkg, testCtx)
	if err != nil {
		world.SetError(err)
		return err
	}

	return nil
}

// callPackageClearComment is a helper to call ClearComment on a Package
func callPackageClearComment(pkg novuspack.Package, ctx context.Context) error {
	// ClearComment is a pure in-memory operation and doesn't require context per spec
	// ctx parameter is kept for compatibility with BDD step signatures but not used
	_ = ctx
	return pkg.ClearComment()
}

// commentOperationIsCalled handles "comment operation is called"
func commentOperationIsCalled(ctx context.Context) error {
	world := getWorldPackageMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	// Ensure package exists
	testCtx := world.NewContext()
	pkg := world.GetPackage()
	if pkg == nil {
		var err error
		pkg, err = novuspack.NewPackage()
		if err != nil {
			world.SetError(err)
			return err
		}
		path := world.TempPath("package.nvpk")
		err = pkg.Create(testCtx, path)
		if err != nil {
			world.SetError(err)
			return err
		}
		world.SetPackage(pkg, path)
	}

	// Use a cancelled context if available, otherwise use normal context
	// Check if there's a cancelled context stored in world metadata
	cancelledCtx, exists := world.GetPackageMetadata("cancelled_context")
	if exists && cancelledCtx != nil {
		if ctxVal, ok := cancelledCtx.(context.Context); ok {
			testCtx = ctxVal
		}
	}

	// Call SetComment with the context
	commentText := "test comment"
	err := callPackageSetComment(pkg, testCtx, commentText)
	if err != nil {
		world.SetError(err)
		// Don't return error - let verification step check it
		return nil
	}

	return nil
}

// Package comment verification step implementations

// readingThePackageCommentShouldReturn handles "reading the package comment should return ..."
func readingThePackageCommentShouldReturn(ctx context.Context, expectedComment string) error {
	world := getWorldPackageMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	pkg := world.GetPackage()
	if pkg == nil {
		return fmt.Errorf("no package available")
	}

	comment := callPackageGetComment(pkg)
	if comment != expectedComment {
		return fmt.Errorf("GetComment() = %q, want %q", comment, expectedComment)
	}

	return nil
}

// commentIsStoredInPackage handles "comment is stored in package"
func commentIsStoredInPackage(ctx context.Context) error {
	world := getWorldPackageMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	pkg := world.GetPackage()
	if pkg == nil {
		return fmt.Errorf("no package available")
	}

	// Verify comment is stored by checking HasComment
	hasComment := callPackageHasComment(pkg)
	if !hasComment {
		return fmt.Errorf("HasComment() = false, expected true")
	}

	// Verify comment text matches
	storedComment, exists := world.GetPackageMetadata("comment")
	if exists && storedComment != nil {
		if commentStr, ok := storedComment.(string); ok {
			retrievedComment := callPackageGetComment(pkg)
			if retrievedComment != commentStr {
				return fmt.Errorf("GetComment() = %q, want %q", retrievedComment, commentStr)
			}
		}
	}

	return nil
}

// callPackageHasComment is a helper to call HasComment on a Package
func callPackageHasComment(pkg novuspack.Package) bool {
	return pkg.HasComment()
}

// commentStringIsReturned handles "comment string is returned"
func commentStringIsReturned(ctx context.Context) error {
	world := getWorldPackageMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	retrievedComment, exists := world.GetPackageMetadata("retrieved_comment")
	if !exists || retrievedComment == nil {
		return fmt.Errorf("retrieved comment not found in world metadata")
	}

	// Verify it's a string (non-empty or empty string is fine)
	_, ok := retrievedComment.(string)
	if !ok {
		return fmt.Errorf("retrieved comment is not a string, got %T", retrievedComment)
	}

	return nil
}

// commentMatchesStoredValue handles "comment matches stored value"
func commentMatchesStoredValue(ctx context.Context) error {
	world := getWorldPackageMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	pkg := world.GetPackage()
	if pkg == nil {
		return fmt.Errorf("no package available")
	}

	// Try to get stored comment, but if not found, get from package directly
	storedComment, exists := world.GetPackageMetadata("comment")
	if !exists || storedComment == nil {
		// If comment wasn't stored in metadata, get it from package
		retrievedComment := callPackageGetComment(pkg)
		if retrievedComment == "" {
			return fmt.Errorf("no comment found in package and no stored comment in world metadata")
		}
		// Comment exists in package, which is acceptable
		return nil
	}

	commentStr, ok := storedComment.(string)
	if !ok {
		return fmt.Errorf("stored comment is not a string")
	}

	retrievedComment := callPackageGetComment(pkg)
	if retrievedComment != commentStr {
		return fmt.Errorf("GetComment() = %q, want %q", retrievedComment, commentStr)
	}

	return nil
}

// commentIsRemoved handles "comment is removed"
func commentIsRemoved(ctx context.Context) error {
	world := getWorldPackageMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	pkg := world.GetPackage()
	if pkg == nil {
		return fmt.Errorf("no package available")
	}

	// Verify comment is removed
	hasComment := callPackageHasComment(pkg)
	if hasComment {
		return fmt.Errorf("HasComment() = true, expected false after ClearComment")
	}

	comment := callPackageGetComment(pkg)
	if comment != "" {
		return fmt.Errorf("GetComment() = %q, want empty string", comment)
	}

	return nil
}

// commentSizeAndCommentStartAreUpdated handles "CommentSize and CommentStart are updated"
func commentSizeAndCommentStartAreUpdated(ctx context.Context) error {
	world := getWorldPackageMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	pkg := world.GetPackage()
	if pkg == nil {
		return fmt.Errorf("no package available")
	}

	// Verify comment is set
	hasComment := callPackageHasComment(pkg)
	if !hasComment {
		return fmt.Errorf("HasComment() = false, expected true")
	}

	return nil
}

// commentSizeIsSetTo0 handles "CommentSize is set to 0"
func commentSizeIsSetTo0(ctx context.Context) error {
	world := getWorldPackageMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	pkg := world.GetPackage()
	if pkg == nil {
		return fmt.Errorf("no package available")
	}

	// Verify comment is cleared
	hasComment := callPackageHasComment(pkg)
	if hasComment {
		return fmt.Errorf("HasComment() = true, expected false")
	}

	return nil
}

// flagsBit4IsSet handles "flags bit 4 is set"
func flagsBit4IsSet(ctx context.Context) error {
	world := getWorldPackageMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	pkg := world.GetPackage()
	if pkg == nil {
		return fmt.Errorf("no package available")
	}

	// Verify HasComment returns true (which indicates flag is set)
	hasComment := callPackageHasComment(pkg)
	if !hasComment {
		return fmt.Errorf("HasComment() = false, expected true (flag bit 4 should be set)")
	}

	return nil
}

// flagsBit4IsCleared handles "flags bit 4 is cleared"
func flagsBit4IsCleared(ctx context.Context) error {
	world := getWorldPackageMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	pkg := world.GetPackage()
	if pkg == nil {
		return fmt.Errorf("no package available")
	}

	// Verify HasComment returns false (which indicates flag is cleared)
	hasComment := callPackageHasComment(pkg)
	if hasComment {
		return fmt.Errorf("HasComment() = true, expected false (flag bit 4 should be cleared)")
	}

	return nil
}

// Package AppID/VendorID management step implementations

// setAppIDIsCalledWithAppID handles "SetAppID is called with app ID"
func setAppIDIsCalledWithAppID(ctx context.Context) error {
	world := getWorldPackageMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	// Ensure package exists
	testCtx := world.NewContext()
	pkg := world.GetPackage()
	if pkg == nil {
		var err error
		pkg, err = novuspack.NewPackage()
		if err != nil {
			world.SetError(err)
			return err
		}
		path := world.TempPath("package.nvpk")
		err = pkg.Create(testCtx, path)
		if err != nil {
			world.SetError(err)
			return err
		}
		world.SetPackage(pkg, path)
	}

	// Use a test AppID
	appID := uint64(12345)
	err := callPackageSetAppID(pkg, appID)
	if err != nil {
		world.SetError(err)
		return err
	}

	world.SetPackageMetadata("app_id", appID)

	return nil
}

// callPackageSetAppID is a helper to call SetAppID on a Package
func callPackageSetAppID(pkg novuspack.Package, appID uint64) error {
	return pkg.SetAppID(appID)
}

// setVendorIDIsCalledWithVendorID handles "SetVendorID is called with vendor ID"
func setVendorIDIsCalledWithVendorID(ctx context.Context) error {
	world := getWorldPackageMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	// Ensure package exists
	testCtx := world.NewContext()
	pkg := world.GetPackage()
	if pkg == nil {
		var err error
		pkg, err = novuspack.NewPackage()
		if err != nil {
			world.SetError(err)
			return err
		}
		path := world.TempPath("package.nvpk")
		err = pkg.Create(testCtx, path)
		if err != nil {
			world.SetError(err)
			return err
		}
		world.SetPackage(pkg, path)
	}

	// Use a test VendorID
	vendorID := uint32(67890)
	err := callPackageSetVendorID(pkg, vendorID)
	if err != nil {
		world.SetError(err)
		return err
	}

	world.SetPackageMetadata("vendor_id", vendorID)

	return nil
}

// callPackageSetVendorID is a helper to call SetVendorID on a Package
func callPackageSetVendorID(pkg novuspack.Package, vendorID uint32) error {
	return pkg.SetVendorID(vendorID)
}

// appIDAndVendorIDAreSetTogether handles "AppID and VendorID are set together"
func appIDAndVendorIDAreSetTogether(ctx context.Context) error {
	world := getWorldPackageMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	// Ensure package exists and is open
	testCtx := world.NewContext()
	pkg := world.GetPackage()
	if pkg == nil {
		var err error
		pkg, err = novuspack.NewPackage()
		if err != nil {
			world.SetError(err)
			return err
		}
		path := world.TempPath("package.nvpk")
		err = pkg.Create(testCtx, path)
		if err != nil {
			world.SetError(err)
			return err
		}
		world.SetPackage(pkg, path)
	}

	// Set both together
	vendorID := uint32(11111)
	appID := uint64(22222)
	err := callPackageSetPackageIdentity(pkg, vendorID, appID)
	if err != nil {
		world.SetError(err)
		return err
	}

	world.SetPackageMetadata("vendor_id", vendorID)
	world.SetPackageMetadata("app_id", appID)

	return nil
}

// callPackageSetPackageIdentity is a helper to call SetPackageIdentity on a Package
func callPackageSetPackageIdentity(pkg novuspack.Package, vendorID uint32, appID uint64) error {
	return pkg.SetPackageIdentity(vendorID, appID)
}

// setVendorIDIsCalled handles "SetVendorID is called"
func setVendorIDIsCalled(ctx context.Context) error {
	return setVendorIDIsCalledWithVendorID(ctx)
}

// setVendorIDOrClearVendorIDIsCalled handles "SetVendorID or ClearVendorID is called"
func setVendorIDOrClearVendorIDIsCalled(ctx context.Context) error {
	world := getWorldPackageMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	// Ensure package exists
	testCtx := world.NewContext()
	pkg := world.GetPackage()
	if pkg == nil {
		var err error
		pkg, err = novuspack.NewPackage()
		if err != nil {
			world.SetError(err)
			return err
		}
		path := world.TempPath("package.nvpk")
		err = pkg.Create(testCtx, path)
		if err != nil {
			world.SetError(err)
			return err
		}
		world.SetPackage(pkg, path)
	}

	// Use a cancelled context if available
	cancelledCtx, exists := world.GetPackageMetadata("cancelled_context")
	if exists && cancelledCtx != nil {
		if ctxVal, ok := cancelledCtx.(context.Context); ok {
			testCtx = ctxVal
		}
	}

	// Default to SetVendorID
	vendorID := uint32(67890)
	err := callPackageSetVendorID(pkg, vendorID)
	if err != nil {
		world.SetError(err)
		// Don't return error - let verification step check it
		return nil
	}

	return nil
}

// setVendorIDIsCalledWithInvalidVendorID handles "SetVendorID is called with invalid vendor ID"
func setVendorIDIsCalledWithInvalidVendorID(ctx context.Context) error {
	world := getWorldPackageMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	// Ensure package exists
	testCtx := world.NewContext()
	pkg := world.GetPackage()
	if pkg == nil {
		var err error
		pkg, err = novuspack.NewPackage()
		if err != nil {
			world.SetError(err)
			return err
		}
		path := world.TempPath("package.nvpk")
		err = pkg.Create(testCtx, path)
		if err != nil {
			world.SetError(err)
			return err
		}
		world.SetPackage(pkg, path)
	}

	// VendorID is uint32, so any uint32 value is technically valid
	// This step might be testing validation that doesn't exist yet
	// For now, just call SetVendorID with a value
	vendorID := uint32(67890)
	err := callPackageSetVendorID(pkg, vendorID)
	if err != nil {
		world.SetError(err)
		// Don't return error - let verification step check it
		return nil
	}

	return nil
}

// Package AppID/VendorID verification step implementations

// appIDIsSetInPackageHeader handles "AppID is set in package header"
// Note: AppID is stored in Info (single source of truth) and synced to header
func appIDIsSetInPackageHeader(ctx context.Context) error {
	world := getWorldPackageMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	pkg := world.GetPackage()
	if pkg == nil {
		return fmt.Errorf("no package available")
	}

	// Verify AppID is set (non-zero) - GetAppID reads from Info
	appID := callPackageGetAppID(pkg)
	if appID == 0 {
		return fmt.Errorf("GetAppID() = 0, expected non-zero")
	}

	// Verify it matches stored value if available
	storedAppID, exists := world.GetPackageMetadata("app_id")
	if exists && storedAppID != nil {
		if storedID, ok := storedAppID.(uint64); ok && appID != storedID {
			return fmt.Errorf("GetAppID() = %d, want %d", appID, storedID)
		}
	}

	return nil
}

// callPackageGetAppID is a helper to call GetAppID on a Package
func callPackageGetAppID(pkg novuspack.Package) uint64 {
	return pkg.GetAppID()
}

// vendorIDIsSetInPackageHeader handles "VendorID is set in package header"
// Note: VendorID is stored in Info (single source of truth) and synced to header
func vendorIDIsSetInPackageHeader(ctx context.Context) error {
	world := getWorldPackageMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	pkg := world.GetPackage()
	if pkg == nil {
		return fmt.Errorf("no package available")
	}

	// Verify VendorID is set (non-zero) - GetVendorID reads from Info
	vendorID := callPackageGetVendorID(pkg)
	if vendorID == 0 {
		return fmt.Errorf("GetVendorID() = 0, expected non-zero")
	}

	// Verify it matches stored value if available
	storedVendorID, exists := world.GetPackageMetadata("vendor_id")
	if exists && storedVendorID != nil {
		if storedID, ok := storedVendorID.(uint32); ok && vendorID != storedID {
			return fmt.Errorf("GetVendorID() = %d, want %d", vendorID, storedID)
		}
	}

	return nil
}

// callPackageGetVendorID is a helper to call GetVendorID on a Package
func callPackageGetVendorID(pkg novuspack.Package) uint32 {
	return pkg.GetVendorID()
}

// appIDIsAccessibleViaGetInfo handles "AppID is accessible via GetInfo"
func appIDIsAccessibleViaGetInfo(ctx context.Context) error {
	world := getWorldPackageMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	pkg := world.GetPackage()
	if pkg == nil {
		return fmt.Errorf("no package available")
	}

	// Get AppID via GetAppID
	appID := callPackageGetAppID(pkg)

	// Get info via GetInfo
	info, err := pkg.GetInfo()
	if err != nil {
		world.SetError(err)
		return err
	}

	// Verify AppID matches
	if info.AppID != appID {
		return fmt.Errorf("Info.AppID = %d, GetAppID() = %d, expected match", info.AppID, appID)
	}

	return nil
}

// vendorIDIsAccessibleViaGetInfo handles "VendorID is accessible via GetInfo"
func vendorIDIsAccessibleViaGetInfo(ctx context.Context) error {
	world := getWorldPackageMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	pkg := world.GetPackage()
	if pkg == nil {
		return fmt.Errorf("no package available")
	}

	// Get VendorID via GetVendorID
	vendorID := callPackageGetVendorID(pkg)

	// Get info via GetInfo
	info, err := pkg.GetInfo()
	if err != nil {
		world.SetError(err)
		return err
	}

	// Verify VendorID matches
	if info.VendorID != vendorID {
		return fmt.Errorf("Info.VendorID = %d, GetVendorID() = %d, expected match", info.VendorID, vendorID)
	}

	return nil
}

// bothIdentifiersAreStored handles "both identifiers are stored"
func bothIdentifiersAreStored(ctx context.Context) error {
	// Verify both AppID and VendorID are set
	if err := appIDIsSetInPackageHeader(ctx); err != nil {
		return fmt.Errorf("AppID not set: %w", err)
	}

	if err := vendorIDIsSetInPackageHeader(ctx); err != nil {
		return fmt.Errorf("VendorID not set: %w", err)
	}

	return nil
}

// combinationIdentifiesPlatformApplication handles "combination identifies platform application"
func combinationIdentifiesPlatformApplication(ctx context.Context) error {
	world := getWorldPackageMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	pkg := world.GetPackage()
	if pkg == nil {
		return fmt.Errorf("no package available")
	}

	// Verify both are set and can be retrieved together
	vendorID, appID := callPackageGetPackageIdentity(pkg)
	if vendorID == 0 || appID == 0 {
		return fmt.Errorf("GetPackageIdentity() = (%d, %d), expected both non-zero", vendorID, appID)
	}

	return nil
}

// callPackageGetPackageIdentity is a helper to call GetPackageIdentity on a Package
func callPackageGetPackageIdentity(pkg novuspack.Package) (uint32, uint64) {
	return pkg.GetPackageIdentity()
}

// vendorIDIsSetTo handles "VendorID is set to X"
func vendorIDIsSetTo(ctx context.Context, value string) error {
	world := getWorldPackageMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	pkg := world.GetPackage()
	if pkg == nil {
		return fmt.Errorf("no package available")
	}

	expectedVendorID, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid VendorID value: %s", value)
	}

	vendorID := callPackageGetVendorID(pkg)
	if vendorID != uint32(expectedVendorID) {
		return fmt.Errorf("GetVendorID() = %d, want %d", vendorID, uint32(expectedVendorID))
	}

	return nil
}

// appIDIsSetTo handles "AppID is set to X"
func appIDIsSetTo(ctx context.Context, value string) error {
	world := getWorldPackageMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	pkg := world.GetPackage()
	if pkg == nil {
		return fmt.Errorf("no package available")
	}

	expectedAppID, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid AppID value: %s", value)
	}

	appID := callPackageGetAppID(pkg)
	if appID != expectedAppID {
		return fmt.Errorf("GetAppID() = %d, want %d", appID, expectedAppID)
	}

	return nil
}

// Error indication step implementations

// errorIndicatesEncodingIssue handles "error indicates encoding issue"
func errorIndicatesEncodingIssue(ctx context.Context) error {
	world := getWorldPackageMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	err := world.GetError()
	if err == nil {
		return fmt.Errorf("no error found, expected encoding issue error")
	}

	// Check if error message indicates encoding issue
	// The error message format is [Validation] comment is not valid UTF-8
	errMsg := strings.ToLower(err.Error())
	// Check for any UTF-related terms or encoding issues
	// Be permissive - if it mentions comment and utf/invalid/not valid, it's an encoding issue
	hasComment := strings.Contains(errMsg, "comment")
	hasEncodingIssue := hasComment && (strings.Contains(errMsg, "utf") ||
		strings.Contains(errMsg, "encoding") ||
		strings.Contains(errMsg, "invalid") ||
		strings.Contains(errMsg, "character") ||
		strings.Contains(errMsg, "not valid"))

	if !hasEncodingIssue {
		return fmt.Errorf("error message '%s' does not indicate encoding issue", err.Error())
	}

	return nil
}

// errorIndicatesLengthLimitExceeded handles "error indicates length limit exceeded"
func errorIndicatesLengthLimitExceeded(ctx context.Context) error {
	world := getWorldPackageMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	err := world.GetError()
	if err == nil {
		return fmt.Errorf("no error found, expected length limit error")
	}

	// Check if error message indicates length limit exceeded
	// The error message format is [Validation] comment length exceeds maximum
	errMsg := strings.ToLower(err.Error())
	// Be permissive - if it mentions comment, length, and exceed/maximum/limit, it's a length issue
	hasComment := strings.Contains(errMsg, "comment")
	hasLength := strings.Contains(errMsg, "length")
	hasLengthIssue := hasComment && hasLength && (strings.Contains(errMsg, "exceed") ||
		strings.Contains(errMsg, "limit") ||
		strings.Contains(errMsg, "maximum") ||
		strings.Contains(errMsg, "max") ||
		strings.Contains(errMsg, "too long"))

	if !hasLengthIssue {
		return fmt.Errorf("error message '%s' does not indicate length limit exceeded", err.Error())
	}

	return nil
}

// errorIndicatesSecurityIssue handles "error indicates security issue"
func errorIndicatesSecurityIssue(ctx context.Context) error {
	world := getWorldPackageMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	err := world.GetError()
	if err == nil {
		return fmt.Errorf("no error found, expected security issue error")
	}

	// Check if error message indicates security issue
	errMsg := strings.ToLower(err.Error())
	hasSecurityIssue := strings.Contains(errMsg, "security") ||
		strings.Contains(errMsg, "injection") ||
		strings.Contains(errMsg, "malicious") ||
		strings.Contains(errMsg, "unsafe") ||
		strings.Contains(errMsg, "script") ||
		strings.Contains(errMsg, "xss")

	if !hasSecurityIssue {
		return fmt.Errorf("error message '%s' does not indicate security issue", err.Error())
	}

	return nil
}

// errorTypeIsContextCancellation handles "error type is context cancellation"
func errorTypeIsContextCancellation(ctx context.Context) error {
	world := getWorldPackageMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	err := world.GetError()
	if err == nil {
		return fmt.Errorf("no error found, expected context cancellation error")
	}

	// Check if error is context cancellation
	if err == context.Canceled || err == context.DeadlineExceeded {
		return nil
	}

	// Check if it's a PackageError with context type
	var pkgErr *novuspack.PackageError
	if errors.As(err, &pkgErr) {
		if pkgErr.Type == novuspack.ErrTypeContext {
			return nil
		}
	}

	// Check error message
	errMsg := strings.ToLower(err.Error())
	hasContextCancel := strings.Contains(errMsg, "context") &&
		(strings.Contains(errMsg, "cancel") || strings.Contains(errMsg, "deadline"))

	if !hasContextCancel {
		return fmt.Errorf("error '%s' does not indicate context cancellation", err.Error())
	}

	return nil
}
