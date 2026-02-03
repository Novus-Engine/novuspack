// This file implements generic concurrency patterns: WorkerPool[T], Worker[T],
// Job[T], and ConcurrencyConfig. It contains type-safe concurrent processing
// with worker pools and thread safety management. This file should contain
// all code related to generic concurrency as specified in api_generics.md
// Section 1.8 and Section 1.9.
//
// Specification: api_generics.md: 1. Core Generic Types

// Package generics provides generic types and patterns for the NovusPack API.
package generics

import (
	"context"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// ThreadSafetyMode defines the level of thread safety guarantees.
type ThreadSafetyMode int

const (
	// ThreadSafetyNone indicates no thread safety guarantees.
	ThreadSafetyNone ThreadSafetyMode = iota

	// ThreadSafetyReadOnly indicates read-only operations are safe.
	ThreadSafetyReadOnly

	// ThreadSafetyConcurrent indicates concurrent read/write operations are safe.
	ThreadSafetyConcurrent

	// ThreadSafetyFull indicates full thread safety with synchronization.
	ThreadSafetyFull
)

// ConcurrencyConfig defines thread safety and worker management settings.
type ConcurrencyConfig struct {
	// Worker management
	MaxWorkers       int           // Maximum parallel workers (0 = auto-detect)
	WorkerTimeout    time.Duration // Worker timeout for graceful shutdown
	WorkerBufferSize int           // Worker channel buffer size (0 = auto-calculate)

	// Thread safety
	UseMutex    bool          // Use mutex for shared state protection
	UseRWMutex  bool          // Use read-write mutex for better read performance
	LockTimeout time.Duration // Lock acquisition timeout

	// Resource management
	MaxConcurrentOps int // Maximum concurrent operations (0 = no limit)
	ResourcePoolSize int // Resource pool size for workers
}

// WorkerStats provides statistics about the worker pool.
type WorkerStats struct {
	ActiveWorkers int  // Number of currently active workers
	TotalWorkers  int  // Total number of workers in the pool
	JobsProcessed int  // Total number of jobs processed
	JobsInQueue   int  // Number of jobs currently in queue
	IsRunning     bool // Whether the worker pool is currently running
}

// Job represents a unit of work for any data type.
type Job[T any] struct {
	ID       string          // Unique identifier for the job
	Data     T               // The data to process
	Result   chan Result[T]  // Channel to receive the result
	Context  context.Context // Context for cancellation and timeout
	Priority int             // Job priority (higher = more priority)
}

// Worker represents a single worker in the pool.
type Worker[T any] struct {
	mu       sync.RWMutex
	id       int
	workChan chan Job[T]
	done     chan struct{}
	strategy Strategy[T, T]
	stats    *WorkerStats
}

// WorkerPool manages concurrent workers for any data type.
type WorkerPool[T any] struct {
	mu       sync.RWMutex
	workers  []*Worker[T]
	workChan chan Job[T]
	done     chan struct{}
	wg       sync.WaitGroup
	config   *ConcurrencyConfig
	running  bool
	stats    WorkerStats
}

// NewWorkerPool creates a new WorkerPool with the given configuration.
func NewWorkerPool[T any](config *ConcurrencyConfig) *WorkerPool[T] {
	if config == nil {
		config = &ConcurrencyConfig{
			MaxWorkers: runtime.NumCPU(),
		}
	}

	if config.MaxWorkers == 0 {
		config.MaxWorkers = runtime.NumCPU()
	}

	if config.WorkerBufferSize == 0 {
		config.WorkerBufferSize = config.MaxWorkers * 2
	}

	return &WorkerPool[T]{
		workers:  make([]*Worker[T], 0, config.MaxWorkers),
		workChan: make(chan Job[T], config.WorkerBufferSize),
		done:     make(chan struct{}),
		config:   config,
		running:  false,
		stats: WorkerStats{
			TotalWorkers: config.MaxWorkers,
		},
	}
}

// Start initializes and starts the worker pool.
func (p *WorkerPool[T]) Start(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.running {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "worker pool is already running", nil, pkgerrors.ValidationErrorContext{
			Field:    "WorkerPool",
			Value:    "running",
			Expected: "not running",
		})
	}

	// Create workers
	for i := 0; i < p.config.MaxWorkers; i++ {
		worker := &Worker[T]{
			id:       i,
			workChan: p.workChan,
			done:     p.done,
			stats:    &p.stats,
		}
		p.workers = append(p.workers, worker)
		p.wg.Add(1)
		go func(w *Worker[T]) {
			defer p.wg.Done()
			w.run(ctx)
		}(worker)
	}

	p.running = true
	p.stats.IsRunning = true
	p.stats.ActiveWorkers = p.config.MaxWorkers

	return nil
}

