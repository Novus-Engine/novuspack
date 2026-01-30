@domain:file_mgmt @REQ-PIPELINE-025 @REQ-PIPELINE-026 @REQ-CRYPTO-013 @spec(api_file_mgmt_transform_pipelines.md#17-temporary-file-security)
Feature: Temporary File Security Requirements

  @REQ-PIPELINE-025 @REQ-CRYPTO-013 @happy
  Scenario: Temp files for encrypted content secured on disk
    Given a pipeline processing encrypted content
    And transformation creates temporary files
    When general security rule applies
    Then temp files for encrypted content are encrypted on disk
    And sensitive data is protected at rest
    And security follows context-aware rules

  @REQ-PIPELINE-025 @REQ-CRYPTO-015 @happy
  Scenario: Decrypt operations exception for temp file security
    Given user decrypting content to write to disk
    And user intent is to decrypt
    When decrypt transformation creates temp file
    Then unencrypted temp files are acceptable
    And temp files in same directory as destination
    And security exception reflects user intent

  @REQ-PIPELINE-026 @REQ-CRYPTO-014 @happy
  Scenario: Secure cleanup with overwrites for sensitive files
    Given temporary files containing sensitive data
    When cleanup is performed
    Then sensitive temp files are overwritten before removal
    And secure cleanup prevents data recovery
    And overwrites follow security best practices

  @REQ-PIPELINE-026 @happy
  Scenario: OS temp directory permissions as baseline
    Given temporary files created for transformations
    Then files created with OS temp directory permissions
    And permissions provide baseline security
    And appropriate for temporary staging files
    And consistent with OS security model

  @REQ-PIPELINE-025 @happy
  Scenario: Atomic operations for temp file security
    Given transformation stage creating output
    When temporary file is written
    Then operation is atomic where possible
    And partial writes are minimized
    And security is maintained during writes
