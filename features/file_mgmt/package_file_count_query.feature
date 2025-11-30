@domain:file_mgmt @m2 @REQ-FILEMGMT-227 @REQ-FILEMGMT-228 @REQ-FILEMGMT-229 @REQ-FILEMGMT-230 @REQ-FILEMGMT-231 @spec(api_file_management.md#926-getfilecount)
Feature: Package file count query

  @REQ-FILEMGMT-228 @REQ-FILEMGMT-229 @happy
  Scenario: GetFileCount returns total number of files
    Given an open package
    And the package contains 5 files
    When GetFileCount is called
    Then the count 5 is returned
    And no error occurs

  @REQ-FILEMGMT-228 @REQ-FILEMGMT-229 @happy
  Scenario: GetFileCount returns zero for empty package
    Given an open package
    And the package contains no files
    When GetFileCount is called
    Then the count 0 is returned
    And no error occurs

  @REQ-FILEMGMT-231 @happy
  Scenario: GetFileCount updates after file addition
    Given an open package with 3 files
    And GetFileCount returns 3
    When a new file is added
    Then GetFileCount returns 4

  @REQ-FILEMGMT-231 @happy
  Scenario: GetFileCount updates after file removal
    Given an open package with 5 files
    And GetFileCount returns 5
    When a file is removed
    Then GetFileCount returns 4

  @REQ-FILEMGMT-231 @happy
  Scenario: GetFileCount is accurate for large packages
    Given an open package with 1000 files
    When GetFileCount is called
    Then the count 1000 is returned
    And count matches actual file entries

  @REQ-FILEMGMT-229 @error
  Scenario: GetFileCount respects context cancellation
    Given an open package
    And a cancelled context
    When GetFileCount is called
    Then a structured context error is returned
    And error type is context cancellation

  @REQ-FILEMGMT-230 @happy
  Scenario: GetFileCount returns integer count efficiently
    Given an open package with many files
    When GetFileCount is called
    Then the count is returned quickly
    And performance is acceptable for large packages
