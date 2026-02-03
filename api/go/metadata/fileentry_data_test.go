package metadata

import (
	"bytes"
	"context"
	"os"
	"testing"
)

// TestLoadData tests LoadData method
//
//nolint:gocognit,gocyclo // table-driven load cases
func TestLoadData(t *testing.T) {
	// Create a temporary source file
	sourceFile, err := os.CreateTemp("", "novuspack-source-*")
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}
	defer func() {
		_ = sourceFile.Close()
		_ = os.Remove(sourceFile.Name())
	}()

	testData := []byte("test data for loading")
	if _, err := sourceFile.Write(testData); err != nil {
		t.Fatalf("Failed to write to source file: %v", err)
	}

	tests := []struct {
		name    string
		setup   func() *FileEntry
		wantErr bool
	}{
		{
			name:    "no source file",
			setup:   NewFileEntry,
			wantErr: true,
		},
		{
			name: "already loaded",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				fe.SetData([]byte("already loaded"))
				return fe
			},
			wantErr: false,
		},
		{
			name: "successful load from source",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				fe.setSourceFile(sourceFile, 0, int64(len(testData)))
				return fe
			},
			wantErr: false,
		},
		{
			name: "cancelled context",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				fe.setSourceFile(sourceFile, 0, int64(len(testData)))
				return fe
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fe := tt.setup()
			ctx := context.Background()
			if tt.name == "cancelled context" {
				cancelledCtx, cancel := context.WithCancel(context.Background())
				cancel()
				ctx = cancelledCtx
			}

			err := fe.LoadData(ctx)

			if (err != nil) != tt.wantErr {
				t.Errorf("LoadData() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && tt.name == "successful load from source" {
				if !fe.IsDataLoaded {
					t.Error("LoadData() did not set IsDataLoaded flag")
				}

				if len(fe.Data) != len(testData) {
					t.Errorf("LoadData() Data length = %d, want %d", len(fe.Data), len(testData))
				}

				if fe.ProcessingState != ProcessingStateComplete {
					t.Errorf("LoadData() ProcessingState = %v, want ProcessingStateComplete", fe.ProcessingState)
				}
			}
		})
	}

	// Test incomplete read (EOF but not enough bytes)
	t.Run("incomplete read", func(t *testing.T) {
		fe := NewFileEntry()
		// Set source size larger than actual file content
		fe.setSourceFile(sourceFile, 0, int64(len(testData)+10))
		err := fe.LoadData(context.Background())
		if err == nil {
			t.Error("LoadData() with incomplete read should return error")
		}
		if fe.ProcessingState != ProcessingStateError {
			t.Errorf("LoadData() ProcessingState = %v, want ProcessingStateError", fe.ProcessingState)
		}
	})

	// Test read error (not EOF)
	t.Run("read error", func(t *testing.T) {
		fe := NewFileEntry()
		// Create a file that will cause read error
		badFile, err := os.Open(os.DevNull)
		if err != nil {
			t.Fatalf("Failed to open /dev/null: %v", err)
		}
		defer func() { _ = badFile.Close() }()
		fe.setSourceFile(badFile, 0, 100)
		err = fe.LoadData(context.Background())
		if err == nil {
			t.Error("LoadData() with read error should return error")
		}
		if fe.ProcessingState != ProcessingStateError {
			t.Errorf("LoadData() ProcessingState = %v, want ProcessingStateError", fe.ProcessingState)
		}
	})
}

// TestUnloadData tests UnloadData method
func TestUnloadData(t *testing.T) {
	// Test normal context
	fe := NewFileEntry()
	fe.SetData([]byte("test data"))

	fe.UnloadData()

	if fe.IsDataLoaded {
		t.Error("UnloadData() did not clear IsDataLoaded flag")
	}

	if len(fe.Data) != 0 {
		t.Error("UnloadData() did not clear Data field")
	}

	if fe.ProcessingState != ProcessingStateIdle {
		t.Errorf("UnloadData() ProcessingState = %v, want ProcessingStateIdle", fe.ProcessingState)
	}
}

