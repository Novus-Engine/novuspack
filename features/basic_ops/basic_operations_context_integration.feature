@domain:basic_ops @m2 @REQ-API_BASIC-017 @spec(api_basic_operations.md#02-context-integration)
Feature: Basic Operations Context Integration

  @REQ-API_BASIC-017 @happy
  Scenario: All methods accept context.Context as first parameter
    Given a NovusPack package operation
    When method is called with context
    Then context is accepted as first parameter
    And operation uses context for cancellation and timeout
    And context integration follows Go best practices

  @REQ-API_BASIC-017 @happy
  Scenario: Context supports request cancellation
    Given a package operation with context
    And a cancellable context
    When context is cancelled during operation
    Then operation is cancelled gracefully
    And structured context error is returned
    And resources are cleaned up

  @REQ-API_BASIC-017 @happy
  Scenario: Context supports timeout handling
    Given a package operation with context
    And a context with timeout
    When operation exceeds timeout
    Then operation is cancelled
    And context timeout error is returned
    And timeout prevents indefinite blocking

  @REQ-API_BASIC-017 @happy
  Scenario: Context supports request-scoped values
    Given a package operation with context
    And context contains request-scoped values
    When operation accesses context values
    Then values are available during operation
    And values are request-scoped

  @REQ-API_BASIC-017 @error
  Scenario: Cancelled context returns structured error
    Given a package operation
    And a cancelled context
    When operation is called with cancelled context
    Then ErrTypeContext error is returned
    And error indicates context cancellation

  @REQ-API_BASIC-017 @error
  Scenario: Timeout context returns structured error
    Given a package operation
    And a context with expired timeout
    When operation is called with timed-out context
    Then ErrTypeContext error is returned
    And error indicates context timeout
