@domain:file_format @m2 @REQ-FILEFMT-047 @spec(package_file_format.md#412-file-entry-structure-requirements)
Feature: File Entry Structure Requirements

  @REQ-FILEFMT-047 @happy
  Scenario: File entry structure requirements define entry format rules
    Given a file entry
    When file entry structure requirements are examined
    Then entry format rules are defined
    And structure supports unique file identification
    And structure supports version tracking

  @REQ-FILEFMT-047 @happy
  Scenario: File entry structure supports unique file identification
    Given a file entry
    When structure requirements are examined
    Then structure includes unique 64-bit FileID
    And FileID provides stable identification
    And FileID enables efficient file tracking

  @REQ-FILEFMT-047 @happy
  Scenario: File entry structure supports version tracking
    Given a file entry
    When structure requirements are examined
    Then structure includes FileVersion and MetadataVersion
    And dual versioning tracks content and metadata changes independently
    And version tracking enables granular change detection

  @REQ-FILEFMT-047 @happy
  Scenario: File entry structure supports multiple paths with per-path metadata
    Given a file entry
    When structure requirements are examined
    Then structure supports multiple paths pointing to same content
    And each path can have its own metadata (permissions, timestamps)
    And multiple path support enables path aliasing

  @REQ-FILEFMT-047 @happy
  Scenario: File entry structure supports hash-based content identification
    Given a file entry
    When structure requirements are examined
    Then structure includes multiple hash types for different purposes
    And content hashes support deduplication
    And integrity hashes support verification
    And fast lookup hashes support quick identification

  @REQ-FILEFMT-047 @happy
  Scenario: File entry structure supports security metadata
    Given a file entry
    When structure requirements are examined
    Then structure includes encryption and compression metadata
    And per-file security and optimization settings are supported
    And security metadata enables per-file configuration