// TestGetData tests GetData method
//
//nolint:gocognit // table-driven get cases
func TestGetData(t *testing.T) {
	// Create a temporary source file
	sourceFile, err := os.CreateTemp("", "novuspack-source-*")
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}
	defer func() {
		_ = sourceFile.Close()
		_ = os.Remove(sourceFile.Name())
	}()

	testData := []byte("test data for GetData")
	if _, err := sourceFile.Write(testData); err != nil {
		t.Fatalf("Failed to write to source file: %v", err)
	}

	tests := []struct {
		name    string
		setup   func() *FileEntry
		wantErr bool
	}{
		{
			name: "data already loaded",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				fe.SetData([]byte("test data"))
				return fe
			},
			wantErr: false,
		},
		{
			name:    "no data available",
			setup:   NewFileEntry,
			wantErr: true,
		},
		{
			name: "load from source file",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				fe.setSourceFile(sourceFile, 0, int64(len(testData)))
				return fe
			},
			wantErr: false,
		},
		{
			name: "load from temp file via ReadFromTempFile",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				err := fe.WriteToTempFile(context.Background(), testData)
				if err != nil {
					t.Fatalf("WriteToTempFile() error = %v", err)
				}
				// Read from temp file to load data
				data, err := fe.ReadFromTempFile(context.Background(), 0, int64(len(testData)))
				if err != nil {
					t.Fatalf("ReadFromTempFile() error = %v", err)
				}
				fe.SetData(data)
				return fe
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fe := tt.setup()
			data, err := fe.GetData()

			if (err != nil) != tt.wantErr {
				t.Errorf("GetData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(data) == 0 {
					t.Error("GetData() returned empty data")
				}

				if tt.name == "load from source file" && !bytes.Equal(data, testData) {
					t.Errorf("GetData() data = %q, want %q", string(data), string(testData))
				}
			}

			// Cleanup temp file if created
			if tt.name == "load from temp file" {
				_ = fe.CleanupTempFile(context.Background()) //nolint:errcheck // cleanup best-effort in test
			}
		})
	}
}

// TestSetData tests SetData method
func TestSetData(t *testing.T) {
	tests := []struct {
		name      string
		data      []byte
		wantState ProcessingState
	}{
		{
			name:      "with data",
			data:      []byte("test data"),
			wantState: ProcessingStateComplete,
		},
		{
			name:      "empty data",
			data:      []byte{},
			wantState: ProcessingStateComplete,
		},
		{
			name:      "nil data",
			data:      nil,
			wantState: ProcessingStateComplete,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fe := NewFileEntry()
			fe.SetData(tt.data)
			checkSetDataResult(t, fe, tt.data, tt.wantState)
		})
	}
}

// checkSetDataResult verifies FileEntry state after SetData; used to reduce TestSetData complexity.
func checkSetDataResult(t *testing.T, fe *FileEntry, data []byte, wantState ProcessingState) {
	t.Helper()
	if !fe.IsDataLoaded {
		t.Error("SetData() did not set IsDataLoaded flag")
	}

	if data == nil {
		data = []byte{}
	}
	if len(fe.Data) != len(data) {
		t.Errorf("SetData() Data length = %d, want %d", len(fe.Data), len(data))
	}
	if fe.ProcessingState != wantState {
		t.Errorf("SetData() ProcessingState = %v, want %v", fe.ProcessingState, wantState)
	}
}

// TestGetProcessingState tests GetProcessingState method
func TestGetProcessingState(t *testing.T) {
	fe := NewFileEntry()

	if fe.GetProcessingState() != ProcessingStateIdle {
		t.Errorf("GetProcessingState() = %v, want ProcessingStateIdle", fe.GetProcessingState())
	}

	fe.SetProcessingState(ProcessingStateLoading)
	if fe.GetProcessingState() != ProcessingStateLoading {
		t.Errorf("GetProcessingState() = %v, want ProcessingStateLoading", fe.GetProcessingState())
	}
}

// TestSetProcessingState tests SetProcessingState method
func TestSetProcessingState(t *testing.T) {
	fe := NewFileEntry()

	fe.SetProcessingState(ProcessingStateProcessing)
	if fe.ProcessingState != ProcessingStateProcessing {
		t.Errorf("SetProcessingState() ProcessingState = %v, want ProcessingStateProcessing", fe.ProcessingState)
	}
}

