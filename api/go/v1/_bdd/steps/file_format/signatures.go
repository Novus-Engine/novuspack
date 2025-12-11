//go:build bdd

// Package file_format provides BDD step definitions for NovusPack file format domain testing.
//
// Domain: file_format
// Tags: @domain:file_format, @phase:2
package file_format

import (
	"bytes"
	"context"
	"fmt"

	"github.com/cucumber/godog"
	novuspack "github.com/novus-engine/novuspack/api/go/v1"
	"github.com/novus-engine/novuspack/api/go/v1/_bdd/contextkeys"
)

// worldFileFormatSignature is an interface for world methods needed by signature steps
type worldFileFormatSignature interface {
	SetSignature(*novuspack.Signature)
	GetSignature() *novuspack.Signature
	SetError(error)
	GetError() error
	SetPackageMetadata(string, interface{})
	GetPackageMetadata(string) (interface{}, bool)
}

// getWorldFileFormatSignature extracts the World from context with signature methods
func getWorldFileFormatSignature(ctx context.Context) worldFileFormatSignature {
	w := ctx.Value(contextkeys.WorldContextKey)
	if w == nil {
		return nil
	}
	if wf, ok := w.(worldFileFormatSignature); ok {
		return wf
	}
	return nil
}

// Helper functions are defined in file_entry.go to avoid duplication

// RegisterFileFormatSignatureSteps registers step definitions for signature parsing.
func RegisterFileFormatSignatureSteps(ctx *godog.ScenarioContext) {
	// Signature type steps
	ctx.Step(`^a signature block$`, aSignatureBlock)
	ctx.Step(`^SignatureType is examined$`, signatureTypeIsExamined)
	ctx.Step(`^SignatureType equals 0x01 for ML-DSA$`, signatureTypeEquals0x01ForMLDSA)
	ctx.Step(`^SignatureType equals 0x02 for SLH-DSA$`, signatureTypeEquals0x02ForSLHDSA)
	ctx.Step(`^SignatureType equals 0x03 for PGP$`, signatureTypeEquals0x03ForPGP)
	ctx.Step(`^SignatureType equals 0x04 for X\.509$`, signatureTypeEquals0x04ForX509)
	ctx.Step(`^SignatureType is a 32-bit unsigned integer$`, signatureTypeIsA32BitUnsignedInteger)
	ctx.Step(`^a NovusPack file containing a signature block$`, aNovusPackFileContainingASignatureBlock)
	ctx.Step(`^a NovusPack file containing multiple signature blocks$`, aNovusPackFileContainingMultipleSignatureBlocks)

	// Signature ReadFrom/WriteTo/NewSignature steps
	ctx.Step(`^NewSignature is called$`, newSignatureIsCalled)
	ctx.Step(`^a Signature is returned$`, aSignatureIsReturned)
	ctx.Step(`^Signature is in initialized state$`, signatureIsInInitializedState)
	ctx.Step(`^signature all fields are zero or empty$`, allSignatureFieldsAreZeroOrEmpty)
	ctx.Step(`^a Signature with values$`, aSignatureWithValues)
	ctx.Step(`^a Signature with all fields set$`, aSignatureWithAllFieldsSet)
	ctx.Step(`^a Signature without comment$`, aSignatureWithoutComment)
	ctx.Step(`^signature WriteTo is called with writer$`, signatureWriteToIsCalledWithWriter)
	ctx.Step(`^signature is written to writer$`, signatureIsWrittenToWriter)
	ctx.Step(`^header is written first \(18 bytes\)$`, signatureHeaderIsWrittenFirst18Bytes)
	ctx.Step(`^comment follows header if present$`, commentFollowsHeaderIfPresent)
	ctx.Step(`^signature data follows comment$`, signatureDataFollowsComment)
	ctx.Step(`^written data matches signature content$`, writtenDataMatchesSignatureContent)
	ctx.Step(`^a reader with valid signature data$`, aReaderWithValidSignatureData)
	ctx.Step(`^signature ReadFrom is called with reader$`, signatureReadFromIsCalledWithReader)
	ctx.Step(`^signature is read from reader$`, signatureIsReadFromReader)
	ctx.Step(`^signature fields match reader data$`, signatureFieldsMatchReaderData)
	ctx.Step(`^signature is valid$`, signatureIsValid)
	ctx.Step(`^all signature fields are preserved$`, allSignatureFieldsArePreserved)
	ctx.Step(`^CommentLength equals 0$`, commentLengthEquals0)
	ctx.Step(`^a reader with less than 18 bytes of signature data$`, aReaderWithLessThan18BytesOfSignatureData)
	ctx.Step(`^error indicates read failure$`, signatureErrorIndicatesReadFailure)
}

