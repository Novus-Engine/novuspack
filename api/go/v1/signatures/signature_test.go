package signatures

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"strings"
	"testing"
)

// TestSignatureValidation verifies validation logic
func TestSignatureValidation(t *testing.T) {
	tests := []struct {
		name      string
		signature Signature
		wantErr   bool
	}{
		{
			"Valid signature",
			Signature{
				SignatureType: SignatureTypeMLDSA,
				SignatureSize: 32,
				SignatureData: make([]byte, 32),
				CommentLength: 0,
			},
			false,
		},
		{
			"Zero signature type",
			Signature{
				SignatureType: 0,
				SignatureSize: 32,
				SignatureData: make([]byte, 32),
			},
			true,
		},
		{
			"Nil signature data",
			Signature{
				SignatureType: SignatureTypeMLDSA,
				SignatureSize: 32,
				SignatureData: nil,
			},
			true,
		},
		{
			"Size mismatch",
			Signature{
				SignatureType: SignatureTypeMLDSA,
				SignatureSize: 32,
				SignatureData: make([]byte, 16),
			},
			true,
		},
		{
			"With comment",
			Signature{
				SignatureType:    SignatureTypeMLDSA,
				SignatureSize:    32,
				SignatureData:    make([]byte, 32),
				CommentLength:    4,
				SignatureComment: "Test",
			},
			false,
		},
		{
			"Non-zero SignatureSize with zero CommentLength",
			Signature{
				SignatureType:    SignatureTypeMLDSA,
				SignatureSize:    32,
				SignatureData:    make([]byte, 32),
				CommentLength:    0,
				SignatureComment: "",
			},
			false,
		},
		{
			"Comment length mismatch",
			Signature{
				SignatureType:    SignatureTypeMLDSA,
				SignatureSize:    32,
				SignatureData:    make([]byte, 32),
				CommentLength:    10,
				SignatureComment: "Test",
			},
			true,
		},
		{
			"Comment length non-zero but comment is empty",
			Signature{
				SignatureType:    SignatureTypeMLDSA,
				SignatureSize:    32,
				SignatureData:    make([]byte, 32),
				CommentLength:    5,
				SignatureComment: "",
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.signature.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestSignatureSizeCalculation verifies size calculation
func TestSignatureSizeCalculation(t *testing.T) {
	tests := []struct {
		name          string
		signatureSize uint32
		commentLength uint16
		wantSize      int
	}{
		{"No comment", 32, 0, 50},
		{"With comment", 32, 10, 60},
		{"Large signature", 2420, 0, 2438},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sig := Signature{
				SignatureSize: tt.signatureSize,
				CommentLength: tt.commentLength,
			}

			if sig.Size() != tt.wantSize {
				t.Errorf("Size() = %d, want %d", sig.Size(), tt.wantSize)
			}
		})
	}
}

// TestSignatureFlags verifies flag operations
func TestSignatureFlags(t *testing.T) {
	sig := Signature{}

	// Test setting flag
	sig.SetFlag(0x01)
	if !sig.HasFlag(0x01) {
		t.Error("Expected flag 0x01 to be set")
	}

	// Test clearing flag
	sig.ClearFlag(0x01)
	if sig.HasFlag(0x01) {
		t.Error("Expected flag 0x01 to be cleared")
	}
}

// TestNewSignature verifies NewSignature initializes correctly
func TestNewSignature(t *testing.T) {
	sig := NewSignature()

	if sig == nil {
		t.Fatal("NewSignature() returned nil")
	}

	// Verify all fields are zero or empty
	if sig.SignatureType != 0 {
		t.Errorf("SignatureType = %d, want 0", sig.SignatureType)
	}
	if sig.SignatureSize != 0 {
		t.Errorf("SignatureSize = %d, want 0", sig.SignatureSize)
	}
	if sig.SignatureFlags != 0 {
		t.Errorf("SignatureFlags = %d, want 0", sig.SignatureFlags)
	}
	if sig.SignatureTimestamp != 0 {
		t.Errorf("SignatureTimestamp = %d, want 0", sig.SignatureTimestamp)
	}
	if sig.CommentLength != 0 {
		t.Errorf("CommentLength = %d, want 0", sig.CommentLength)
	}
	if len(sig.SignatureData) != 0 {
		t.Errorf("SignatureData length = %d, want 0", len(sig.SignatureData))
	}
}

// TestSignatureReadFrom verifies ReadFrom deserialization
// Specification: ../../docs/tech_specs/package_file_format.md Section 7.1 - Signature Structure
func TestSignatureReadFrom(t *testing.T) {
	tests := []struct {
		name    string
		sig     Signature
		wantErr bool
	}{
		{
			"Valid signature without comment",
			Signature{
				SignatureType:      SignatureTypeMLDSA,
				SignatureSize:      64,
				SignatureFlags:     0,
				SignatureTimestamp: 1638360000,
				CommentLength:      0,
				SignatureData:      make([]byte, 64),
			},
			false,
		},
		{
			"Valid signature with comment",
			Signature{
				SignatureType:      SignatureTypeMLDSA,
				SignatureSize:      128,
				SignatureFlags:     0x0101,
				SignatureTimestamp: 1638360000,
				CommentLength:      20,
				SignatureComment:   "Test signature comment",
				SignatureData:      make([]byte, 128),
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure CommentLength matches actual comment length
			tt.sig.CommentLength = uint16(len(tt.sig.SignatureComment))
			// Ensure SignatureSize matches actual data length
			tt.sig.SignatureSize = uint32(len(tt.sig.SignatureData))

			// Serialize using WriteTo
			var writeBuf bytes.Buffer
			_, writeErr := tt.sig.WriteTo(&writeBuf)
			if writeErr != nil {
				t.Fatalf("WriteTo() error = %v", writeErr)
			}

			// Deserialize using ReadFrom
			var sig Signature
			n, err := sig.ReadFrom(&writeBuf)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				expectedSize := tt.sig.Size()
				if n != int64(expectedSize) {
					t.Errorf("ReadFrom() read %d bytes, want %d", n, expectedSize)
				}

				// Verify all fields match
				if sig.SignatureType != tt.sig.SignatureType {
					t.Errorf("SignatureType = %d, want %d", sig.SignatureType, tt.sig.SignatureType)
				}
				if sig.SignatureSize != tt.sig.SignatureSize {
					t.Errorf("SignatureSize = %d, want %d", sig.SignatureSize, tt.sig.SignatureSize)
				}
				if sig.CommentLength != tt.sig.CommentLength {
					t.Errorf("CommentLength = %d, want %d", sig.CommentLength, tt.sig.CommentLength)
				}
				if sig.SignatureComment != tt.sig.SignatureComment {
					t.Errorf("SignatureComment = %q, want %q", sig.SignatureComment, tt.sig.SignatureComment)
				}
				if len(sig.SignatureData) != len(tt.sig.SignatureData) {
					t.Errorf("SignatureData length = %d, want %d", len(sig.SignatureData), len(tt.sig.SignatureData))
				}

				// Verify validation passes
				if err := sig.Validate(); err != nil {
					t.Errorf("ReadFrom() signature validation failed: %v", err)
				}

				// Verify SignatureFlags and SignatureTimestamp match
				if sig.SignatureFlags != tt.sig.SignatureFlags {
					t.Errorf("SignatureFlags = %d, want %d", sig.SignatureFlags, tt.sig.SignatureFlags)
				}
				if sig.SignatureTimestamp != tt.sig.SignatureTimestamp {
					t.Errorf("SignatureTimestamp = %d, want %d", sig.SignatureTimestamp, tt.sig.SignatureTimestamp)
				}
			}
		})
	}
}

// TestSignatureReadFromIncompleteData verifies ReadFrom handles incomplete data
func TestSignatureReadFromIncompleteData(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"No data", []byte{}},
		{"Partial header", make([]byte, 8)},
		{"Almost complete header", make([]byte, 17)},
		{"Header but no data", func() []byte {
			buf := new(bytes.Buffer)
			_ = binary.Write(buf, binary.LittleEndian, uint32(1))  // SignatureType
			_ = binary.Write(buf, binary.LittleEndian, uint32(64)) // SignatureSize
			_ = binary.Write(buf, binary.LittleEndian, uint32(0))  // SignatureFlags
			_ = binary.Write(buf, binary.LittleEndian, uint32(0))  // SignatureTimestamp
			_ = binary.Write(buf, binary.LittleEndian, uint16(0))  // CommentLength
			return buf.Bytes()                                     // Only 18 bytes, but SignatureSize says 64 bytes needed
		}()},
		{"Header with comment but incomplete comment", func() []byte {
			buf := new(bytes.Buffer)
			_ = binary.Write(buf, binary.LittleEndian, uint32(1))  // SignatureType
			_ = binary.Write(buf, binary.LittleEndian, uint32(64)) // SignatureSize
			_ = binary.Write(buf, binary.LittleEndian, uint32(0))  // SignatureFlags
			_ = binary.Write(buf, binary.LittleEndian, uint32(0))  // SignatureTimestamp
			_ = binary.Write(buf, binary.LittleEndian, uint16(10)) // CommentLength
			buf.WriteString("test")                                // Only 4 bytes of 10
			return buf.Bytes()
		}()},
		{"Header with comment but incomplete signature data", func() []byte {
			buf := new(bytes.Buffer)
			_ = binary.Write(buf, binary.LittleEndian, uint32(1))  // SignatureType
			_ = binary.Write(buf, binary.LittleEndian, uint32(64)) // SignatureSize
			_ = binary.Write(buf, binary.LittleEndian, uint32(0))  // SignatureFlags
			_ = binary.Write(buf, binary.LittleEndian, uint32(0))  // SignatureTimestamp
			_ = binary.Write(buf, binary.LittleEndian, uint16(10)) // CommentLength
			comment := "test comment"
			buf.Write([]byte(comment[:10])) // Exactly 10 bytes
			partialData := make([]byte, 30) // Only 30 bytes of 64 signature data
			buf.Write(partialData)
			return buf.Bytes()
		}()},
		{"Header with comment but incomplete signature data (exact boundary)", func() []byte {
			buf := new(bytes.Buffer)
			_ = binary.Write(buf, binary.LittleEndian, uint32(1))  // SignatureType
			_ = binary.Write(buf, binary.LittleEndian, uint32(64)) // SignatureSize
			_ = binary.Write(buf, binary.LittleEndian, uint32(0))  // SignatureFlags
			_ = binary.Write(buf, binary.LittleEndian, uint32(0))  // SignatureTimestamp
			_ = binary.Write(buf, binary.LittleEndian, uint16(10)) // CommentLength
			comment := "test comment"
			buf.Write([]byte(comment[:10])) // Exactly 10 bytes
			partialData := make([]byte, 63) // Only 63 bytes of 64 signature data (exact boundary)
			buf.Write(partialData)
			return buf.Bytes()
		}()},
		{"Header with comment but incomplete comment (exact boundary)", func() []byte {
			buf := new(bytes.Buffer)
			_ = binary.Write(buf, binary.LittleEndian, uint32(1))  // SignatureType
			_ = binary.Write(buf, binary.LittleEndian, uint32(64)) // SignatureSize
			_ = binary.Write(buf, binary.LittleEndian, uint32(0))  // SignatureFlags
			_ = binary.Write(buf, binary.LittleEndian, uint32(0))  // SignatureTimestamp
			_ = binary.Write(buf, binary.LittleEndian, uint16(10)) // CommentLength
			comment := "test comment"
			buf.Write([]byte(comment[:9])) // Only 9 bytes of 10 (exact boundary)
			return buf.Bytes()
		}()},
		{"Header with comment but no signature data when SignatureSize > 0", func() []byte {
			buf := new(bytes.Buffer)
			_ = binary.Write(buf, binary.LittleEndian, uint32(1))  // SignatureType
			_ = binary.Write(buf, binary.LittleEndian, uint32(64)) // SignatureSize
			_ = binary.Write(buf, binary.LittleEndian, uint32(0))  // SignatureFlags
			_ = binary.Write(buf, binary.LittleEndian, uint32(0))  // SignatureTimestamp
			_ = binary.Write(buf, binary.LittleEndian, uint16(10)) // CommentLength
			comment := "test comment"
			buf.Write([]byte(comment[:10])) // Exactly 10 bytes
			// No signature data
			return buf.Bytes()
		}()},
		{"Incomplete SignatureSize read", func() []byte {
			buf := new(bytes.Buffer)
			_ = binary.Write(buf, binary.LittleEndian, uint32(1)) // SignatureType
			// Only 2 bytes of SignatureSize (need 4)
			buf.Write([]byte{0x40, 0x00})
			return buf.Bytes()
		}()},
		{"Incomplete SignatureFlags read", func() []byte {
			buf := new(bytes.Buffer)
			_ = binary.Write(buf, binary.LittleEndian, uint32(1))  // SignatureType
			_ = binary.Write(buf, binary.LittleEndian, uint32(64)) // SignatureSize
			// Only 2 bytes of SignatureFlags (need 4)
			buf.Write([]byte{0x00, 0x00})
			return buf.Bytes()
		}()},
		{"Incomplete SignatureTimestamp read", func() []byte {
			buf := new(bytes.Buffer)
			_ = binary.Write(buf, binary.LittleEndian, uint32(1))  // SignatureType
			_ = binary.Write(buf, binary.LittleEndian, uint32(64)) // SignatureSize
			_ = binary.Write(buf, binary.LittleEndian, uint32(0))  // SignatureFlags
			// Only 2 bytes of SignatureTimestamp (need 4)
			buf.Write([]byte{0x00, 0x00})
			return buf.Bytes()
		}()},
		{"Incomplete CommentLength read", func() []byte {
			buf := new(bytes.Buffer)
			_ = binary.Write(buf, binary.LittleEndian, uint32(1))  // SignatureType
			_ = binary.Write(buf, binary.LittleEndian, uint32(64)) // SignatureSize
			_ = binary.Write(buf, binary.LittleEndian, uint32(0))  // SignatureFlags
			_ = binary.Write(buf, binary.LittleEndian, uint32(0))  // SignatureTimestamp
			// Only 1 byte of CommentLength (need 2)
			buf.Write([]byte{0x00})
			return buf.Bytes()
		}()},
		{"Valid signature with zero SignatureSize and zero CommentLength", func() []byte {
			buf := new(bytes.Buffer)
			_ = binary.Write(buf, binary.LittleEndian, uint32(1)) // SignatureType
			_ = binary.Write(buf, binary.LittleEndian, uint32(0)) // SignatureSize = 0
			_ = binary.Write(buf, binary.LittleEndian, uint32(0)) // SignatureFlags
			_ = binary.Write(buf, binary.LittleEndian, uint32(0)) // SignatureTimestamp
			_ = binary.Write(buf, binary.LittleEndian, uint16(0)) // CommentLength = 0
			// Complete header, no data (valid case)
			return buf.Bytes()
		}()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sig Signature
			r := bytes.NewReader(tt.data)
			_, err := sig.ReadFrom(r)

			// Check if this is a valid case (zero sizes)
			isValidZeroCase := strings.Contains(tt.name, "Valid signature with zero")
			if isValidZeroCase {
				if err != nil {
					t.Errorf("ReadFrom() expected success for valid zero-size signature, got error: %v", err)
				}
				// Verify the signature was read correctly
				if sig.SignatureType != 1 {
					t.Errorf("SignatureType = %d, want 1", sig.SignatureType)
				}
				if sig.SignatureSize != 0 {
					t.Errorf("SignatureSize = %d, want 0", sig.SignatureSize)
				}
				if sig.CommentLength != 0 {
					t.Errorf("CommentLength = %d, want 0", sig.CommentLength)
				}
			} else {
				if err == nil {
					t.Errorf("ReadFrom() expected error for incomplete data, got nil")
				}
			}
		})
	}
}

