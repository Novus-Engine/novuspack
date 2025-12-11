package metadata

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

// TestPackageCommentValidation verifies validation logic
func TestPackageCommentValidation(t *testing.T) {
	tests := []struct {
		name    string
		comment PackageComment
		wantErr bool
	}{
		{
			"Valid comment",
			PackageComment{
				CommentLength: 5,
				Comment:       "Test\x00",
				Reserved:      [3]uint8{0, 0, 0},
			},
			false,
		},
		{
			"Empty comment",
			PackageComment{
				CommentLength: 0,
				Comment:       "",
				Reserved:      [3]uint8{0, 0, 0},
			},
			false,
		},
		{
			"Length mismatch",
			PackageComment{
				CommentLength: 0,
				Comment:       "Test",
			},
			true,
		},
		{
			"Exceeds max length",
			PackageComment{
				CommentLength: MaxCommentLength + 1,
				Comment:       "Test",
			},
			true,
		},
		{
			"Non-zero reserved",
			PackageComment{
				CommentLength: 4,
				Comment:       "Test",
				Reserved:      [3]uint8{1, 0, 0},
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.comment.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestPackageCommentSizeCalculation verifies size calculation
func TestPackageCommentSizeCalculation(t *testing.T) {
	tests := []struct {
		name          string
		commentLength uint32
		wantSize      int
	}{
		{"Empty comment", 0, 7},
		{"Short comment", 10, 17},
		{"Long comment", 1000, 1007},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comment := PackageComment{
				CommentLength: tt.commentLength,
			}

			if comment.Size() != tt.wantSize {
				t.Errorf("Size() = %d, want %d", comment.Size(), tt.wantSize)
			}
		})
	}
}

// TestPackageCommentIsEmpty verifies IsEmpty function
func TestPackageCommentIsEmpty(t *testing.T) {
	tests := []struct {
		name      string
		comment   PackageComment
		wantEmpty bool
	}{
		{"Empty comment", PackageComment{CommentLength: 0}, true},
		{"Non-empty comment", PackageComment{CommentLength: 5, Comment: "test\x00"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.comment.IsEmpty() != tt.wantEmpty {
				t.Errorf("IsEmpty() = %v, want %v", tt.comment.IsEmpty(), tt.wantEmpty)
			}
		})
	}
}

// TestPackageCommentSetComment verifies SetComment function
func TestPackageCommentSetComment(t *testing.T) {
	tests := []struct {
		name      string
		comment   string
		wantErr   bool
		wantLen   uint32
		checkNull bool
	}{
		{"Empty comment", "", false, 1, true},
		{"Simple comment", "test", false, 5, true},
		{"Comment with newline", "test\ncomment", false, 13, true},
		{"Comment with tab", "test\tcomment", false, 13, true},
		{"Comment already null-terminated", "test\x00", false, 5, true},
		{"Comment with embedded null", "test\x00middle", true, 0, false},
		{"Comment exceeding max length", strings.Repeat("a", MaxCommentLength), true, 0, false},
		{"Comment at max length boundary", strings.Repeat("a", MaxCommentLength-1), false, MaxCommentLength, true},
		{"Comment with invalid UTF-8", "\xFF\xFE", true, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pc PackageComment
			err := pc.SetComment(tt.comment)

			if (err != nil) != tt.wantErr {
				t.Errorf("SetComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if pc.CommentLength != tt.wantLen {
					t.Errorf("CommentLength = %d, want %d", pc.CommentLength, tt.wantLen)
				}

				if tt.checkNull {
					commentBytes := []byte(pc.Comment)
					if len(commentBytes) == 0 || commentBytes[len(commentBytes)-1] != 0x00 {
						t.Errorf("Comment is not null-terminated")
					}
				}

				// Reserved should be zero
				for i, b := range pc.Reserved {
					if b != 0 {
						t.Errorf("Reserved[%d] = %d, want 0", i, b)
					}
				}
			}
		})
	}
}

// TestPackageCommentGetComment verifies GetComment function
func TestPackageCommentGetComment(t *testing.T) {
	tests := []struct {
		name        string
		comment     PackageComment
		wantComment string
	}{
		{"Empty comment", PackageComment{CommentLength: 0, Comment: ""}, ""},
		{"Null-terminated comment", PackageComment{CommentLength: 5, Comment: "test\x00"}, "test"},
		{"Comment without null terminator", PackageComment{CommentLength: 4, Comment: "test"}, "test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.comment.GetComment()
			if got != tt.wantComment {
				t.Errorf("GetComment() = %q, want %q", got, tt.wantComment)
			}
		})
	}
}

