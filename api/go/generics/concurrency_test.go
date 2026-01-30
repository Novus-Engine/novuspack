package generics

import (
	"context"
	"testing"
	"time"
)

// TestThreadSafetyMode tests ThreadSafetyMode constants
func TestThreadSafetyMode(t *testing.T) {
	if ThreadSafetyNone != 0 {
		t.Error("ThreadSafetyNone should be 0")
	}
	if ThreadSafetyReadOnly != 1 {
		t.Error("ThreadSafetyReadOnly should be 1")
	}
	if ThreadSafetyConcurrent != 2 {
		t.Error("ThreadSafetyConcurrent should be 2")
	}
	if ThreadSafetyFull != 3 {
		t.Error("ThreadSafetyFull should be 3")
	}
}

// TestConcurrencyConfig tests ConcurrencyConfig
func TestConcurrencyConfig(t *testing.T) {
	config := &ConcurrencyConfig{
		MaxWorkers:       4,
		WorkerTimeout:    5 * time.Second,
		WorkerBufferSize: 10,
		UseMutex:         true,
		UseRWMutex:       true,
		LockTimeout:      1 * time.Second,
		MaxConcurrentOps: 100,
		ResourcePoolSize: 20,
	}

	if config.MaxWorkers != 4 {
		t.Errorf("MaxWorkers should be 4, got %d", config.MaxWorkers)
	}
	if config.WorkerTimeout != 5*time.Second {
		t.Errorf("WorkerTimeout should be 5s, got %v", config.WorkerTimeout)
	}
}

// TestNewWorkerPool tests WorkerPool creation
func TestNewWorkerPool(t *testing.T) {
	config := &ConcurrencyConfig{
		MaxWorkers: 2,
	}
	pool := NewWorkerPool[string](config)

	//nolint:staticcheck // SA5011: false positive - t.Fatal exits, so pool is not nil after check
	if pool == nil {
		t.Fatal("NewWorkerPool should not return nil")
	}
	//nolint:staticcheck // SA5011: false positive - t.Fatal exits, so pool is not nil after check
	if pool.config.MaxWorkers != 2 {
		t.Errorf("WorkerPool should have MaxWorkers=2, got %d", pool.config.MaxWorkers)
	}
}

func TestNewWorkerPool_DefaultConfig(t *testing.T) {
	pool := NewWorkerPool[string](nil)
	//nolint:staticcheck // SA5011: false positive - t.Fatal exits, so pool is not nil after check
	if pool == nil {
		t.Fatal("NewWorkerPool should not return nil with nil config")
	}
	//nolint:staticcheck // SA5011: false positive - t.Fatal exits, so pool is not nil after check
	if pool.config.MaxWorkers == 0 {
		t.Error("NewWorkerPool should set default MaxWorkers")
	}
}

// TestWorkerPool_StartStop tests starting and stopping worker pool
func TestWorkerPool_StartStop(t *testing.T) {
	config := &ConcurrencyConfig{
		MaxWorkers:    2,
		WorkerTimeout: 1 * time.Second,
	}
	pool := NewWorkerPool[string](config)
	ctx := context.Background()

	// Start pool
	err := pool.Start(ctx)
	if err != nil {
		t.Fatalf("Start should succeed, got error: %v", err)
	}

	stats := pool.GetWorkerStats()
	if !stats.IsRunning {
		t.Error("Worker pool should be running after Start")
	}
	if stats.ActiveWorkers != 2 {
		t.Errorf("ActiveWorkers should be 2, got %d", stats.ActiveWorkers)
	}

	// Stop pool
	err = pool.Stop(ctx)
	if err != nil {
		t.Fatalf("Stop should succeed, got error: %v", err)
	}

	stats = pool.GetWorkerStats()
	if stats.IsRunning {
		t.Error("Worker pool should not be running after Stop")
	}
}

func TestWorkerPool_StartTwice(t *testing.T) {
	pool := NewWorkerPool[string](&ConcurrencyConfig{MaxWorkers: 1})
	ctx := context.Background()

	err := pool.Start(ctx)
	if err != nil {
		t.Fatalf("First Start should succeed, got error: %v", err)
	}

	err = pool.Start(ctx)
	if err == nil {
		t.Error("Second Start should return error")
	}

	_ = pool.Stop(ctx) // Error handling not needed in test cleanup
}

func TestWorkerPool_StopNotRunning(t *testing.T) {
	pool := NewWorkerPool[string](&ConcurrencyConfig{MaxWorkers: 1})
	ctx := context.Background()

	err := pool.Stop(ctx)
	if err == nil {
		t.Error("Stop on non-running pool should return error")
	}
}