// TestSetSourceFile tests the internal setSourceFile helper.
func TestSetSourceFile(t *testing.T) {
	fe := NewFileEntry()
	testFile, _ := os.Open(os.DevNull)
	defer func() { _ = testFile.Close() }() //nolint:errcheck // Close on exit - error is non-critical

	fe.setSourceFile(testFile, 10, 100)

	if fe.SourceFile != testFile {
		t.Error("setSourceFile() did not set SourceFile")
	}

	if fe.SourceOffset != 10 {
		t.Errorf("setSourceFile() SourceOffset = %d, want 10", fe.SourceOffset)
	}

	if fe.SourceSize != 100 {
		t.Errorf("setSourceFile() SourceSize = %d, want 100", fe.SourceSize)
	}
}

// TestGetSourceFile tests the internal getSourceFile helper.
func TestGetSourceFile(t *testing.T) {
	fe := NewFileEntry()
	testFile, _ := os.Open(os.DevNull)
	defer func() { _ = testFile.Close() }() //nolint:errcheck // Close on exit - error is non-critical

	fe.setSourceFile(testFile, 10, 100)

	file, offset, size := fe.getSourceFile()
	if file != testFile {
		t.Error("getSourceFile() returned wrong file")
	}

	if offset != 10 {
		t.Errorf("getSourceFile() offset = %d, want 10", offset)
	}

	if size != 100 {
		t.Errorf("getSourceFile() size = %d, want 100", size)
	}
}

// TestSetTempPath tests the internal setTempPath helper.
func TestSetTempPath(t *testing.T) {
	fe := NewFileEntry()
	testPath := "/tmp/test"

	fe.setTempPath(testPath)

	if fe.TempFilePath != testPath {
		t.Errorf("setTempPath() TempFilePath = %q, want %q", fe.TempFilePath, testPath)
	}

	if !fe.IsTempFile {
		t.Error("setTempPath() did not set IsTempFile flag")
	}

	fe.setTempPath("")
	if fe.IsTempFile {
		t.Error("setTempPath() with empty path did not clear IsTempFile flag")
	}
}

// TestGetTempPath tests the internal getTempPath helper.
func TestGetTempPath(t *testing.T) {
	fe := NewFileEntry()
	testPath := "/tmp/test"

	fe.setTempPath(testPath)

	if fe.getTempPath() != testPath {
		t.Errorf("getTempPath() = %q, want %q", fe.getTempPath(), testPath)
	}
}

// TestCreateTempFile tests CreateTempFile method
func TestCreateTempFile(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		wantErr bool
	}{
		{
			name:    "successful creation",
			ctx:     context.Background(),
			wantErr: false,
		},
		{
			name: "cancelled context",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fe := NewFileEntry()

			err := fe.CreateTempFile(tt.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTempFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if fe.TempFilePath == "" {
					t.Error("CreateTempFile() did not set TempFilePath")
				}

				if !fe.IsTempFile {
					t.Error("CreateTempFile() did not set IsTempFile flag")
				}

				// Cleanup
				_ = fe.CleanupTempFile(context.Background()) //nolint:errcheck // Cleanup error is non-critical in test
			}
		})
	}
}

