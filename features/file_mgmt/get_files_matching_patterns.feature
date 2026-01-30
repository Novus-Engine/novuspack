@domain:file_mgmt @m2 @REQ-FILEMGMT-027 @spec(api_file_mgmt_addition.md#24-addfilepattern-package-method)
Feature: Get files matching patterns

  @happy
  Scenario: GetPatterns gets files matching patterns
    Given an open package with files matching patterns
    When GetPatterns is called with pattern
    Then file entries matching pattern are returned
    And pattern matching works correctly
    And all matching files are included

  @happy
  Scenario: GetPatterns returns empty list when no matches
    Given an open package with files
    When GetPatterns is called with pattern that matches nothing
    Then empty list is returned

  @error
  Scenario: GetPatterns handles invalid patterns gracefully
    Given an open package
    When GetPatterns is called with invalid pattern
    Then structured validation error is returned
    And error indicates invalid pattern

  @REQ-FILEMGMT-037 @REQ-FILEMGMT-041 @error
  Scenario: GetPatterns respects context cancellation
    Given an open package with files
    And a cancelled context
    When GetPatterns is called
    Then structured context error is returned
    And error type is context cancellation
