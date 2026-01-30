@domain:file_mgmt @extraction @REQ-API_BASIC-144 @REQ-API_BASIC-145 @spec(api_basic_operations.md#19-package-session-base-management)
Feature: Session base runtime property for extraction

  @REQ-API_BASIC-144 @REQ-API_BASIC-145 @happy
  Scenario: SetSessionBase sets runtime property for extraction
    Given an open NovusPack package
    When SetSessionBase is called with "/tmp/extract"
    Then session base is set to "/tmp/extract"
    And GetSessionBase returns "/tmp/extract"

  @REQ-API_BASIC-145 @happy
  Scenario: GetSessionBase returns current runtime property for extraction
    Given an open NovusPack package
    And session base is set to "/tmp/extract"
    When GetSessionBase is called
    Then "/tmp/extract" is returned

  @REQ-API_BASIC-145 @happy
  Scenario: GetSessionBase returns empty string when not set
    Given an open NovusPack package
    And session base is not set
    When GetSessionBase is called
    Then empty string is returned

  @REQ-API_BASIC-144 @happy
  Scenario: SetSessionBase validates path format for extraction
    Given an open NovusPack package
    When SetSessionBase is called with valid path "/tmp/extract"
    Then session base is set successfully
    And no error is returned

  @REQ-API_BASIC-144 @error
  Scenario: SetSessionBase returns error for invalid path
    Given an open NovusPack package
    When SetSessionBase is called with invalid path
    Then ErrTypeValidation error is returned
    And error indicates invalid path format

  @REQ-API_BASIC-144 @happy
  Scenario: ExtractPath updates runtime session base from options
    Given an open NovusPack package
    And file "data.txt" exists
    When ExtractPath is called with SessionBase option "/tmp/extract"
    Then package runtime session base is set to "/tmp/extract"
    And GetSessionBase returns "/tmp/extract"

  @REQ-API_BASIC-144 @happy
  Scenario: Session base persists across multiple ExtractPath calls
    Given an open NovusPack package
    And session base is set to "/tmp/extract"
    And files "a.txt" and "b.txt" exist
    When ExtractPath is called for "a.txt"
    And ExtractPath is called for "b.txt"
    Then both files are extracted to "/tmp/extract/"
    And session base remains "/tmp/extract"
