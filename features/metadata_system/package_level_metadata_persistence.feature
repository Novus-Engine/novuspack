@domain:metadata_system @m2 @REQ-METASYS-002 @spec(metadata.md#2-package-metadata-file-specification)
Feature: Package-level metadata persistence

  @happy
  Scenario: Package metadata fields are persisted and validated
    Given an open package
    When I set package-level metadata fields
    Then fields should be persisted and validated per schema

  @happy
  Scenario: Package metadata schema is defined
    Given package metadata system
    When metadata schema is examined
    Then schema defines required fields
    And schema defines optional fields
    And schema defines validation rules

  @happy
  Scenario: Package metadata is stored in special files
    Given a package
    When package metadata is set
    Then metadata is stored in special metadata files
    And file types 65000-65535 are used
    And metadata is accessible

  @error
  Scenario: Invalid package metadata is rejected
    Given a package
    When invalid metadata is set
    Then structured validation error is returned
    And error indicates schema violation