// TestCleanupTempFile tests CleanupTempFile method
//
//nolint:gocognit // table-driven cleanup cases
func TestCleanupTempFile(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *FileEntry
		wantErr bool
	}{
		{
			name: "cleanup existing temp file",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				err := fe.CreateTempFile(context.Background())
				if err != nil {
					t.Fatalf("CreateTempFile() error = %v", err)
				}
				return fe
			},
			wantErr: false,
		},
		{
			name:    "cleanup with no temp file",
			setup:   NewFileEntry,
			wantErr: false, // Should not error if no temp file
		},
		{
			name: "cleanup with cancelled context",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				err := fe.CreateTempFile(context.Background())
				if err != nil {
					t.Fatalf("CreateTempFile() error = %v", err)
				}
				return fe
			},
			wantErr: true, // Context cancellation returns error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fe := tt.setup()
			tempPath := fe.TempFilePath

			ctx := context.Background()
			if tt.name == "cleanup with cancelled context" {
				cancelledCtx, cancel := context.WithCancel(context.Background())
				cancel()
				ctx = cancelledCtx
			}

			err := fe.CleanupTempFile(ctx)

			if (err != nil) != tt.wantErr {
				t.Errorf("CleanupTempFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.name == "cleanup existing temp file" {
				if fe.TempFilePath != "" {
					t.Error("CleanupTempFile() did not clear TempFilePath")
				}

				if fe.IsTempFile {
					t.Error("CleanupTempFile() did not clear IsTempFile flag")
				}

				// Verify file is deleted
				if _, err := os.Stat(tempPath); err == nil {
					t.Error("CleanupTempFile() did not delete temp file")
				}
			}
		})
	}
}

// TestStreamToTempFile tests StreamToTempFile method
func TestStreamToTempFile(t *testing.T) {
	// Create a temporary source file
	sourceFile, err := os.CreateTemp("", "novuspack-source-*")
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}
	defer func() {
		_ = sourceFile.Close()
		_ = os.Remove(sourceFile.Name())
	}()

	testData := []byte("test data for streaming")
	if _, err := sourceFile.Write(testData); err != nil {
		t.Fatalf("Failed to write to source file: %v", err)
	}

	fe := NewFileEntry()
	fe.setSourceFile(sourceFile, 0, int64(len(testData)))

	// Test successful streaming
	err = fe.StreamToTempFile(context.Background())
	if err != nil {
		t.Fatalf("StreamToTempFile() error = %v", err)
	}

	if fe.TempFilePath == "" {
		t.Error("StreamToTempFile() did not set TempFilePath")
	}

	if !fe.IsTempFile {
		t.Error("StreamToTempFile() did not set IsTempFile flag")
	}

	if fe.ProcessingState != ProcessingStateComplete {
		t.Errorf("StreamToTempFile() ProcessingState = %v, want ProcessingStateComplete", fe.ProcessingState)
	}

	// Verify data was written
	readData, err := os.ReadFile(fe.TempFilePath)
	if err != nil {
		t.Fatalf("Failed to read temp file: %v", err)
	}

	if !bytes.Equal(readData, testData) {
		t.Errorf("StreamToTempFile() data = %q, want %q", string(readData), string(testData))
	}

	// Cleanup
	_ = fe.CleanupTempFile(context.Background()) //nolint:errcheck // cleanup best-effort in test

	// Test with no source file
	testStreamToTempFileNoSource(t)

	// Test with cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	fe3 := NewFileEntry()
	fe3.setSourceFile(sourceFile, 0, int64(len(testData)))
	err = fe3.StreamToTempFile(ctx)
	if err == nil {
		t.Error("StreamToTempFile() with cancelled context should return error")
	}

	// Test with existing temp file path (should reuse)
	fe4 := NewFileEntry()
	err = fe4.CreateTempFile(context.Background())
	if err != nil {
		t.Fatalf("CreateTempFile() error = %v", err)
	}
	tempPath := fe4.TempFilePath
	fe4.setSourceFile(sourceFile, 0, int64(len(testData)))
	err = fe4.StreamToTempFile(context.Background())
	if err != nil {
		t.Fatalf("StreamToTempFile() with existing temp file error = %v", err)
	}
	if fe4.TempFilePath != tempPath {
		t.Error("StreamToTempFile() should reuse existing temp file path")
	}
	_ = fe4.CleanupTempFile(context.Background()) //nolint:errcheck // cleanup best-effort in test

	// Test seek error
	fe5 := NewFileEntry()
	// Use a closed file to cause seek error
	closedFile, _ := os.Open(os.DevNull)
	_ = closedFile.Close()
	fe5.setSourceFile(closedFile, 0, 10)
	err = fe5.StreamToTempFile(context.Background())
	if err == nil {
		t.Error("StreamToTempFile() with seek error should return error")
	}

	// Test copy error (incomplete copy)
	testStreamToTempFileIncompleteCopy(t, sourceFile, testData)
}