// TestNewPackageComment verifies NewPackageComment function
func TestNewPackageComment(t *testing.T) {
	pc := NewPackageComment()

	if pc == nil {
		t.Fatal("NewPackageComment() returned nil")
	}

	if pc.CommentLength != 0 {
		t.Errorf("CommentLength = %d, want 0", pc.CommentLength)
	}
	if pc.Comment != "" {
		t.Errorf("Comment = %q, want empty", pc.Comment)
	}
	for i, b := range pc.Reserved {
		if b != 0 {
			t.Errorf("Reserved[%d] = %d, want 0", i, b)
		}
	}

	// Verify it's equivalent to an empty comment state
	if !pc.IsEmpty() {
		t.Errorf("IsEmpty() = false, want true for new PackageComment")
	}

	// Verify it passes validation
	if err := pc.Validate(); err != nil {
		t.Errorf("Validate() error = %v, want nil", err)
	}
}

// TestPackageCommentClear verifies Clear function
func TestPackageCommentClear(t *testing.T) {
	pc := PackageComment{
		CommentLength: 10,
		Comment:       "test\x00",
		Reserved:      [3]uint8{1, 2, 3},
	}

	pc.Clear()

	if pc.CommentLength != 0 {
		t.Errorf("CommentLength = %d, want 0", pc.CommentLength)
	}
	if pc.Comment != "" {
		t.Errorf("Comment = %q, want empty", pc.Comment)
	}
	for i, b := range pc.Reserved {
		if b != 0 {
			t.Errorf("Reserved[%d] = %d, want 0", i, b)
		}
	}
}

// TestPackageCommentReadFrom verifies ReadFrom function
func TestPackageCommentReadFrom(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		wantErr     bool
		wantLen     uint32
		wantComment string
	}{
		{
			"Empty comment",
			[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			false,
			0,
			"",
		},
		{
			"Simple comment",
			[]byte{
				0x05, 0x00, 0x00, 0x00, // CommentLength = 5
				0x74, 0x65, 0x73, 0x74, 0x00, // "test\x00"
				0x00, 0x00, 0x00, // Reserved
			},
			false,
			5,
			"test",
		},
		{
			"Comment with newline",
			[]byte{
				0x0D, 0x00, 0x00, 0x00, // CommentLength = 13
				0x74, 0x65, 0x73, 0x74, 0x0A, 0x63, 0x6F, 0x6D, 0x6D, 0x65, 0x6E, 0x74, 0x00, // "test\ncomment\x00"
				0x00, 0x00, 0x00, // Reserved
			},
			false,
			13,
			"test\ncomment",
		},
		{
			"Incomplete data",
			[]byte{0x05, 0x00, 0x00, 0x00},
			true,
			0,
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pc PackageComment
			r := bytes.NewReader(tt.data)
			_, err := pc.ReadFrom(r)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if pc.CommentLength != tt.wantLen {
					t.Errorf("CommentLength = %d, want %d", pc.CommentLength, tt.wantLen)
				}

				gotComment := pc.GetComment()
				if gotComment != tt.wantComment {
					t.Errorf("GetComment() = %q, want %q", gotComment, tt.wantComment)
				}
			}
		})
	}
}

