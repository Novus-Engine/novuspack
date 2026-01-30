@domain:basic_ops @m2 @REQ-API_BASIC-218 @spec(api_basic_operations.md#7211-read-only-enforcement-mechanism)
Feature: Read-only enforcement mechanism

  @REQ-API_BASIC-218 @happy
  Scenario: Read-only enforcement uses a wrapper without duplicating parsing logic
    Given a package opened in read-only mode
    When read-only enforcement is applied
    Then enforcement is implemented via a wrapper around a package instance
    And the wrapper avoids duplicating parsing logic
    And write operations are blocked with structured errors
    And read operations reuse the underlying read logic
    And enforcement behavior aligns with the OpenPackageReadOnly design

