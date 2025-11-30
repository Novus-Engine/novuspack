@domain:writing @m2 @REQ-WRITE-001 @spec(api_writing.md#1-safewrite---atomic-package-writing)
Feature: SafeWrite atomic package writing

  @happy
  Scenario: SafeWrite uses temp file and atomic rename
    Given a package pending write operations
    When I perform a safe write
    Then a temp file should be used and an atomic rename should finalize

  @happy
  Scenario: SafeWrite creates temp file in same directory
    Given a package to be written
    When SafeWrite is called
    Then temporary file is created in same directory as target
    And temp file has unique name
    And temp file is used for writing

  @happy
  Scenario: SafeWrite streams data for large files
    Given a package with large file content
    When SafeWrite is called
    Then data is streamed from source
    And streaming handles large content efficiently
    And memory usage is controlled

  @happy
  Scenario: SafeWrite uses in-memory data for small files
    Given a package with small file content
    When SafeWrite is called
    Then data is written from memory
    And in-memory writing is efficient
    And memory thresholds are respected

  @happy
  Scenario: SafeWrite atomically renames temp file to target
    Given a package written to temp file
    When SafeWrite completes successfully
    Then temp file is atomically renamed to target path
    And original file is replaced atomically
    And no partial writes are possible

  @happy
  Scenario: SafeWrite compresses package content when specified
    Given a package to be written
    When SafeWrite is called with compression type
    Then package content is compressed
    And file entries, data, and index are compressed
    And header, comment, and signatures remain uncompressed

  @happy
  Scenario: SafeWrite handles new package creation
    Given a new package
    When SafeWrite is called
    Then new package file is created
    And package structure is written correctly
    And package is ready for use

  @happy
  Scenario: SafeWrite handles complete package rewrite
    Given an existing package requiring complete rewrite
    When SafeWrite is called
    Then complete package is rewritten
    And new package replaces old package atomically
    And package integrity is maintained

  @happy
  Scenario: SafeWrite handles defragmentation operations
    Given a package requiring defragmentation
    When SafeWrite is called
    Then defragmentation is performed
    And package structure is optimized
    And write operation is atomic

  @happy
  Scenario: SafeWrite cleans up temp file on failure
    Given a package write operation that fails
    When SafeWrite encounters an error
    Then temp file is automatically cleaned up
    And no temporary files remain
    And package state is unchanged

  @happy
  Scenario: SafeWrite validates target directory
    Given a package write operation
    When SafeWrite is called
    Then target directory existence is validated
    And directory permissions are checked
    And validation occurs before writing

  @error
  Scenario: SafeWrite fails if directory does not exist
    Given a package write operation to non-existent directory
    When SafeWrite is called
    Then structured validation error is returned
    And error indicates directory not found

  @error
  Scenario: SafeWrite fails on I/O errors
    Given a package write operation
    And file system I/O errors occur
    When SafeWrite is called
    Then structured I/O error is returned
    And temp file is cleaned up
    And package state is unchanged

  @error
  Scenario: SafeWrite respects context cancellation
    Given a long-running package write operation
    And a cancelled context
    When SafeWrite is called
    Then structured context error is returned
    And temp file is cleaned up
    And operation is cancelled