// TestPackageCommentWriteTo verifies WriteTo function
func TestPackageCommentWriteTo(t *testing.T) {
	tests := []struct {
		name     string
		comment  PackageComment
		wantErr  bool
		wantSize int64
	}{
		{
			"Empty comment",
			PackageComment{CommentLength: 0, Comment: "", Reserved: [3]uint8{0, 0, 0}},
			false,
			7,
		},
		{
			"Simple comment",
			PackageComment{CommentLength: 5, Comment: "test\x00", Reserved: [3]uint8{0, 0, 0}},
			false,
			12,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			n, err := tt.comment.WriteTo(&buf)

			if (err != nil) != tt.wantErr {
				t.Errorf("WriteTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if n != tt.wantSize {
					t.Errorf("WriteTo() wrote %d bytes, want %d", n, tt.wantSize)
				}

				// Verify we can read it back
				var pc PackageComment
				r := bytes.NewReader(buf.Bytes())
				_, readErr := pc.ReadFrom(r)
				if readErr != nil {
					t.Errorf("Failed to read back written data: %v", readErr)
				}

				if pc.CommentLength != tt.comment.CommentLength {
					t.Errorf("Read back CommentLength = %d, want %d", pc.CommentLength, tt.comment.CommentLength)
				}
			}
		})
	}
}

// TestPackageCommentRoundTrip verifies round-trip serialization
func TestPackageCommentRoundTrip(t *testing.T) {
	tests := []struct {
		name    string
		comment string
	}{
		{"Empty", ""},
		{"Simple", "test"},
		{"With newline", "test\ncomment"},
		{"With tab", "test\tcomment"},
		{"Unicode", "test 测试 comment"},
		{"Long comment", strings.Repeat("a", 1000)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pc1 PackageComment
			if err := pc1.SetComment(tt.comment); err != nil {
				t.Fatalf("SetComment() error = %v", err)
			}

			var buf bytes.Buffer
			if _, err := pc1.WriteTo(&buf); err != nil {
				t.Fatalf("WriteTo() error = %v", err)
			}

			var pc2 PackageComment
			if _, err := pc2.ReadFrom(&buf); err != nil {
				t.Fatalf("ReadFrom() error = %v", err)
			}

			if pc1.CommentLength != pc2.CommentLength {
				t.Errorf("CommentLength mismatch: %d != %d", pc1.CommentLength, pc2.CommentLength)
			}

			if pc1.GetComment() != pc2.GetComment() {
				t.Errorf("Comment mismatch: %q != %q", pc1.GetComment(), pc2.GetComment())
			}

			// Validate the read comment
			if err := pc2.Validate(); err != nil {
				t.Errorf("Read comment validation failed: %v", err)
			}
		})
	}
}