// Signature type step implementations

func aSignatureBlock(ctx context.Context) error {
	// TODO: Create a signature block
	return nil
}

func signatureTypeIsExamined(ctx context.Context) error {
	// TODO: Examine SignatureType
	return nil
}

func signatureTypeEquals0x01ForMLDSA(ctx context.Context) error {
	// TODO: Verify SignatureType equals 0x01 for ML-DSA
	return nil
}

func signatureTypeEquals0x02ForSLHDSA(ctx context.Context) error {
	// TODO: Verify SignatureType equals 0x02 for SLH-DSA
	return nil
}

func signatureTypeEquals0x03ForPGP(ctx context.Context) error {
	// TODO: Verify SignatureType equals 0x03 for PGP
	return nil
}

func signatureTypeEquals0x04ForX509(ctx context.Context) error {
	// TODO: Verify SignatureType equals 0x04 for X.509
	return nil
}

func signatureTypeIsA32BitUnsignedInteger(ctx context.Context) error {
	// TODO: Verify SignatureType is a 32-bit unsigned integer
	return nil
}

func aNovusPackFileContainingASignatureBlock(ctx context.Context) error {
	// TODO: Create a NovusPack file containing a signature block
	return godog.ErrPending
}

func aNovusPackFileContainingMultipleSignatureBlocks(ctx context.Context) error {
	// TODO: Create a NovusPack file containing multiple signature blocks
	return godog.ErrPending
}

// Signature ReadFrom/WriteTo/NewSignature step implementations

