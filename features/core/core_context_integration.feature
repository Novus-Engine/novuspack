@domain:core @m2 @REQ-CORE-015 @REQ-CORE-016 @REQ-CORE-017 @REQ-FILEMGMT-037 @REQ-FILEMGMT-041 @REQ-WRITE-008 @REQ-WRITE-011 @REQ-COMPR-016 @REQ-SIG-015 @REQ-SIG-019 @REQ-STREAM-013 @REQ-STREAM-017 @REQ-DEDUP-003 @REQ-DEDUP-005 @REQ-SEC-007 @REQ-SEC-011 @REQ-GEN-005 @REQ-META-011 @REQ-META-014 @REQ-VALID-003 @spec(api_core.md#02-context-integration)
Feature: Core Context Integration

  @REQ-CORE-015 @happy
  Scenario: All public methods accept context.Context as first parameter
    Given a NovusPack API operation
    When public methods are called
    Then all methods accept context.Context as first parameter
    And context parameter is required
    And context follows standard Go patterns

  @REQ-CORE-016 @happy
  Scenario: Context cancellation must be checked and respected in all operations
    Given a NovusPack API operation
    And a context with cancellation support
    When operation is performed
    Then context cancellation is checked during operation
    And cancellation is respected
    And operation can be cancelled via context

  @REQ-CORE-017 @happy
  Scenario: Context timeout errors returned as structured context errors
    Given a NovusPack API operation
    And a context that times out
    When operation is performed
    Then context timeout error is returned
    And error is structured context error
    And error follows structured error format

  @REQ-CORE-015 @REQ-CORE-016 @REQ-CORE-017 @happy
  Scenario: Context integration supports request cancellation and timeout handling
    Given a NovusPack API operation
    When context integration is used
    Then request cancellation is supported
    And timeout handling is supported
    And graceful shutdown is enabled

  @REQ-CORE-015 @REQ-CORE-016 @REQ-CORE-017 @happy
  Scenario: Context integration follows 2025 Go best practices
    Given a NovusPack API operation
    When context integration is examined
    Then API is compatible with modern Go applications
    And API follows Go standard context patterns
    And API integrates with Go frameworks
