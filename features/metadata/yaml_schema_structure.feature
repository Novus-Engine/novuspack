@domain:metadata @m2 @REQ-META-039 @spec(metadata.md#22-yaml-schema-structure)
Feature: YAML Schema Structure

  @REQ-META-039 @happy
  Scenario: YAML schema structure defines package metadata schema v1.0
    Given a NovusPack package
    When YAML schema structure is examined
    Then Package Information fields are defined
    And Game-Specific Metadata fields are defined
    And Asset Metadata fields are defined
    And Security Metadata fields are defined
    And Custom Metadata field is defined

  @REQ-META-039 @happy
  Scenario: Package Information fields in schema
    Given a NovusPack package
    When Package Information fields are examined
    Then name field is string type
    And version field is string type
    And description field is string type
    And author field is string type
    And license field is string type
    And created field is ISO8601-timestamp type
    And modified field is ISO8601-timestamp type

  @REQ-META-039 @happy
  Scenario: Game-Specific Metadata fields in schema
    Given a NovusPack package
    When Game-Specific Metadata fields are examined
    Then engine field is string type
    And platform field is array of strings type
    And genre field is string type
    And rating field is string type
    And requirements object contains min_ram, min_storage, graphics, and os fields

  @REQ-META-039 @happy
  Scenario: Asset and Security Metadata fields in schema
    Given a NovusPack package
    When Asset and Security Metadata fields are examined
    Then Asset Metadata contains textures, sounds, models, scripts, and total_size integer fields
    And Security Metadata contains encryption_level, signature_type string fields
    And Security Metadata contains security_scan, trusted_source boolean fields

  @REQ-META-039 @error
  Scenario: YAML schema structure validates field types
    Given a NovusPack package
    When invalid field types are provided
    Then schema validation detects type mismatches
    And appropriate errors are returned
