@domain:file_mgmt @m2 @REQ-FILEMGMT-179 @spec(api_file_management.md#71-purpose)
Feature: CompressFile Purpose and Usage

  @REQ-FILEMGMT-179 @happy
  Scenario: File compression operations manage file-level compression
    Given an open NovusPack package
    And a valid context
    When file compression operations are used
    Then file-level compression operations are available
    And existing files in package can be compressed
    And compression can be managed per file

  @REQ-FILEMGMT-179 @happy
  Scenario: File compression operations support per-file compression
    Given an open NovusPack package
    And a valid context
    And files exist in the package
    When file compression operations are performed
    Then individual files can be compressed
    And compression settings can be applied per file
    And file compression state can be managed

  @REQ-FILEMGMT-179 @happy
  Scenario: File compression operations provide compression management
    Given an open NovusPack package
    And a valid context
    When compression operations are performed
    Then compression can be applied to files
    Then compression can be removed from files
    And compression information can be retrieved

  @REQ-FILEMGMT-179 @error
  Scenario: File compression operations require package to be open
    Given a closed NovusPack package
    And a valid context
    When compression operations are attempted
    Then a structured error is returned
    And error indicates package is not open
    And error follows structured error format