// TestWorkerPool_SubmitJob tests job submission
func TestWorkerPool_SubmitJob(t *testing.T) {
	config := &ConcurrencyConfig{
		MaxWorkers: 1,
	}
	pool := NewWorkerPool[string](config)
	ctx := context.Background()

	err := pool.Start(ctx)
	if err != nil {
		t.Fatalf("Start should succeed, got error: %v", err)
	}
	defer func() { _ = pool.Stop(ctx) }() // Stop on exit - error is non-critical

	resultChan := make(chan Result[string], 1)
	job := Job[string]{
		ID:      "test-job",
		Data:    "test-data",
		Result:  resultChan,
		Context: ctx,
	}

	err = pool.SubmitJob(ctx, job)
	if err != nil {
		t.Fatalf("SubmitJob should succeed, got error: %v", err)
	}
}

func TestWorkerPool_SubmitJobNotRunning(t *testing.T) {
	pool := NewWorkerPool[string](&ConcurrencyConfig{MaxWorkers: 1})
	ctx := context.Background()

	resultChan := make(chan Result[string], 1)
	job := Job[string]{
		ID:      "test-job",
		Data:    "test-data",
		Result:  resultChan,
		Context: ctx,
	}

	err := pool.SubmitJob(ctx, job)
	if err == nil {
		t.Error("SubmitJob on non-running pool should return error")
	}
}

// TestWorkerPool_GetWorkerStats tests statistics
func TestWorkerPool_GetWorkerStats(t *testing.T) {
	config := &ConcurrencyConfig{
		MaxWorkers: 3,
	}
	pool := NewWorkerPool[string](config)

	stats := pool.GetWorkerStats()
	if stats.TotalWorkers != 3 {
		t.Errorf("TotalWorkers should be 3, got %d", stats.TotalWorkers)
	}
	if stats.IsRunning {
		t.Error("Worker pool should not be running initially")
	}
}

// TestProcessConcurrently tests concurrent processing
func TestProcessConcurrently(t *testing.T) {
	ctx := context.Background()
	items := []string{"item1", "item2", "item3"}

	processor := &testStrategy{
		name:         "test",
		strategyType: "test",
	}

	config := &ConcurrencyConfig{
		MaxWorkers: 2,
	}

	results, err := ProcessConcurrently(ctx, items, processor, config)
	if err != nil {
		t.Fatalf("ProcessConcurrently should succeed, got error: %v", err)
	}

	if len(results) != len(items) {
		t.Errorf("ProcessConcurrently should return %d results, got %d", len(items), len(results))
	}

	for i, result := range results {
		if !result.IsOk() {
			t.Errorf("Result %d should be Ok", i)
		}
		val, err := result.Unwrap()
		if err != nil {
			t.Errorf("Result %d should not have error, got %v", i, err)
		}
		expected := "processed: " + items[i]
		if val != expected {
			t.Errorf("Result %d should be %q, got %q", i, expected, val)
		}
	}
}

func TestProcessConcurrently_EmptySlice(t *testing.T) {
	ctx := context.Background()
	items := []string{}

	processor := &testStrategy{
		name:         "test",
		strategyType: "test",
	}

	results, err := ProcessConcurrently(ctx, items, processor, nil)
	if err != nil {
		t.Fatalf("ProcessConcurrently should succeed with empty slice, got error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("ProcessConcurrently should return empty results for empty slice, got %d", len(results))
	}
}

func TestProcessConcurrently_ErrorHandling(t *testing.T) {
	ctx := context.Background()
	items := []string{"item1", "error", "item3"}

	processor := &testStrategy{
		name:         "test",
		strategyType: "test",
	}

	config := &ConcurrencyConfig{
		MaxWorkers: 2,
	}

	results, err := ProcessConcurrently(ctx, items, processor, config)
	if err != nil {
		t.Fatalf("ProcessConcurrently should succeed even with errors, got error: %v", err)
	}

	if len(results) != len(items) {
		t.Errorf("ProcessConcurrently should return %d results, got %d", len(items), len(results))
	}

	// Check that the error result is properly handled
	result := results[1]
	if !result.IsErr() {
		t.Error("Result for 'error' item should be Err")
	}
}

// TestJob tests Job structure
func TestJob(t *testing.T) {
	ctx := context.Background()
	resultChan := make(chan Result[string], 1)

	job := Job[string]{
		ID:       "test-id",
		Data:     "test-data",
		Result:   resultChan,
		Context:  ctx,
		Priority: 10,
	}

	if job.ID != "test-id" {
		t.Errorf("Job ID should be 'test-id', got %q", job.ID)
	}
	if job.Data != "test-data" {
		t.Errorf("Job Data should be 'test-data', got %q", job.Data)
	}
	if job.Priority != 10 {
		t.Errorf("Job Priority should be 10, got %d", job.Priority)
	}
}
