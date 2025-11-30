@domain:file_mgmt @m2 @REQ-FILEMGMT-204 @spec(api_file_management.md#9211-parameters)
Feature: GetFileByFileID Parameter Specification

  @REQ-FILEMGMT-204 @happy
  Scenario: GetFileByFileID parameters include context and fileID
    Given an open NovusPack package
    And a valid context
    And a 64-bit file identifier
    When GetFileByFileID is called
    Then context parameter is accepted for cancellation and timeout handling
    And fileID parameter is accepted as unique 64-bit identifier
    And parameters are validated

  @REQ-FILEMGMT-204 @happy
  Scenario: GetFileByFileID context supports cancellation
    Given an open NovusPack package
    And a context that can be cancelled
    And a 64-bit file identifier
    When GetFileByFileID is called
    And context is cancelled
    Then operation respects context cancellation
    And structured context error is returned

  @REQ-FILEMGMT-204 @happy
  Scenario: GetFileByFileID context supports timeout handling
    Given an open NovusPack package
    And a context with timeout
    And a 64-bit file identifier
    When GetFileByFileID is called
    And timeout is exceeded
    Then operation respects context timeout
    And structured context timeout error is returned

  @REQ-FILEMGMT-204 @error
  Scenario: GetFileByFileID handles package not open errors
    Given a closed NovusPack package
    And a valid context
    And a 64-bit file identifier
    When GetFileByFileID is called
    Then a structured error is returned
    And error indicates package is not open
    And error follows structured error format

  @REQ-FILEMGMT-204 @error
  Scenario: GetFileByFileID handles invalid fileID
    Given an open NovusPack package
    And a valid context
    And a non-existent file identifier
    When GetFileByFileID is called with non-existent fileID
    Then FileEntry is not found
    And boolean false is returned
    And no error is returned for not found case