// testStreamToTempFileNoSource verifies StreamToTempFile returns error when no source file is set.
func testStreamToTempFileNoSource(t *testing.T) {
	t.Helper()
	fe := NewFileEntry()
	err := fe.StreamToTempFile(context.Background())
	if err == nil {
		t.Error("StreamToTempFile() with no source file should return error")
	}
}

// testStreamToTempFileIncompleteCopy verifies StreamToTempFile returns error when copy is incomplete.
func testStreamToTempFileIncompleteCopy(t *testing.T, sourceFile *os.File, testData []byte) {
	t.Helper()
	fe := NewFileEntry()
	fe.setSourceFile(sourceFile, 0, int64(len(testData)+100)) // Request more than available
	err := fe.StreamToTempFile(context.Background())
	if err == nil {
		t.Error("StreamToTempFile() with incomplete copy should return error")
	}
}

// TestWriteToTempFile tests WriteToTempFile method
//
//nolint:gocognit,gocyclo // table-driven write cases
func TestWriteToTempFile(t *testing.T) {
	fe := NewFileEntry()
	testData := []byte("test data to write")

	// Test successful write
	err := fe.WriteToTempFile(context.Background(), testData)
	if err != nil {
		t.Fatalf("WriteToTempFile() error = %v", err)
	}

	if fe.TempFilePath == "" {
		t.Error("WriteToTempFile() did not set TempFilePath")
	}

	if !fe.IsTempFile {
		t.Error("WriteToTempFile() did not set IsTempFile flag")
	}

	if fe.ProcessingState != ProcessingStateComplete {
		t.Errorf("WriteToTempFile() ProcessingState = %v, want ProcessingStateComplete", fe.ProcessingState)
	}

	// Verify data was written
	readData, err := os.ReadFile(fe.TempFilePath)
	if err != nil {
		t.Fatalf("Failed to read temp file: %v", err)
	}

	if !bytes.Equal(readData, testData) {
		t.Errorf("WriteToTempFile() data = %q, want %q", string(readData), string(testData))
	}

	// Cleanup
	_ = fe.CleanupTempFile(context.Background()) //nolint:errcheck // cleanup best-effort in test

	// Test with cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	fe2 := NewFileEntry()
	err = fe2.WriteToTempFile(ctx, testData)
	if err == nil {
		t.Error("WriteToTempFile() with cancelled context should return error")
	}

	// Test with existing temp file path (should reuse)
	fe3 := NewFileEntry()
	err = fe3.CreateTempFile(context.Background())
	if err != nil {
		t.Fatalf("CreateTempFile() error = %v", err)
	}
	tempPath := fe3.TempFilePath
	err = fe3.WriteToTempFile(context.Background(), testData)
	if err != nil {
		t.Fatalf("WriteToTempFile() with existing temp file error = %v", err)
	}
	if fe3.TempFilePath != tempPath {
		t.Error("WriteToTempFile() should reuse existing temp file path")
	}
	_ = fe3.CleanupTempFile(context.Background()) //nolint:errcheck // cleanup best-effort in test

	// Test write error (simulate by using invalid temp file path after creation)
	fe4 := NewFileEntry()
	err = fe4.CreateTempFile(context.Background())
	if err != nil {
		t.Fatalf("CreateTempFile() error = %v", err)
	}
	// Remove the file to cause write error
	_ = os.Remove(fe4.TempFilePath)
	// Try to write - should fail when opening file
	err = fe4.WriteToTempFile(context.Background(), testData)
	// This might succeed if it recreates the file, or fail if it tries to open
	// Accept either outcome, but verify state
	if err == nil {
		if fe4.ProcessingState != ProcessingStateComplete {
			t.Errorf("WriteToTempFile() ProcessingState = %v, want ProcessingStateComplete", fe4.ProcessingState)
		}
	} else {
		if fe4.ProcessingState != ProcessingStateError {
			t.Errorf("WriteToTempFile() ProcessingState = %v, want ProcessingStateError", fe4.ProcessingState)
		}
	}
	_ = fe4.CleanupTempFile(context.Background()) //nolint:errcheck // cleanup best-effort in test

	// Test file open error (use invalid directory path)
	fe5 := NewFileEntry()
	fe5.TempFilePath = "/nonexistent/directory/file.txt" // Invalid path
	err = fe5.WriteToTempFile(context.Background(), testData)
	if err == nil {
		t.Error("WriteToTempFile() with invalid path should return error")
	}
	// Note: When file open fails, ProcessingState is not set to Error (only set on write failure)
	// So state remains at initial value (ProcessingStateIdle)
	if fe5.ProcessingState != ProcessingStateIdle {
		t.Errorf("WriteToTempFile() ProcessingState = %v, want ProcessingStateIdle (file open error doesn't set state)", fe5.ProcessingState)
	}
}

