@domain:basic_ops @m2 @REQ-API_BASIC-028 @spec(api_basic_operations.md#412-newpackage-example-usage)
Feature: NewPackage Example Usage

  @REQ-API_BASIC-028 @happy
  Scenario: NewPackage example demonstrates package instance creation
    Given a context for package operations
    When NewPackage is called
    Then Package instance is created
    And instance is ready for Create or Open
    And example demonstrates constructor usage

  @REQ-API_BASIC-028 @happy
  Scenario: NewPackage example shows typical workflow
    Given NewPackage creates package instance
    When Create or Open is called next
    Then package workflow is demonstrated
    And example shows proper initialization pattern
    And usage pattern is clear

  @REQ-API_BASIC-028 @happy
  Scenario: NewPackage example follows standard pattern with defer
    Given a code example using NewPackage
    When example code is examined
    Then NewPackage is called with context
    And defer package.Close() is used
    And example demonstrates resource cleanup pattern

  @REQ-API_BASIC-028 @happy
  Scenario: NewPackage example demonstrates in-memory package creation
    Given a code example using NewPackage
    When example code is examined
    Then package is created in memory
    And no file I/O occurs during creation
    And package must be written to persist

  @REQ-API_BASIC-028 @happy
  Scenario: NewPackage example shows simple initialization
    Given a code example using NewPackage
    When example code is examined
    Then example shows simple one-line package creation
    And example demonstrates minimal setup required
    And example is clear and concise
