@skip @domain:compression @m2 @spec(api_package_compression.md#12-error-handling)
Feature: Compression Error Handling

# This feature captures high-level expectations for compression error behavior and recovery.
# More detailed runnable scenarios live in dedicated compression error handling feature files.

  @REQ-COMPR-110 @error
  Scenario: Compression failures return structured compression errors
    Given a compression operation fails while processing package content
    When the operation returns an error
    Then the error is a structured error
    And the error category indicates a compression failure

  @REQ-COMPR-060 @error
  Scenario: Compression operations support recovery and cleanup on error
    Given a compression operation fails partway through processing
    When the operation aborts due to the failure
    Then partial state is cleaned up
    And the package remains in a consistent state for subsequent operations