func newSignatureIsCalled(ctx context.Context) (context.Context, error) {
	world := getWorldFileFormatSignature(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	sig := novuspack.NewSignature()
	world.SetSignature(sig)
	return ctx, nil
}

func aSignatureIsReturned(ctx context.Context) error {
	world := getWorldFileFormatSignature(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	sig := world.GetSignature()
	if sig == nil {
		return fmt.Errorf("no Signature returned")
	}
	return nil
}

func signatureIsInInitializedState(ctx context.Context) error {
	world := getWorldFileFormatSignature(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	sig := world.GetSignature()
	if sig == nil {
		return fmt.Errorf("no signature available")
	}
	// Verify initialization - all fields should be zero or empty
	if sig.SignatureType != 0 {
		return fmt.Errorf("SignatureType = %d, want 0", sig.SignatureType)
	}
	if sig.SignatureSize != 0 {
		return fmt.Errorf("SignatureSize = %d, want 0", sig.SignatureSize)
	}
	if sig.CommentLength != 0 {
		return fmt.Errorf("CommentLength = %d, want 0", sig.CommentLength)
	}
	if len(sig.SignatureData) != 0 {
		return fmt.Errorf("SignatureData length = %d, want 0", len(sig.SignatureData))
	}
	return nil
}

func allSignatureFieldsAreZeroOrEmpty(ctx context.Context) error {
	return signatureIsInInitializedState(ctx)
}

func aSignatureWithValues(ctx context.Context) error {
	world := getWorldFileFormatSignature(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Create a signature with test values
	sig := &novuspack.Signature{
		SignatureType:      novuspack.SignatureTypeMLDSA,
		SignatureSize:      64,
		SignatureFlags:     0,
		SignatureTimestamp: 1638360000,
		CommentLength:      0,
		SignatureData:      make([]byte, 64),
	}
	sig.SignatureSize = uint32(len(sig.SignatureData))
	// Store original for round-trip comparison
	world.SetPackageMetadata("signature_original", sig)
	world.SetSignature(sig)
	return nil
}

func aSignatureWithAllFieldsSet(ctx context.Context) error {
	world := getWorldFileFormatSignature(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Create a signature with all fields set
	sig := &novuspack.Signature{
		SignatureType:      novuspack.SignatureTypeMLDSA,
		SignatureSize:      128,
		SignatureFlags:     0x0101,
		SignatureTimestamp: 1638360000,
		CommentLength:      20,
		SignatureComment:   "Test signature comment",
		SignatureData:      make([]byte, 128),
	}
	world.SetSignature(sig)
	return nil
}

func aSignatureWithoutComment(ctx context.Context) error {
	world := getWorldFileFormatSignature(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Create a signature without comment
	sig := &novuspack.Signature{
		SignatureType:      novuspack.SignatureTypeMLDSA,
		SignatureSize:      64,
		SignatureFlags:     0,
		SignatureTimestamp: 1638360000,
		CommentLength:      0,
		SignatureData:      make([]byte, 64),
	}
	sig.SignatureSize = uint32(len(sig.SignatureData))
	// Store original for round-trip comparison
	world.SetPackageMetadata("signature_original", sig)
	world.SetSignature(sig)
	return nil
}

func signatureWriteToIsCalledWithWriter(ctx context.Context) (context.Context, error) {
	world := getWorldFileFormatSignature(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	sig := world.GetSignature()
	if sig == nil {
		return ctx, fmt.Errorf("no signature available")
	}
	// Serialize using WriteTo
	var buf bytes.Buffer
	_, err := sig.WriteTo(&buf)
	if err != nil {
		world.SetError(err)
		return ctx, fmt.Errorf("WriteTo failed: %w", err)
	}
	// Store serialized data
	world.SetPackageMetadata("signature_serialized", buf.Bytes())
	return ctx, nil
}

func signatureIsWrittenToWriter(ctx context.Context) error {
	world := getWorldFileFormatSignature(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Verify serialized data exists
	data, exists := world.GetPackageMetadata("signature_serialized")
	if !exists {
		return fmt.Errorf("signature was not serialized")
	}
	if _, ok := data.([]byte); !ok {
		return fmt.Errorf("serialized data is not a byte slice")
	}
	return nil
}

func signatureHeaderIsWrittenFirst18Bytes(ctx context.Context) error {
	world := getWorldFileFormatSignature(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Verify header is written first
	return nil
}

func commentFollowsHeaderIfPresent(ctx context.Context) error {
	world := getWorldFileFormatSignature(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Verify comment follows header if present
	sig := world.GetSignature()
	if sig == nil {
		return fmt.Errorf("no signature available")
	}
	if sig.CommentLength > 0 {
		// Comment should be after the 18-byte header
		data, exists := world.GetPackageMetadata("signature_serialized")
		if !exists {
			return fmt.Errorf("signature was not serialized")
		}
		buf, ok := data.([]byte)
		if !ok {
			return fmt.Errorf("serialized data is not a byte slice")
		}
		if len(buf) < 18+int(sig.CommentLength) {
			return fmt.Errorf("serialized data does not include comment")
		}
	}
	return nil
}

func signatureDataFollowsComment(ctx context.Context) error {
	world := getWorldFileFormatSignature(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Verify signature data follows comment
	sig := world.GetSignature()
	if sig == nil {
		return fmt.Errorf("no signature available")
	}
	data, exists := world.GetPackageMetadata("signature_serialized")
	if !exists {
		return fmt.Errorf("signature was not serialized")
	}
	buf, ok := data.([]byte)
	if !ok {
		return fmt.Errorf("serialized data is not a byte slice")
	}
	// Signature data should be after header (18) + comment (CommentLength)
	expectedMinSize := 18 + int(sig.CommentLength) + int(sig.SignatureSize)
	if len(buf) < expectedMinSize {
		return fmt.Errorf("serialized data does not include signature data")
	}
	return nil
}

func writtenDataMatchesSignatureContent(ctx context.Context) error {
	world := getWorldFileFormatSignature(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Verify written data matches signature content
	originalSig := world.GetSignature()
	if originalSig == nil {
		return fmt.Errorf("no original signature available")
	}
	data, exists := world.GetPackageMetadata("signature_serialized")
	if !exists {
		return fmt.Errorf("signature was not serialized")
	}
	buf, ok := data.([]byte)
	if !ok {
		return fmt.Errorf("serialized data is not a byte slice")
	}
	// Deserialize and compare
	var readSig novuspack.Signature
	_, err := readSig.ReadFrom(bytes.NewReader(buf))
	if err != nil {
		return fmt.Errorf("failed to read back serialized data: %w", err)
	}
	// Compare key fields
	if readSig.SignatureType != originalSig.SignatureType {
		return fmt.Errorf("SignatureType mismatch: %d != %d", readSig.SignatureType, originalSig.SignatureType)
	}
	if readSig.SignatureSize != originalSig.SignatureSize {
		return fmt.Errorf("SignatureSize mismatch: %d != %d", readSig.SignatureSize, originalSig.SignatureSize)
	}
	return nil
}

func aReaderWithValidSignatureData(ctx context.Context) (context.Context, error) {
	world := getWorldFileFormatSignature(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// Create a valid signature and serialize it
	sig := &novuspack.Signature{
		SignatureType:      novuspack.SignatureTypeMLDSA,
		SignatureSize:      64,
		SignatureFlags:     0,
		SignatureTimestamp: 1638360000,
		CommentLength:      0,
		SignatureData:      make([]byte, 64),
	}
	sig.SignatureSize = uint32(len(sig.SignatureData))
	var buf bytes.Buffer
	_, err := sig.WriteTo(&buf)
	if err != nil {
		return ctx, fmt.Errorf("failed to serialize signature: %w", err)
	}
	world.SetPackageMetadata("signature_reader_data", buf.Bytes())
	return ctx, nil
}

func signatureReadFromIsCalledWithReader(ctx context.Context) (context.Context, error) {
	world := getWorldFileFormatSignature(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// Get reader data
	data, exists := world.GetPackageMetadata("signature_reader_data")
	if !exists {
		return ctx, fmt.Errorf("no reader data available")
	}
	buf, ok := data.([]byte)
	if !ok {
		return ctx, fmt.Errorf("reader data is not a byte slice")
	}
	// Read signature using ReadFrom
	sig := &novuspack.Signature{}
	_, err := sig.ReadFrom(bytes.NewReader(buf))
	if err != nil {
		world.SetError(err)
		// Return nil to allow error scenarios to continue and check for the error
		return ctx, nil
	}
	world.SetSignature(sig)
	return ctx, nil
}

func signatureIsReadFromReader(ctx context.Context) error {
	world := getWorldFileFormatSignature(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	sig := world.GetSignature()
	if sig == nil {
		return fmt.Errorf("no signature available")
	}
	// TODO: Verify ReadFrom was called and succeeded
	return nil
}

func signatureFieldsMatchReaderData(ctx context.Context) error {
	world := getWorldFileFormatSignature(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Verify signature fields match the original data
	sig := world.GetSignature()
	if sig == nil {
		return fmt.Errorf("no signature available")
	}
	// Verify signature has expected values from the test data
	if sig.SignatureType != novuspack.SignatureTypeMLDSA {
		return fmt.Errorf("SignatureType mismatch: %d != %d", sig.SignatureType, novuspack.SignatureTypeMLDSA)
	}
	if sig.SignatureSize != 64 {
		return fmt.Errorf("SignatureSize mismatch: %d != 64", sig.SignatureSize)
	}
	return nil
}

func signatureIsValid(ctx context.Context) error {
	world := getWorldFileFormatSignature(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	sig := world.GetSignature()
	if sig == nil {
		return fmt.Errorf("no signature available")
	}
	err := sig.Validate()
	if err != nil {
		world.SetError(err)
		return fmt.Errorf("signature validation failed: %w", err)
	}
	return nil
}

func allSignatureFieldsArePreserved(ctx context.Context) error {
	world := getWorldFileFormatSignature(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Get original signature
	originalData, exists := world.GetPackageMetadata("signature_original")
	if !exists {
		// If no original stored, just verify current signature is valid
		return signatureIsValid(ctx)
	}
	originalSig, ok := originalData.(*novuspack.Signature)
	if !ok {
		return fmt.Errorf("original signature is not a Signature")
	}
	// Get deserialized signature
	readSig := world.GetSignature()
	if readSig == nil {
		return fmt.Errorf("no deserialized signature available")
	}
	// Compare key fields
	if readSig.SignatureType != originalSig.SignatureType {
		return fmt.Errorf("SignatureType not preserved: %d != %d", readSig.SignatureType, originalSig.SignatureType)
	}
	if readSig.SignatureSize != originalSig.SignatureSize {
		return fmt.Errorf("SignatureSize not preserved: %d != %d", readSig.SignatureSize, originalSig.SignatureSize)
	}
	if readSig.CommentLength != originalSig.CommentLength {
		return fmt.Errorf("CommentLength not preserved: %d != %d", readSig.CommentLength, originalSig.CommentLength)
	}
	return nil
}

func commentLengthEquals0(ctx context.Context) error {
	world := getWorldFileFormatSignature(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	sig := world.GetSignature()
	if sig == nil {
		return fmt.Errorf("no signature available")
	}
	if sig.CommentLength != 0 {
		return fmt.Errorf("CommentLength = %d, want 0", sig.CommentLength)
	}
	return nil
}

func aReaderWithLessThan18BytesOfSignatureData(ctx context.Context) (context.Context, error) {
	world := getWorldFileFormatSignature(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// Create reader with less than 18 bytes
	incompleteData := make([]byte, 10)
	world.SetPackageMetadata("signature_reader_data", incompleteData)
	return ctx, nil
}

func signatureErrorIndicatesReadFailure(ctx context.Context) error {
	world := getWorldFileFormatSignature(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	err := world.GetError()
	if err == nil {
		return fmt.Errorf("expected error but got none")
	}
	errMsg := err.Error()
	if !contains(errMsg, "read") && !contains(errMsg, "incomplete") {
		return fmt.Errorf("error message '%s' does not indicate read failure", errMsg)
	}
	return nil
}
