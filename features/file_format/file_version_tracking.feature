@domain:file_format @m2 @REQ-FILEFMT-049 @spec(package_file_format.md#4122-file-version-tracking)
Feature: File Version Tracking

  @REQ-FILEFMT-049 @happy
  Scenario: File version tracking provides version management
    Given a file entry
    When file version fields are examined
    Then FileVersion tracks file content changes
    And MetadataVersion tracks file metadata changes
    And dual versioning enables granular change detection

  @REQ-FILEFMT-049 @happy
  Scenario: FileVersion tracks file content changes independently
    Given a file entry
    And FileVersion has current value
    When file content is modified
    Then FileVersion is incremented
    And MetadataVersion remains unchanged
    And version change indicates file data modification

  @REQ-FILEFMT-049 @happy
  Scenario: MetadataVersion tracks file metadata changes independently
    Given a file entry
    And MetadataVersion has current value
    When file metadata is modified (paths, tags, compression, encryption)
    Then MetadataVersion is incremented
    And FileVersion remains unchanged
    And version change indicates metadata modification

  @REQ-FILEFMT-049 @happy
  Scenario: Dual versioning enables granular change detection
    Given a file entry
    And file has version history
    When version fields are compared
    Then FileVersion enables content change detection
    And MetadataVersion enables metadata change detection
    And granular change tracking is supported

  @REQ-FILEFMT-049 @happy
  Scenario: File version fields have initial value of 1 for new files
    Given a new file entry
    When file version fields are examined
    Then FileVersion is set to 1
    And MetadataVersion is set to 1
    And initial versions indicate new file entry

  @REQ-FILEFMT-049 @happy
  Scenario: File version tracking supports incremental operations
    Given a file entry
    And file has version history
    When incremental operations use version fields
    Then version fields enable efficient incremental operations
    And change detection uses version comparison
    And conflict resolution uses version information

  @REQ-FILEFMT-049 @happy
  Scenario: Package-level metadata is tracked separately
    Given a file entry
    And package metadata is modified
    When file version fields are examined
    Then FileVersion and MetadataVersion are unchanged
    And package-level metadata is tracked by package MetadataVersion
    And file metadata versions track only file-level changes
