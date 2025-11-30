@domain:file_mgmt @m2 @REQ-FILEMGMT-128 @spec(api_file_management.md#132-encryption-management)
Feature: Encryption Management Best Practices

  @REQ-FILEMGMT-128 @happy
  Scenario: Encryption management best practices define encryption patterns
    Given an open NovusPack package
    When encryption management best practices are applied
    Then encryption patterns are defined
    And best practices guide encryption usage

  @REQ-FILEMGMT-128 @REQ-FILEMGMT-129 @happy
  Scenario: Choose appropriate encryption types for sensitive data
    Given an open package
    And sensitive file data
    When encryption is applied
    Then AES-256-GCM can be used for sensitive data
    And ML-KEM can be used for post-quantum security
    And encryption type matches security requirements

  @REQ-FILEMGMT-128 @REQ-FILEMGMT-130 @happy
  Scenario: Secure key management follows best practices
    Given an open package
    And ML-KEM keys are available
    When encryption keys are managed
    Then keys are generated at appropriate security levels
    And keys are securely stored
    And sensitive data is cleared after use
    And key management follows security best practices

  @REQ-FILEMGMT-128 @happy
  Scenario: Encryption management supports per-file encryption selection
    Given an open package
    When files are added with encryption
    Then encryption can be selected per file
    And ML-KEM encryption can be used per file
    And AES-256-GCM encryption can be used per file
    And encryption type matches file requirements

  @REQ-FILEMGMT-128 @error
  Scenario: Encryption management handles invalid key errors
    Given an open package
    And invalid or missing encryption keys
    When encryption operation is attempted
    Then structured encryption error is returned
    And error indicates key management issue
    And error follows structured error format

  @REQ-FILEMGMT-128 @error
  Scenario: Encryption management respects context cancellation
    Given an open package
    And a cancelled context
    When encryption operation is called
    Then structured context error is returned
    And error follows structured error format