// errorReader is a reader that returns a non-EOF error
type errorReader struct {
	err error
}

func (r *errorReader) Read(p []byte) (int, error) {
	if r.err == nil {
		r.err = errors.New("read error")
	}
	return 0, r.err
}

// partialReader reads some bytes then returns an error
type partialReader struct {
	data []byte
	pos  int
	err  error
}

func (r *partialReader) Read(p []byte) (int, error) {
	if r.pos >= len(r.data) {
		if r.err == nil {
			r.err = errors.New("read error")
		}
		return 0, r.err
	}
	n := copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}

// TestSignatureReadFromNonEOFErrors verifies ReadFrom handles non-EOF errors
func TestSignatureReadFromNonEOFErrors(t *testing.T) {
	tests := []struct {
		name    string
		reader  io.Reader
		wantErr bool
	}{
		{
			"Error reader during SignatureType read",
			&errorReader{},
			true,
		},
		{
			"Error reader after partial SignatureType read",
			func() io.Reader {
				buf := new(bytes.Buffer)
				buf.Write([]byte{0x01, 0x00}) // Partial SignatureType (2 of 4 bytes)
				return io.MultiReader(buf, &errorReader{})
			}(),
			true,
		},
		{
			"Partial reader during SignatureSize read",
			func() io.Reader {
				buf := new(bytes.Buffer)
				_ = binary.Write(buf, binary.LittleEndian, uint32(1)) // SignatureType
				return io.MultiReader(buf, &partialReader{data: []byte{0x40, 0x00}, err: errors.New("read error")})
			}(),
			true,
		},
		{
			"Partial reader during SignatureFlags read",
			func() io.Reader {
				buf := new(bytes.Buffer)
				_ = binary.Write(buf, binary.LittleEndian, uint32(1))  // SignatureType
				_ = binary.Write(buf, binary.LittleEndian, uint32(64)) // SignatureSize
				return io.MultiReader(buf, &partialReader{data: []byte{0x00, 0x00}, err: errors.New("read error")})
			}(),
			true,
		},
		{
			"Partial reader during SignatureTimestamp read",
			func() io.Reader {
				buf := new(bytes.Buffer)
				_ = binary.Write(buf, binary.LittleEndian, uint32(1))  // SignatureType
				_ = binary.Write(buf, binary.LittleEndian, uint32(64)) // SignatureSize
				_ = binary.Write(buf, binary.LittleEndian, uint32(0))  // SignatureFlags
				return io.MultiReader(buf, &partialReader{data: []byte{0x00, 0x00}, err: errors.New("read error")})
			}(),
			true,
		},
		{
			"Partial reader during CommentLength read",
			func() io.Reader {
				buf := new(bytes.Buffer)
				_ = binary.Write(buf, binary.LittleEndian, uint32(1))  // SignatureType
				_ = binary.Write(buf, binary.LittleEndian, uint32(64)) // SignatureSize
				_ = binary.Write(buf, binary.LittleEndian, uint32(0))  // SignatureFlags
				_ = binary.Write(buf, binary.LittleEndian, uint32(0))  // SignatureTimestamp
				return io.MultiReader(buf, &partialReader{data: []byte{0x00}, err: errors.New("read error")})
			}(),
			true,
		},
		{
			"Error reader during comment read",
			func() io.Reader {
				buf := new(bytes.Buffer)
				_ = binary.Write(buf, binary.LittleEndian, uint32(1))  // SignatureType
				_ = binary.Write(buf, binary.LittleEndian, uint32(64)) // SignatureSize
				_ = binary.Write(buf, binary.LittleEndian, uint32(0))  // SignatureFlags
				_ = binary.Write(buf, binary.LittleEndian, uint32(0))  // SignatureTimestamp
				_ = binary.Write(buf, binary.LittleEndian, uint16(10)) // CommentLength = 10
				return io.MultiReader(buf, &errorReader{})
			}(),
			true,
		},
		{
			"Error reader during signature data read",
			func() io.Reader {
				buf := new(bytes.Buffer)
				_ = binary.Write(buf, binary.LittleEndian, uint32(1))  // SignatureType
				_ = binary.Write(buf, binary.LittleEndian, uint32(64)) // SignatureSize
				_ = binary.Write(buf, binary.LittleEndian, uint32(0))  // SignatureFlags
				_ = binary.Write(buf, binary.LittleEndian, uint32(0))  // SignatureTimestamp
				_ = binary.Write(buf, binary.LittleEndian, uint16(0))  // CommentLength = 0
				return io.MultiReader(buf, &errorReader{})
			}(),
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sig Signature
			_, err := sig.ReadFrom(tt.reader)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil {
				// For error readers, should not be EOF/incomplete
				if strings.Contains(tt.name, "Error reader") {
					errStr := err.Error()
					if strings.Contains(errStr, "EOF") || strings.Contains(errStr, "incomplete") {
						t.Errorf("ReadFrom() error = %q, want non-EOF error for error reader", errStr)
					}
				}
			}
		})
	}
}

