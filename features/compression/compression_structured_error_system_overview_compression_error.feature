@domain:compression @m2 @REQ-COMPR-091 @spec(api_package_compression.md#14-structured-error-system)
Feature: Compression: Structured Error System (Overview)

  @REQ-COMPR-091 @happy
  Scenario: Compression API uses comprehensive structured error system
    Given a compression operation
    When error occurs during compression
    Then structured error system is used
    And error provides better error categorization
    And error provides context information
    And error provides debugging capabilities

  @REQ-COMPR-091 @happy
  Scenario: Structured error system provides compression error types
    Given a compression operation
    When error occurs
    Then compression-specific error types are provided
    And error types include ErrTypeCompression, ErrTypeValidation, ErrTypeIO
    And error types include ErrTypeContext, ErrTypeCorruption, ErrTypeUnsupported

  @REQ-COMPR-091 @happy
  Scenario: Structured error system links to core error system documentation
    Given a compression operation
    And structured error system is used
    When error system is examined
    Then complete error system documentation is referenced
    And error system follows core structured error patterns
    And error system is consistent across API