// TestPackageCommentValidationEnhanced verifies enhanced validation
func TestPackageCommentValidationEnhanced(t *testing.T) {
	tests := []struct {
		name    string
		comment PackageComment
		wantErr bool
		errMsg  string
	}{
		{
			"Valid empty comment",
			PackageComment{CommentLength: 0, Comment: "", Reserved: [3]uint8{0, 0, 0}},
			false,
			"",
		},
		{
			"Valid comment",
			PackageComment{CommentLength: 5, Comment: "test\x00", Reserved: [3]uint8{0, 0, 0}},
			false,
			"",
		},
		{
			"Missing null terminator",
			PackageComment{CommentLength: 4, Comment: "test", Reserved: [3]uint8{0, 0, 0}},
			true,
			"not null-terminated",
		},
		{
			"Embedded null",
			PackageComment{CommentLength: 6, Comment: "te\x00st\x00", Reserved: [3]uint8{0, 0, 0}},
			true,
			"embedded null",
		},
		{
			"Length mismatch",
			PackageComment{CommentLength: 10, Comment: "test\x00", Reserved: [3]uint8{0, 0, 0}},
			true,
			"length mismatch",
		},
		{
			"Exceeds max length",
			PackageComment{CommentLength: MaxCommentLength + 1, Comment: "test\x00", Reserved: [3]uint8{0, 0, 0}},
			true,
			"exceeds maximum",
		},
		{
			"Non-zero reserved",
			PackageComment{CommentLength: 5, Comment: "test\x00", Reserved: [3]uint8{1, 0, 0}},
			true,
			"reserved byte",
		},
		{
			"Invalid UTF-8",
			PackageComment{CommentLength: 3, Comment: "\xFF\xFE\x00", Reserved: [3]uint8{0, 0, 0}},
			true,
			"not valid UTF-8",
		},
		{
			"Empty comment bytes with length > 0",
			PackageComment{CommentLength: 5, Comment: "", Reserved: [3]uint8{0, 0, 0}},
			true,
			"comment is empty but comment length is",
		},
		{
			"Comment length mismatch",
			PackageComment{CommentLength: 10, Comment: "test\x00", Reserved: [3]uint8{0, 0, 0}},
			true,
			"length mismatch",
		},
		{
			"Reserved byte at index 1 for empty comment",
			PackageComment{CommentLength: 0, Comment: "", Reserved: [3]uint8{0, 1, 0}},
			true,
			"reserved byte 1",
		},
		{
			"Reserved byte at index 2 for empty comment",
			PackageComment{CommentLength: 0, Comment: "", Reserved: [3]uint8{0, 0, 1}},
			true,
			"reserved byte 2",
		},
		{
			"Reserved byte at index 1 for non-empty comment",
			PackageComment{CommentLength: 5, Comment: "test\x00", Reserved: [3]uint8{0, 1, 0}},
			true,
			"reserved byte 1",
		},
		{
			"Reserved byte at index 2 for non-empty comment",
			PackageComment{CommentLength: 5, Comment: "test\x00", Reserved: [3]uint8{0, 0, 1}},
			true,
			"reserved byte 2",
		},
		{
			"CommentLength exactly at max boundary",
			PackageComment{CommentLength: MaxCommentLength, Comment: strings.Repeat("a", MaxCommentLength-1) + "\x00", Reserved: [3]uint8{0, 0, 0}},
			false,
			"",
		},
		{
			"CommentLength equals 1 (just null terminator)",
			PackageComment{CommentLength: 1, Comment: "\x00", Reserved: [3]uint8{0, 0, 0}},
			false,
			"",
		},
		{
			"Embedded null at position 0",
			PackageComment{CommentLength: 6, Comment: "\x00test\x00", Reserved: [3]uint8{0, 0, 0}},
			true,
			"embedded null character at position 0",
		},
		{
			"Embedded null at last position before terminator",
			PackageComment{CommentLength: 6, Comment: "test\x00\x00", Reserved: [3]uint8{0, 0, 0}},
			true,
			"embedded null character at position 4",
		},
		{
			"CommentLength mismatch - shorter than actual",
			PackageComment{CommentLength: 3, Comment: "test\x00", Reserved: [3]uint8{0, 0, 0}},
			true,
			"length mismatch",
		},
		{
			"All reserved bytes non-zero for empty comment",
			PackageComment{CommentLength: 0, Comment: "", Reserved: [3]uint8{1, 2, 3}},
			true,
			"reserved byte",
		},
		{
			"All reserved bytes non-zero for non-empty comment",
			PackageComment{CommentLength: 5, Comment: "test\x00", Reserved: [3]uint8{1, 2, 3}},
			true,
			"reserved byte",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.comment.Validate()

			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil {
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Validate() error = %q, want error containing %q", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

// TestPackageCommentReadFromIncompleteData tests error handling for incomplete data
func TestPackageCommentReadFromIncompleteData(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"No data", []byte{}},
		{"Incomplete length", []byte{0x05}},
		{"Incomplete comment", []byte{0x05, 0x00, 0x00, 0x00, 0x74}},
		{"Incomplete reserved", []byte{0x05, 0x00, 0x00, 0x00, 0x74, 0x65, 0x73, 0x74, 0x00, 0x00}},
		{"Incomplete reserved bytes read", []byte{0x05, 0x00, 0x00, 0x00, 0x74, 0x65, 0x73, 0x74, 0x00, 0x00, 0x00}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pc PackageComment
			r := bytes.NewReader(tt.data)
			_, err := pc.ReadFrom(r)

			if err == nil {
				t.Errorf("ReadFrom() expected error for incomplete data, got nil")
			}
		})
	}
}

// TestPackageCommentReadFromIncompleteReservedBytes tests incomplete reserved bytes read
func TestPackageCommentReadFromIncompleteReservedBytes(t *testing.T) {
	// Test case where comment is read successfully but reserved bytes are incomplete
	data := []byte{
		0x00, 0x00, 0x00, 0x00, // CommentLength = 0 (empty comment)
		0x00, 0x00, // Only 2 bytes of reserved (need 3)
	}

	var pc PackageComment
	r := bytes.NewReader(data)
	_, err := pc.ReadFrom(r)

	if err == nil {
		t.Errorf("ReadFrom() expected error for incomplete reserved bytes, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "failed to read reserved bytes") {
		t.Errorf("ReadFrom() error = %q, want error containing 'failed to read reserved bytes'", err.Error())
	}
}

// TestPackageCommentWriteToWriterError tests error handling
func TestPackageCommentWriteToWriterError(t *testing.T) {
	tests := []struct {
		name    string
		comment PackageComment
		writer  io.Writer
	}{
		{
			"Error writing comment length",
			PackageComment{CommentLength: 5, Comment: "test\x00", Reserved: [3]uint8{0, 0, 0}},
			&failingWriter{maxWrite: 0},
		},
		{
			"Error writing comment (partial write)",
			PackageComment{CommentLength: 5, Comment: "test\x00", Reserved: [3]uint8{0, 0, 0}},
			&failingWriter{maxWrite: 4},
		},
		{
			"Error writing reserved bytes (partial write)",
			PackageComment{CommentLength: 5, Comment: "test\x00", Reserved: [3]uint8{0, 0, 0}},
			&failingWriter{maxWrite: 9},
		},
		{
			"Error writing reserved bytes",
			PackageComment{CommentLength: 0, Comment: "", Reserved: [3]uint8{0, 0, 0}},
			&failingWriter{maxWrite: 4},
		},
		{
			"Error writing empty comment reserved",
			PackageComment{CommentLength: 0, Comment: "", Reserved: [3]uint8{0, 0, 0}},
			&failingWriter{maxWrite: 6},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.comment.WriteTo(tt.writer)
			if err == nil {
				t.Errorf("WriteTo() expected error from failing writer, got nil")
			}
		})
	}
}

// TestPackageCommentWriteToIncompleteWrite tests incomplete write scenarios
func TestPackageCommentWriteToIncompleteWrite(t *testing.T) {
	// Test incomplete comment write
	comment := PackageComment{CommentLength: 5, Comment: "test\x00", Reserved: [3]uint8{0, 0, 0}}
	w := &partialWriter{maxWrite: 8} // Can write length (4) but not full comment (5)
	_, err := comment.WriteTo(w)
	if err == nil {
		t.Errorf("WriteTo() expected error for incomplete comment write, got nil")
	}
}

// failingWriter is a writer that fails after writing maxWrite bytes
type failingWriter struct {
	maxWrite int
	written  int
}

func (w *failingWriter) Write(p []byte) (int, error) {
	if w.written >= w.maxWrite {
		return 0, io.ErrShortWrite
	}
	n := len(p)
	if w.written+n > w.maxWrite {
		n = w.maxWrite - w.written
	}
	w.written += n
	return n, nil
}

// partialWriter is a writer that writes partial data without error
// This simulates a writer that successfully writes but doesn't write all bytes
type partialWriter struct {
	maxWrite int
	written  int
}

func (w *partialWriter) Write(p []byte) (int, error) {
	if w.written >= w.maxWrite {
		return 0, nil // Return 0 bytes written, no error (triggers incomplete write check)
	}
	n := len(p)
	if w.written+n > w.maxWrite {
		n = w.maxWrite - w.written
	}
	w.written += n
	// Return partial write without error - this will trigger the incomplete write check
	return n, nil
}