// TestSignatureWriteTo verifies WriteTo serialization
// Specification: ../../docs/tech_specs/package_file_format.md Section 7.1 - Signature Structure
func TestSignatureWriteTo(t *testing.T) {
	tests := []struct {
		name    string
		sig     Signature
		wantErr bool
	}{
		{
			"Valid signature without comment",
			Signature{
				SignatureType:      SignatureTypeMLDSA,
				SignatureSize:      64,
				SignatureFlags:     0,
				SignatureTimestamp: 1638360000,
				CommentLength:      0,
				SignatureData:      make([]byte, 64),
			},
			false,
		},
		{
			"Valid signature with comment",
			Signature{
				SignatureType:      SignatureTypeMLDSA,
				SignatureSize:      128,
				SignatureFlags:     0x0101,
				SignatureTimestamp: 1638360000,
				CommentLength:      20,
				SignatureComment:   "Test signature comment",
				SignatureData:      make([]byte, 128),
			},
			false,
		},
		{
			"Valid signature with zero CommentLength",
			Signature{
				SignatureType:      SignatureTypeMLDSA,
				SignatureSize:      32,
				SignatureFlags:     0,
				SignatureTimestamp: 1638360000,
				CommentLength:      0,
				SignatureComment:   "",
				SignatureData:      make([]byte, 32),
			},
			false,
		},
		{
			"Valid signature with large data",
			Signature{
				SignatureType:      SignatureTypeMLDSA,
				SignatureSize:      256,
				SignatureFlags:     0x0101,
				SignatureTimestamp: 1638360000,
				CommentLength:      0,
				SignatureComment:   "",
				SignatureData:      make([]byte, 256),
			},
			false,
		},
		{
			"Valid signature with comment and large data",
			Signature{
				SignatureType:      SignatureTypeMLDSA,
				SignatureSize:      512,
				SignatureFlags:     0x0101,
				SignatureTimestamp: 1638360000,
				CommentLength:      30,
				SignatureComment:   "Comment for large signature data",
				SignatureData:      make([]byte, 512),
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure CommentLength matches actual comment length
			tt.sig.CommentLength = uint16(len(tt.sig.SignatureComment))
			// Ensure SignatureSize matches actual data length
			tt.sig.SignatureSize = uint32(len(tt.sig.SignatureData))

			var buf bytes.Buffer
			n, err := tt.sig.WriteTo(&buf)

			if (err != nil) != tt.wantErr {
				t.Errorf("WriteTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				expectedSize := tt.sig.Size()
				if n != int64(expectedSize) {
					t.Errorf("WriteTo() wrote %d bytes, want %d", n, expectedSize)
				}

				if buf.Len() != expectedSize {
					t.Errorf("WriteTo() buffer size = %d bytes, want %d", buf.Len(), expectedSize)
				}

				// Verify we can read it back
				var sig Signature
				_, readErr := sig.ReadFrom(&buf)
				if readErr != nil {
					t.Errorf("Failed to read back written data: %v", readErr)
				}

				if sig.SignatureType != tt.sig.SignatureType {
					t.Errorf("SignatureType mismatch: %d != %d", sig.SignatureType, tt.sig.SignatureType)
				}
				if sig.SignatureSize != tt.sig.SignatureSize {
					t.Errorf("SignatureSize mismatch: %d != %d", sig.SignatureSize, tt.sig.SignatureSize)
				}
			}
		})
	}
}

// TestSignatureRoundTrip verifies round-trip serialization
func TestSignatureRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		sig  Signature
	}{
		{
			"Signature without comment",
			Signature{
				SignatureType:      SignatureTypeMLDSA,
				SignatureSize:      64,
				SignatureFlags:     0,
				SignatureTimestamp: 1638360000,
				CommentLength:      0,
				SignatureData:      make([]byte, 64),
			},
		},
		{
			"Signature with comment",
			Signature{
				SignatureType:      SignatureTypeMLDSA,
				SignatureSize:      128,
				SignatureFlags:     0x0101,
				SignatureTimestamp: 1638360000,
				CommentLength:      25,
				SignatureComment:   "Test signature comment",
				SignatureData:      make([]byte, 128),
			},
		},
		{
			"PGP signature",
			Signature{
				SignatureType:      SignatureTypePGP,
				SignatureSize:      256,
				SignatureFlags:     0x0202,
				SignatureTimestamp: 1638361000,
				CommentLength:      0,
				SignatureData:      make([]byte, 256),
			},
		},
		{
			"Signature with maximum timestamp",
			Signature{
				SignatureType:      SignatureTypeMLDSA,
				SignatureSize:      64,
				SignatureFlags:     0xFFFF,
				SignatureTimestamp: 0xFFFFFFFF,
				CommentLength:      0,
				SignatureComment:   "",
				SignatureData:      make([]byte, 64),
			},
		},
		{
			"Signature with all fields set",
			Signature{
				SignatureType:      SignatureTypePGP,
				SignatureSize:      128,
				SignatureFlags:     0x0101,
				SignatureTimestamp: 1638360000,
				CommentLength:      50,
				SignatureComment:   strings.Repeat("Test", 12) + "XX", // 50 chars
				SignatureData:      make([]byte, 128),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure CommentLength matches actual comment length
			tt.sig.CommentLength = uint16(len(tt.sig.SignatureComment))
			// Ensure SignatureSize matches actual data length
			tt.sig.SignatureSize = uint32(len(tt.sig.SignatureData))

			// Write
			var buf bytes.Buffer
			if _, err := tt.sig.WriteTo(&buf); err != nil {
				t.Fatalf("WriteTo() error = %v", err)
			}

			// Read
			var sig Signature
			if _, err := sig.ReadFrom(&buf); err != nil {
				t.Fatalf("ReadFrom() error = %v", err)
			}

			// Compare all fields
			if sig.SignatureType != tt.sig.SignatureType {
				t.Errorf("SignatureType mismatch: %d != %d", sig.SignatureType, tt.sig.SignatureType)
			}
			if sig.SignatureSize != tt.sig.SignatureSize {
				t.Errorf("SignatureSize mismatch: %d != %d", sig.SignatureSize, tt.sig.SignatureSize)
			}
			if sig.SignatureFlags != tt.sig.SignatureFlags {
				t.Errorf("SignatureFlags mismatch: %d != %d", sig.SignatureFlags, tt.sig.SignatureFlags)
			}
			if sig.SignatureTimestamp != tt.sig.SignatureTimestamp {
				t.Errorf("SignatureTimestamp mismatch: %d != %d", sig.SignatureTimestamp, tt.sig.SignatureTimestamp)
			}
			if sig.CommentLength != tt.sig.CommentLength {
				t.Errorf("CommentLength mismatch: %d != %d", sig.CommentLength, tt.sig.CommentLength)
			}
			if sig.SignatureComment != tt.sig.SignatureComment {
				t.Errorf("SignatureComment mismatch: %q != %q", sig.SignatureComment, tt.sig.SignatureComment)
			}
			if len(sig.SignatureData) != len(tt.sig.SignatureData) {
				t.Errorf("SignatureData length mismatch: %d != %d", len(sig.SignatureData), len(tt.sig.SignatureData))
			}

			// Validate
			if err := sig.Validate(); err != nil {
				t.Errorf("Round-trip signature validation failed: %v", err)
			}
		})
	}
}

