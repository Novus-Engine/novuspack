@skip @domain:security @m2 @spec(api_security.md#322-example-usage)
Feature: Security Usage Examples

# This feature captures usage-oriented security scenarios derived from the security specs.
# Detailed runnable scenarios live in the dedicated security feature files.

  @documentation
  Scenario: Security status reporting summarizes checksums and validation errors
    Given an open package
    When the caller requests a security status report
    Then the report indicates whether checksums are present and valid
    And the report includes any validation errors as descriptive strings

  @documentation
  Scenario: Security status reporting is v1-compatible when signature validation is deferred
    Given signature validation is deferred to v2
    When the caller requests a security validation result in v1
    Then signature-related fields are returned as zero values
    And checksum-related fields are still populated when available
