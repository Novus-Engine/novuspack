@domain:compression @m2 @REQ-COMPR-096 @spec(api_package_compression.md#1431-creating-compression-errors)
Feature: Creating Compression Errors

  @REQ-COMPR-096 @happy
  Scenario: Compression errors are created with rich context
    Given a compression operation that fails
    When structured compression error is created
    Then error includes ErrTypeCompression type
    And error includes algorithm context (e.g. Zstd)
    And error includes compression level context (e.g. 6)
    And error includes input size context
    And error includes operation context

  @REQ-COMPR-096 @happy
  Scenario: Compression errors wrap underlying errors
    Given a compression operation with underlying error
    When structured compression error is created
    Then original error is wrapped
    And error chain is maintained
    And error unwrapping is supported

  @REQ-COMPR-096 @happy
  Scenario: Compression errors provide debugging information
    Given compression errors with context
    When error is inspected
    Then error provides debugging information
    And context aids in problem diagnosis
    And error details enable troubleshooting

  @REQ-COMPR-096 @error
  Scenario: Creating compression errors demonstrates error creation pattern
    Given compression operations that may fail
    When compression errors are created
    Then error creation pattern is demonstrated
    And structured error system is used
    And context information is included