// Helper types are defined in fileformat/test_helpers.go
// We need to import them or define them here for the signatures package
// Since they're in a different package, we'll define them here
type failingWriter struct {
	maxBytes int
	written  int
}

func (w *failingWriter) Write(p []byte) (int, error) {
	if w.written >= w.maxBytes {
		return 0, errors.New("write failed")
	}
	remaining := w.maxBytes - w.written
	if len(p) > remaining {
		w.written = w.maxBytes
		return remaining, errors.New("write failed")
	}
	w.written += len(p)
	return len(p), nil
}

type errorWriter struct {
	err error
}

func (w *errorWriter) Write(p []byte) (int, error) {
	if w.err == nil {
		w.err = errors.New("write error")
	}
	return 0, w.err
}

type incompleteWriter struct {
	maxWrite int
	written  int
}

func (w *incompleteWriter) Write(p []byte) (int, error) {
	if w.written >= w.maxWrite {
		return 0, errors.New("write failed")
	}
	toWrite := len(p)
	if toWrite > w.maxWrite-w.written {
		toWrite = w.maxWrite - w.written
	}
	w.written += toWrite
	return toWrite, nil
}

// TestSignatureWriteToErrorPaths verifies WriteTo error handling
func TestSignatureWriteToErrorPaths(t *testing.T) {
	tests := []struct {
		name      string
		sig       Signature
		writer    io.Writer
		wantErr   bool
		errSubstr string
	}{
		{
			"Error writer during header write",
			Signature{
				SignatureType: SignatureTypeMLDSA,
				SignatureSize: 64,
				SignatureData: make([]byte, 64),
				CommentLength: 0,
			},
			&errorWriter{},
			true,
			"failed to write",
		},
		{
			"Failing writer during header write",
			Signature{
				SignatureType: SignatureTypeMLDSA,
				SignatureSize: 64,
				SignatureData: make([]byte, 64),
				CommentLength: 0,
			},
			&failingWriter{maxBytes: 10},
			true,
			"failed to write",
		},
		{
			"Failing writer during SignatureType write",
			Signature{
				SignatureType: SignatureTypeMLDSA,
				SignatureSize: 64,
				SignatureData: make([]byte, 64),
				CommentLength: 0,
			},
			&failingWriter{maxBytes: 3},
			true,
			"failed to write",
		},
		{
			"Failing writer during SignatureSize write",
			Signature{
				SignatureType: SignatureTypeMLDSA,
				SignatureSize: 64,
				SignatureData: make([]byte, 64),
				CommentLength: 0,
			},
			&failingWriter{maxBytes: 7},
			true,
			"failed to write",
		},
		{
			"Failing writer during SignatureFlags write",
			Signature{
				SignatureType: SignatureTypeMLDSA,
				SignatureSize: 64,
				SignatureData: make([]byte, 64),
				CommentLength: 0,
			},
			&failingWriter{maxBytes: 11},
			true,
			"failed to write",
		},
		{
			"Failing writer during SignatureTimestamp write",
			Signature{
				SignatureType: SignatureTypeMLDSA,
				SignatureSize: 64,
				SignatureData: make([]byte, 64),
				CommentLength: 0,
			},
			&failingWriter{maxBytes: 15},
			true,
			"failed to write",
		},
		{
			"Failing writer during CommentLength write",
			Signature{
				SignatureType: SignatureTypeMLDSA,
				SignatureSize: 64,
				SignatureData: make([]byte, 64),
				CommentLength: 0,
			},
			&failingWriter{maxBytes: 17},
			true,
			"failed to write",
		},
		{
			"Failing writer during comment write",
			Signature{
				SignatureType:    SignatureTypeMLDSA,
				SignatureSize:    64,
				SignatureData:    make([]byte, 64),
				CommentLength:    10,
				SignatureComment: "test comment",
			},
			&failingWriter{maxBytes: 17}, // Allow header (18 bytes) but fail during comment write
			true,
			"failed to write",
		},
		{
			"Incomplete comment write",
			Signature{
				SignatureType:    SignatureTypeMLDSA,
				SignatureSize:    64,
				SignatureData:    make([]byte, 64),
				CommentLength:    10,
				SignatureComment: "test comment",
			},
			&incompleteWriter{maxWrite: 20},
			true,
			"incomplete comment write",
		},
		{
			"Failing writer during signature data write",
			Signature{
				SignatureType: SignatureTypeMLDSA,
				SignatureSize: 64,
				SignatureData: make([]byte, 64),
				CommentLength: 0,
			},
			&failingWriter{maxBytes: 17}, // Allow header (18 bytes) but fail during data write
			true,
			"failed to write",
		},
		{
			"Incomplete signature data write",
			Signature{
				SignatureType: SignatureTypeMLDSA,
				SignatureSize: 64,
				SignatureData: make([]byte, 64),
				CommentLength: 0,
			},
			&incompleteWriter{maxWrite: 30},
			true,
			"incomplete signature data write",
		},
		{
			"Signature with comment and data - failing during data write",
			Signature{
				SignatureType:    SignatureTypeMLDSA,
				SignatureSize:    64,
				SignatureData:    make([]byte, 64),
				CommentLength:    10,
				SignatureComment: "test comment",
			},
			&failingWriter{maxBytes: 28}, // Allow header (18) + comment (10) but fail during data write
			true,
			"failed to write",
		},
		{
			"Signature with comment and data - incomplete data write",
			Signature{
				SignatureType:    SignatureTypeMLDSA,
				SignatureSize:    64,
				SignatureData:    make([]byte, 64),
				CommentLength:    10,
				SignatureComment: "test comment",
			},
			&incompleteWriter{maxWrite: 40}, // Allow header (18) + comment (10) + partial data (12)
			true,
			"incomplete signature data write",
		},
		{
			"Signature with large data - failing during write",
			Signature{
				SignatureType: SignatureTypeMLDSA,
				SignatureSize: 256,
				SignatureData: make([]byte, 256),
				CommentLength: 0,
			},
			&failingWriter{maxBytes: 50}, // Allow header (18) but fail during large data write
			true,
			"failed to write",
		},
		{
			"Signature with large data - incomplete write",
			Signature{
				SignatureType: SignatureTypeMLDSA,
				SignatureSize: 256,
				SignatureData: make([]byte, 256),
				CommentLength: 0,
			},
			&incompleteWriter{maxWrite: 100}, // Allow header (18) + partial data (82)
			true,
			"incomplete signature data write",
		},
		{
			"Signature with comment - incomplete comment write (exact boundary)",
			Signature{
				SignatureType:    SignatureTypeMLDSA,
				SignatureSize:    64,
				SignatureData:    make([]byte, 64),
				CommentLength:    10,
				SignatureComment: "test comment",
			},
			&incompleteWriter{maxWrite: 27}, // Allow header (18) + 9 bytes of comment (need 10)
			true,
			"incomplete comment write",
		},
		{
			"Signature with data - incomplete data write (exact boundary)",
			Signature{
				SignatureType: SignatureTypeMLDSA,
				SignatureSize: 64,
				SignatureData: make([]byte, 64),
				CommentLength: 0,
			},
			&incompleteWriter{maxWrite: 81}, // Allow header (18) + 63 bytes of data (need 64)
			true,
			"incomplete signature data write",
		},
		{
			"Signature with comment and data - incomplete comment write (exact boundary)",
			Signature{
				SignatureType:    SignatureTypeMLDSA,
				SignatureSize:    64,
				SignatureData:    make([]byte, 64),
				CommentLength:    5,
				SignatureComment: "hello",
			},
			&incompleteWriter{maxWrite: 22}, // Allow header (18) + 4 bytes of comment (need 5)
			true,
			"incomplete comment write",
		},
		{
			"Signature with comment and data - incomplete data write (exact boundary)",
			Signature{
				SignatureType:    SignatureTypeMLDSA,
				SignatureSize:    64,
				SignatureData:    make([]byte, 64),
				CommentLength:    5,
				SignatureComment: "hello",
			},
			&incompleteWriter{maxWrite: 86}, // Allow header (18) + comment (5) + 63 bytes of data (need 64)
			true,
			"incomplete signature data write",
		},
		{
			"Signature with very large data - failing during write",
			Signature{
				SignatureType: SignatureTypeMLDSA,
				SignatureSize: 1024,
				SignatureData: make([]byte, 1024),
				CommentLength: 0,
			},
			&failingWriter{maxBytes: 100}, // Allow header (18) but fail during large data write
			true,
			"failed to write",
		},
		{
			"Signature with very large comment - failing during write",
			Signature{
				SignatureType:    SignatureTypeMLDSA,
				SignatureSize:    64,
				SignatureData:    make([]byte, 64),
				CommentLength:    200,
				SignatureComment: strings.Repeat("A", 200),
			},
			&failingWriter{maxBytes: 50}, // Allow header (18) but fail during large comment write
			true,
			"failed to write",
		},
		{
			"Signature with very large comment - incomplete write",
			Signature{
				SignatureType:    SignatureTypeMLDSA,
				SignatureSize:    64,
				SignatureData:    make([]byte, 64),
				CommentLength:    200,
				SignatureComment: strings.Repeat("A", 200),
			},
			&incompleteWriter{maxWrite: 100}, // Allow header (18) + partial comment (82)
			true,
			"incomplete comment write",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// WriteTo updates lengths first
			_, err := tt.sig.WriteTo(tt.writer)

			if (err != nil) != tt.wantErr {
				t.Errorf("WriteTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil {
				if tt.errSubstr != "" {
					errStr := err.Error()
					if !strings.Contains(errStr, tt.errSubstr) {
						t.Errorf("WriteTo() error = %q, want error containing %q", errStr, tt.errSubstr)
					}
				}
			}
		})
	}
}

// Note: WriteTo updates CommentLength and SignatureSize to match actual data before writing,
// so length mismatch tests are not applicable here. The length checks in WriteTo are defensive
// and will always pass since the lengths are updated first.