// TestReadFromTempFile tests ReadFromTempFile method
func TestReadFromTempFile(t *testing.T) {
	fe := NewFileEntry()
	testData := []byte("test data for reading")

	// Create temp file and write data
	err := fe.WriteToTempFile(context.Background(), testData)
	if err != nil {
		t.Fatalf("WriteToTempFile() error = %v", err)
	}

	// Test successful read
	readData, err := fe.ReadFromTempFile(context.Background(), 0, int64(len(testData)))
	if err != nil {
		t.Fatalf("ReadFromTempFile() error = %v", err)
	}

	if !bytes.Equal(readData, testData) {
		t.Errorf("ReadFromTempFile() data = %q, want %q", string(readData), string(testData))
	}

	// Test partial read
	partialData, err := fe.ReadFromTempFile(context.Background(), 0, 4)
	if err != nil {
		t.Fatalf("ReadFromTempFile() partial error = %v", err)
	}

	if string(partialData) != "test" {
		t.Errorf("ReadFromTempFile() partial data = %q, want %q", string(partialData), "test")
	}

	// Test read from offset
	offsetData, err := fe.ReadFromTempFile(context.Background(), 5, 4)
	if err != nil {
		t.Fatalf("ReadFromTempFile() offset error = %v", err)
	}

	if string(offsetData) != "data" {
		t.Errorf("ReadFromTempFile() offset data = %q, want %q", string(offsetData), "data")
	}

	// Cleanup
	_ = fe.CleanupTempFile(context.Background()) //nolint:errcheck // cleanup best-effort in test

	// Test with no temp file
	fe2 := NewFileEntry()
	_, err = fe2.ReadFromTempFile(context.Background(), 0, 10)
	if err == nil {
		t.Error("ReadFromTempFile() with no temp file should return error")
	}

	// Test with cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	fe3 := NewFileEntry()
	err = fe3.WriteToTempFile(context.Background(), testData)
	if err != nil {
		t.Fatalf("WriteToTempFile() error = %v", err)
	}
	_, err = fe3.ReadFromTempFile(ctx, 0, 10)
	if err == nil {
		t.Error("ReadFromTempFile() with cancelled context should return error")
	}
	_ = fe3.CleanupTempFile(context.Background()) //nolint:errcheck // cleanup best-effort in test

	// Test read beyond file size
	fe4 := NewFileEntry()
	err = fe4.WriteToTempFile(context.Background(), testData)
	if err != nil {
		t.Fatalf("WriteToTempFile() error = %v", err)
	}
	_, err = fe4.ReadFromTempFile(context.Background(), 0, int64(len(testData)+100))
	if err == nil {
		t.Error("ReadFromTempFile() beyond file size should return error")
	}
	_ = fe4.CleanupTempFile(context.Background()) //nolint:errcheck // cleanup best-effort in test

	// Test invalid offset
	fe5 := NewFileEntry()
	err = fe5.WriteToTempFile(context.Background(), testData)
	if err != nil {
		t.Fatalf("WriteToTempFile() error = %v", err)
	}
	_, err = fe5.ReadFromTempFile(context.Background(), -1, 10)
	if err == nil {
		t.Error("ReadFromTempFile() with invalid offset should return error")
	}
	_ = fe5.CleanupTempFile(context.Background()) //nolint:errcheck // cleanup best-effort in test
}