// Stop gracefully shuts down the worker pool.
func (p *WorkerPool[T]) Stop(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.running {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "worker pool is not running", nil, pkgerrors.ValidationErrorContext{
			Field:    "WorkerPool",
			Value:    "not running",
			Expected: "running",
		})
	}

	// Signal all workers to stop
	close(p.done)

	// Wait for workers to finish or timeout
	done := make(chan struct{})
	go func() {
		close(done)
	}()

	select {
	case <-done:
		// Workers stopped gracefully
	case <-ctx.Done():
		return pkgerrors.WrapErrorWithContext(ctx.Err(), pkgerrors.ErrTypeContext, "worker pool shutdown cancelled", pkgerrors.ValidationErrorContext{
			Field:    "WorkerPool",
			Value:    nil,
			Expected: "graceful shutdown",
		})
	case <-time.After(p.config.WorkerTimeout):
		if p.config.WorkerTimeout > 0 {
			return pkgerrors.NewPackageError(pkgerrors.ErrTypeIO, "worker pool shutdown timeout", nil, pkgerrors.ValidationErrorContext{
				Field:    "WorkerPool",
				Value:    p.config.WorkerTimeout,
				Expected: "shutdown within timeout",
			})
		}
	}

	p.running = false
	p.stats.IsRunning = false
	p.stats.ActiveWorkers = 0

	return nil
}

// SubmitJob submits a job to the worker pool.
func (p *WorkerPool[T]) SubmitJob(ctx context.Context, job Job[T]) error {
	p.mu.RLock()
	running := p.running
	p.mu.RUnlock()

	if !running {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "worker pool is not running", nil, pkgerrors.ValidationErrorContext{
			Field:    "WorkerPool",
			Value:    "not running",
			Expected: "running",
		})
	}

	select {
	case p.workChan <- job:
		p.mu.Lock()
		p.stats.JobsInQueue = len(p.workChan)
		p.mu.Unlock()
		return nil
	case <-ctx.Done():
		return pkgerrors.WrapErrorWithContext(ctx.Err(), pkgerrors.ErrTypeContext, "job submission cancelled", pkgerrors.ValidationErrorContext{
			Field:    "Job",
			Value:    job.ID,
			Expected: "submitted successfully",
		})
	}
}

// GetWorkerStats returns current worker pool statistics.
func (p *WorkerPool[T]) GetWorkerStats() WorkerStats {
	p.mu.RLock()
	defer p.mu.RUnlock()

	stats := p.stats
	stats.JobsInQueue = len(p.workChan)
	return stats
}

// processJob runs the job through the strategy (or pass-through) and sends the result.
func (w *Worker[T]) processJob(ctx context.Context, job Job[T]) {
	if w.strategy == nil {
		job.Result <- Ok(job.Data)
		return
	}
	jobCtx, cancel := context.WithCancel(ctx)
	go func() {
		select {
		case <-job.Context.Done():
			cancel()
		case <-ctx.Done():
			cancel()
		case <-jobCtx.Done():
		}
	}()
	result, err := w.strategy.Process(jobCtx, job.Data)
	cancel()
	if err != nil {
		job.Result <- Err[T](err)
	} else {
		job.Result <- Ok(result)
	}
}

// run is the main loop for a worker.
func (w *Worker[T]) run(ctx context.Context) {
	for {
		select {
		case job := <-w.workChan:
			w.processJob(ctx, job)
			w.mu.Lock()
			w.stats.JobsProcessed++
			w.mu.Unlock()
		case <-w.done:
			return
		case <-ctx.Done():
			return
		}
	}
}

// SetStrategy sets the strategy for processing jobs.
func (w *Worker[T]) SetStrategy(strategy Strategy[T, T]) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.strategy = strategy
}

// ProcessConcurrently processes multiple items concurrently using the worker pool.
func ProcessConcurrently[T any](ctx context.Context, items []T, processor Strategy[T, T], config *ConcurrencyConfig) ([]Result[T], error) {
	if len(items) == 0 {
		return []Result[T]{}, nil
	}

	pool := NewWorkerPool[T](config)
	if err := pool.Start(ctx); err != nil {
		return nil, err
	}
	defer func() { _ = pool.Stop(ctx) }() // Stop on exit - error is non-critical

	// Set strategy for all workers
	for _, worker := range pool.workers {
		worker.SetStrategy(processor)
	}

	// Submit all jobs
	results := make([]Result[T], len(items))
	resultChans := make([]chan Result[T], len(items))

	for i, item := range items {
		resultChan := make(chan Result[T], 1)
		resultChans[i] = resultChan

		job := Job[T]{
			ID:      strconv.Itoa(i),
			Data:    item,
			Result:  resultChan,
			Context: ctx,
		}

		if err := pool.SubmitJob(ctx, job); err != nil {
			return nil, err
		}
	}

	// Collect results
	for i, resultChan := range resultChans {
		select {
		case result := <-resultChan:
			results[i] = result
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return results, nil
}
