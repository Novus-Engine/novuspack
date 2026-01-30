@domain:file_mgmt @m2 @REQ-FILEMGMT-114 @spec(api_file_mgmt_errors.md#1222-error-type-mapping) @spec(api_file_mgmt_file_entry.md#1221-error-type-categories)
Feature: Error Type Mapping

  @REQ-FILEMGMT-114 @happy
  Scenario: Error type mapping maps sentinel errors to structured errors
    Given an open NovusPack package
    And a valid context
    When sentinel errors occur
    Then sentinel errors are mapped to structured errors
    And error types are correctly mapped
    And error descriptions are preserved

  @REQ-FILEMGMT-114 @happy
  Scenario: Error type mapping provides validation error mapping
    Given an open NovusPack package
    And a valid context
    When validation sentinel errors occur
    Then ErrFileNotFound maps to ErrTypeValidation
    And ErrFileExists maps to ErrTypeValidation
    And ErrInvalidPath maps to ErrTypeValidation
    And ErrInvalidPattern maps to ErrTypeValidation

  @REQ-FILEMGMT-114 @happy
  Scenario: Error type mapping provides encryption error mapping
    Given an open NovusPack package
    And a valid context
    When encryption sentinel errors occur
    Then ErrUnsupportedEncryption maps to ErrTypeEncryption
    And ErrEncryptionFailed maps to ErrTypeEncryption
    And ErrDecryptionFailed maps to ErrTypeEncryption

  @REQ-FILEMGMT-114 @happy
  Scenario: Error type mapping provides I/O and context error mapping
    Given an open NovusPack package
    And a valid context
    When I/O and context sentinel errors occur
    Then ErrIOError maps to ErrTypeIO
    And ErrContextCancelled maps to ErrTypeContext
    And ErrContextTimeout maps to ErrTypeContext

  @REQ-FILEMGMT-114 @happy
  Scenario: Error type mapping enables legacy error compatibility
    Given an open NovusPack package
    And a valid context
    When legacy code uses sentinel errors
    Then sentinel errors are supported
    And sentinel errors can be converted to structured errors
    And backward compatibility is maintained
