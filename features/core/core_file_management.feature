@domain:core @m2 @REQ-CORE-032 @spec(api_core.md#4-file-management)
Feature: Core File Management

  @REQ-CORE-032 @happy
  Scenario: File management defines file operation capabilities
    Given an open NovusPack package
    And a valid context
    When file management is used
    Then file operation capabilities are available
    And files can be added, removed, and managed
    And file operations integrate with core interface

  @REQ-CORE-032 @happy
  Scenario: File management links to File Management API
    Given an open NovusPack package
    And file operations are needed
    When file management is examined
    Then file management references File Management API documentation
    And detailed file operation methods are documented
    And file encryption and deduplication are supported

  @REQ-CORE-032 @happy
  Scenario: File management supports basic file operations
    Given an open NovusPack package
    And a valid context
    When basic file operations are performed
    Then files can be added to the package
    And files can be removed from the package
    And files can be extracted from the package

  @REQ-CORE-032 @happy
  Scenario: File management supports encryption-aware operations
    Given an open NovusPack package
    And a valid context
    When encryption-aware file operations are performed
    Then files can be added with specific encryption types
    And encryption type system is available
    And encryption algorithms can be validated

  @REQ-CORE-032 @happy
  Scenario: File management supports pattern operations
    Given an open NovusPack package
    And a valid context
    When pattern operations are performed
    Then multiple files can be added using patterns
    And pattern matching is supported
    And bulk file operations are available

  @REQ-CORE-032 @happy
  Scenario: File management supports file information and queries
    Given an open NovusPack package
    And a valid context
    When file information is queried
    Then file information can be retrieved
    And file search capabilities are available
    And file existence can be checked

  @REQ-CORE-032 @happy
  Scenario: File management integrates with context
    Given an open NovusPack package
    And a valid context
    When file management operations are performed
    Then all methods accept context.Context
    And context cancellation is supported
    And context timeout handling is supported
