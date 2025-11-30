@domain:file_mgmt @m2 @spec(api_file_management.md#51-extract-file)
Feature: Extract file bytes

  @REQ-FILEMGMT-003 @happy
  Scenario: Extract obeys encryption and validation
    Given a package with an encrypted file "secret.txt"
    When I extract the file with correct keys
    Then the extracted bytes should match the original content

  @REQ-FILEMGMT-037 @REQ-FILEMGMT-038 @error
  Scenario: ExtractFile validates path parameter
    Given an open package
    When ExtractFile is called with empty path
    Then structured validation error is returned
    And error indicates invalid path

  @REQ-FILEMGMT-037 @REQ-FILEMGMT-041 @error
  Scenario: ExtractFile respects context cancellation
    Given an open package with files
    And a cancelled context
    When ExtractFile is called
    Then structured context error is returned
    And error type is context cancellation
