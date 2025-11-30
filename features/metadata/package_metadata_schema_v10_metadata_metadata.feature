@domain:metadata @m2 @REQ-META-040 @spec(metadata.md#221-package-metadata-schema-v10)
Feature: Package Metadata Schema v1.0

  @REQ-META-040 @happy
  Scenario: Package metadata schema v1.0 defines package metadata structure
    Given a NovusPack package
    When package metadata schema v1.0 is examined
    Then schema defines Package Information fields
    And schema defines Game-Specific Metadata fields
    And schema defines Asset Metadata fields
    And schema defines Security Metadata fields
    And schema defines Custom Metadata field

  @REQ-META-040 @happy
  Scenario: Package metadata schema v1.0 Package Information fields
    Given a NovusPack package
    And package metadata schema v1.0
    When Package Information fields are examined
    Then name field is string type
    And version field is string type
    And description field is string type
    And author field is string type
    And license field is string type
    And created field is ISO8601-timestamp type
    And modified field is ISO8601-timestamp type

  @REQ-META-040 @happy
  Scenario: Package metadata schema v1.0 includes Game-Specific and Asset Metadata
    Given a NovusPack package
    And package metadata schema v1.0
    When Game-Specific and Asset Metadata fields are examined
    Then Game-Specific Metadata includes engine, platform, genre, rating, requirements
    And Asset Metadata includes textures, sounds, models, scripts, total_size integer fields

  @REQ-META-040 @happy
  Scenario: Package metadata schema v1.0 includes Security and Custom Metadata
    Given a NovusPack package
    And package metadata schema v1.0
    When Security and Custom Metadata fields are examined
    Then Security Metadata includes encryption_level, signature_type string fields
    And Security Metadata includes security_scan, trusted_source boolean fields
    And Custom Metadata includes extensible key-value pairs object

  @REQ-META-040 @error
  Scenario: Package metadata schema v1.0 validates field types
    Given a NovusPack package
    When invalid schema field types are provided
    Then schema validation detects type mismatches
    And appropriate errors are returned
