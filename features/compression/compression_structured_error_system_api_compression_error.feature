@domain:compression @m2 @REQ-COMPR-092 @spec(api_package_compression.md#141-structured-error-system)
Feature: Compression: Structured Error System (API)

  @REQ-COMPR-092 @happy
  Scenario: Compression API uses comprehensive structured error system
    Given a compression operation
    When error occurs during compression
    Then structured error system is used
    And error provides better error categorization
    And error provides context information
    And error provides debugging capabilities

  @REQ-COMPR-092 @happy
  Scenario: Structured error system links to core error system documentation
    Given a compression operation
    And structured error system is used
    When error system is examined
    Then complete error system documentation is referenced
    And error system follows core structured error patterns
    And error system is consistent across API

  @REQ-COMPR-092 @happy
  Scenario: Structured error system provides error categorization
    Given a compression operation
    When error occurs
    Then error is categorized by type
    And error type indicates error category
    And categorization enables appropriate error handling

  @REQ-COMPR-092 @happy
  Scenario: Structured error system provides error context
    Given a compression operation
    When error occurs
    Then error includes context information
    And context helps identify source of problem
    And context provides details about operation

  @REQ-COMPR-092 @happy
  Scenario: Structured error system provides debugging capabilities
    Given a compression operation
    When error occurs
    Then error provides debugging information
    And error details enable problem diagnosis
    And error system supports troubleshooting

  @REQ-COMPR-092 @error
  Scenario: Structured error system handles compression operation failures
    Given a compression operation
    And compression operation fails
    When error is returned
    Then structured error is provided
    And error includes compression-specific context
    And error follows structured error format
