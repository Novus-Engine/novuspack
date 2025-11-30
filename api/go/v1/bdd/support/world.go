package support

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Package represents a NovusPack package instance
// This will be replaced with the actual API type once implemented
type Package interface {
	Close() error
	IsOpen() bool
}

// World provides test context and state management for BDD scenarios
type World struct {
	WorkDir string
	TempDir string

	// Package state management
	packages map[string]Package // path -> package instance
	mu       sync.RWMutex

	// Current package context (for scenarios that work with a single package)
	currentPackage Package
	currentPath    string

	// Error tracking
	lastError error

	// File state tracking
	files map[string][]byte // path -> content

	// Context storage for timeout and cancellation testing
	testContext context.Context
	testCancel  context.CancelFunc

	// Package metadata storage for verification
	packageMetadata map[string]interface{} // key -> value
}

// NewWorld creates a new test world instance
func NewWorld() (*World, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	tempDir, err := os.MkdirTemp("", "novuspack-bdd-")
	if err != nil {
		return nil, err
	}

	return &World{
		WorkDir:        workDir,
		TempDir:        tempDir,
		packages:       make(map[string]Package),
		files:          make(map[string][]byte),
		packageMetadata: make(map[string]interface{}),
	}, nil
}

// Resolve resolves a path relative to the work directory
func (w *World) Resolve(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(w.WorkDir, path)
}

// TempPath returns a path in the temporary directory
func (w *World) TempPath(name string) string {
	return filepath.Join(w.TempDir, name)
}

// SetPackage sets the current package for the scenario
func (w *World) SetPackage(pkg Package, path string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.currentPackage = pkg
	w.currentPath = path
	if path != "" {
		w.packages[path] = pkg
	}
}

// GetPackage returns the current package
func (w *World) GetPackage() Package {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.currentPackage
}

// GetPackageByPath returns a package by its path
func (w *World) GetPackageByPath(path string) Package {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.packages[path]
}

// SetError records the last error that occurred
func (w *World) SetError(err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.lastError = err
}

// GetError returns the last recorded error
func (w *World) GetError() error {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.lastError
}

// ClearError clears the last recorded error
func (w *World) ClearError() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.lastError = nil
}

// SetFileContent stores file content for later retrieval
func (w *World) SetFileContent(path string, content []byte) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.files[path] = content
}

// GetFileContent retrieves stored file content
func (w *World) GetFileContent(path string) []byte {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.files[path]
}

// HasFile checks if file content is stored
func (w *World) HasFile(path string) bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	_, exists := w.files[path]
	return exists
}

// NewContext creates a new context for package operations
func (w *World) NewContext() context.Context {
	return context.Background()
}

// NewContextWithTimeout creates a context with timeout
func (w *World) NewContextWithTimeout(timeout string) (context.Context, context.CancelFunc) {
	// TODO: Parse timeout string and create context with timeout
	// For now, use default timeout
	duration := 30 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	w.mu.Lock()
	w.testContext = ctx
	w.testCancel = cancel
	w.mu.Unlock()
	return ctx, cancel
}

// GetTestContext returns the stored test context
func (w *World) GetTestContext() context.Context {
	w.mu.RLock()
	defer w.mu.RUnlock()
	if w.testContext != nil {
		return w.testContext
	}
	return context.Background()
}

// Cleanup closes all packages and removes temporary files
func (w *World) Cleanup() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Close all packages
	for path, pkg := range w.packages {
		if pkg != nil {
			_ = pkg.Close()
		}
		delete(w.packages, path)
	}

	w.currentPackage = nil
	w.currentPath = ""
	w.lastError = nil
	w.files = make(map[string][]byte)
	w.packageMetadata = make(map[string]interface{})
	w.testContext = nil
	w.testCancel = nil

	// Remove temporary directory
	if w.TempDir != "" {
		return os.RemoveAll(w.TempDir)
	}
	return nil
}

// SetPackageMetadata stores package metadata for later verification
func (w *World) SetPackageMetadata(key string, value interface{}) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.packageMetadata == nil {
		w.packageMetadata = make(map[string]interface{})
	}
	w.packageMetadata[key] = value
}

// GetPackageMetadata retrieves stored package metadata
func (w *World) GetPackageMetadata(key string) (interface{}, bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	if w.packageMetadata == nil {
		return nil, false
	}
	value, exists := w.packageMetadata[key]
	return value, exists
}
