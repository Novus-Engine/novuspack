@domain:file_mgmt @extraction @concurrency @performance @REQ-FILEMGMT-389 @REQ-FILEMGMT-390 @REQ-FILEMGMT-391 @REQ-FILEMGMT-392 @REQ-FILEMGMT-393 @REQ-FILEMGMT-394 @spec(api_file_mgmt_extraction.md#524-extractoptions-concurrency)
Feature: Concurrent extraction of multiple files

  As a package user
  I want to extract multiple files concurrently
  So that extraction completes faster on multi-core systems

  @REQ-FILEMGMT-389 @happy
  Scenario: Extract multiple files concurrently using worker pool
    Given a package with 100 files
    And EnableConcurrentExtraction is enabled
    And MaxConcurrentExtractions is set to 4
    When I extract all files
    Then files should be extracted using a worker pool
    And up to 4 files should be extracted simultaneously
    And all files should be extracted successfully

  @REQ-FILEMGMT-390 @happy
  Scenario: Thread-safe disk space tracking during concurrent extraction
    Given a package with 100 files
    And EnableConcurrentExtraction is enabled
    When I extract all files concurrently
    Then disk space tracking should be thread-safe
    And space reservations should be coordinated across workers
    And no race conditions should occur

  @REQ-FILEMGMT-391 @happy
  Scenario: Space reservation before extraction prevents race conditions
    Given a package with 100 files
    And EnableConcurrentExtraction is enabled
    When I extract all files concurrently
    Then space should be reserved for each file before extraction begins
    And total reserved space should not exceed available space
    And all reservations should be tracked thread-safely

  @REQ-FILEMGMT-392 @happy
  Scenario: Coordinated cancellation across all workers
    Given a package with 100 files
    And EnableConcurrentExtraction is enabled
    And extraction is in progress
    When context cancellation is triggered
    Then all worker threads should receive cancellation signal
    And all workers should stop extraction gracefully
    And partial files should be cleaned up from all workers

  @REQ-FILEMGMT-393 @happy
  Scenario: Progress callbacks include disk space and worker status
    Given a package with 100 files
    And EnableConcurrentExtraction is enabled
    When I extract all files with progress callback
    Then progress callback should include current file count
    And progress callback should include total file count
    And progress callback should include available disk space
    And progress callback should include worker status

  @REQ-FILEMGMT-394 @happy
  Scenario: Context cancellation cleans up partial files from all workers
    Given a package with 100 files
    And EnableConcurrentExtraction is enabled
    And extraction is in progress
    When context cancellation is triggered
    Then all partial files from all workers should be cleaned up
    And no orphaned files should remain
    And disk space reservations should be released

  @REQ-FILEMGMT-390 @error
  Scenario: Concurrent extraction stops when disk space exhausted
    Given a package with 100 files
    And EnableConcurrentExtraction is enabled
    And disk space becomes exhausted during concurrent extraction
    When space check detects insufficient space
    Then all workers should be cancelled
    And extraction should stop immediately
    And error should indicate disk space exhausted

  @REQ-FILEMGMT-389 @happy
  Scenario: Default concurrent extraction uses CPU core count
    Given a package with 100 files
    And EnableConcurrentExtraction is enabled
    And MaxConcurrentExtractions is not set
    When I extract all files
    Then number of concurrent workers should equal CPU core count
    And extraction should complete successfully
