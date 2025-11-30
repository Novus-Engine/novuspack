@domain:basic_ops @m2 @REQ-API_BASIC-055 @spec(api_basic_operations.md#721-defragment-behavior)
Feature: Defragment behavior

  @REQ-API_BASIC-055 @happy
  Scenario: Defragment removes unused space from deleted files
    Given an open NovusPack package
    And package has deleted files with unused space
    When Defragment is called
    Then unused space is removed
    And deleted file space is reclaimed
    And package structure is optimized

  @REQ-API_BASIC-055 @happy
  Scenario: Defragment reorganizes file entries for optimal access
    Given an open NovusPack package
    When Defragment is called
    Then file entries are reorganized
    And entries are optimized for access
    And package access performance improves

  @REQ-API_BASIC-055 @happy
  Scenario: Defragment compacts data sections
    Given an open NovusPack package
    And package has fragmented data sections
    When Defragment is called
    Then data sections are compacted
    And file size is reduced
    And package structure is more efficient

  @REQ-API_BASIC-055 @happy
  Scenario: Defragment preserves package metadata and signatures
    Given an open NovusPack package
    And package has metadata and signatures
    When Defragment is called
    Then package metadata is preserved
    And digital signatures are preserved
    And package integrity is maintained

  @REQ-API_BASIC-055 @happy
  Scenario: Defragment may take significant time for large packages
    Given an open NovusPack package
    And package is large with many files
    When Defragment is called
    Then operation may take significant time
    And progress can be monitored via context
    And operation completes successfully

  @REQ-API_BASIC-055 @error
  Scenario: Defragment returns error when package is not open
    Given a NovusPack package that is not open
    When Defragment is called
    Then validation error is returned
    And error indicates package is not open

  @REQ-API_BASIC-055 @error
  Scenario: Defragment returns error in read-only mode
    Given an open NovusPack package in read-only mode
    When Defragment is called
    Then validation error is returned
    And error indicates package is read-only

  @REQ-API_BASIC-055 @error
  Scenario: Defragment handles I/O errors during defragmentation
    Given an open NovusPack package
    And file system error occurs during defragmentation
    When Defragment is called
    Then I/O error is returned
    And error indicates file system issue

  @REQ-API_BASIC-055 @error
  Scenario: Defragment respects context cancellation
    Given an open NovusPack package
    And a cancelled context
    When Defragment is called with cancelled context
    Then context error is returned
    And error type is context cancellation
