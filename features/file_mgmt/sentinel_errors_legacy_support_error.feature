@domain:file_mgmt @m2 @REQ-FILEMGMT-113 @spec(api_file_management.md#1221-sentinel-errors-legacy-support)
Feature: Sentinel Errors Legacy Support

  @REQ-FILEMGMT-113 @happy
  Scenario: Sentinel errors provide legacy error support
    Given an open NovusPack package
    And a valid context
    When legacy code uses sentinel errors
    Then sentinel errors are supported
    And legacy error compatibility is maintained
    And backward compatibility is preserved

  @REQ-FILEMGMT-113 @happy
  Scenario: Sentinel errors are defined for common error conditions
    Given an open NovusPack package
    And a valid context
    When common error conditions occur
    Then ErrFileNotFound sentinel error is available
    And ErrFileExists sentinel error is available
    And ErrInvalidPath sentinel error is available
    And ErrPackageNotOpen sentinel error is available
    And other common sentinel errors are available

  @REQ-FILEMGMT-113 @happy
  Scenario: Sentinel errors can be converted to structured errors
    Given an open NovusPack package
    And a valid context
    When sentinel errors are used
    Then sentinel errors can be converted to structured errors
    And error type mapping is available
    And structured error conversion is supported
