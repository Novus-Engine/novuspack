@skip @domain:metadata @m2 @spec(metadata.md#package-information)
Feature: Package Information Metadata

# This feature captures high-level package information expectations from the metadata spec.
# Detailed runnable scenarios live in the dedicated metadata feature files.

  @documentation
  Scenario: Package information includes human-readable identity fields
    Given a package metadata file is present
    When package information is represented in metadata
    Then the metadata includes name, version, description, author, and license fields
    And timestamps are represented using an ISO8601 timestamp format

  @documentation
  Scenario: Package information may include security and custom metadata
    Given a package needs additional metadata beyond identity fields
    When the metadata schema is populated
    Then security metadata may include encryption level and trust signals
    And a custom metadata object may be used for extensible key-value pairs
