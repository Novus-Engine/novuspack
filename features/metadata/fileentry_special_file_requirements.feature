@domain:metadata @m2 @REQ-META-095 @REQ-META-099 @spec(api_metadata.md#8314-fileentry-requirements)
Feature: FileEntry Special File Requirements

  @REQ-META-095 @happy
  Scenario: Special file entry has correct Type field
    Given an open NovusPack package
    And a special metadata file entry
    When FileEntry Type field is examined
    Then Type is set to appropriate special file type
    And Type is 65001 for path metadata
    And Type is 65000 for package metadata

  @REQ-META-095 @happy
  Scenario: Special file entry has no compression
    Given an open NovusPack package
    And a special metadata file entry
    When FileEntry CompressionType field is examined
    Then CompressionType is set to 0
    And file is uncompressed for FastWrite compatibility

  @REQ-META-095 @happy
  Scenario: Special file entry has no encryption
    Given an open NovusPack package
    And a special metadata file entry
    When FileEntry EncryptionType field is examined
    Then EncryptionType is set to 0x00
    And special file is not encrypted

  @REQ-META-095 @happy
  Scenario: Special file entry has required tags
    Given an open NovusPack package
    And a special metadata file entry
    When FileEntry Tags are examined
    Then file_type tag is set to "special_metadata"
    And metadata_type tag indicates metadata type
    And tags identify special file purpose

  @REQ-META-099 @happy
  Scenario: FileEntry path properties link to PathMetadataEntry
    Given an open NovusPack package
    And a file entry with path metadata association
    When FileEntry path properties are examined
    Then PathMetadataEntries map contains path to PathMetadataEntry mappings
    And GetPathMetadataForPath retrieves PathMetadataEntry for specific paths
    And path hierarchy is accessible via PathMetadataEntry.ParentPath

  @REQ-META-099 @happy
  Scenario: PathMetadataEntry inherits tags from path hierarchy
    Given an open NovusPack package
    And a file entry with path metadata association
    And path hierarchy with tags
    When GetInheritedTags is called on PathMetadataEntry
    Then inherited tags contain tags from path hierarchy
    And tags are inherited based on priority
    And inheritance follows path hierarchy via ParentPath

  @REQ-META-099 @happy
  Scenario: PathMetadataEntry effective tags combine path and inherited tags
    Given an open NovusPack package
    And a file entry with own tags
    And path hierarchy with tags
    When GetEffectiveTags is called on PathMetadataEntry
    Then effective tags include PathMetadataEntry tags
    And effective tags include inherited tags from path hierarchy
    And effective tags include associated FileEntry tags
    And path-specific tags override inherited tags

  @REQ-META-011 @error
  Scenario: Special file entry validation fails with invalid Type
    Given an open NovusPack package
    And a file entry with invalid special file type
    When special file entry is validated
    Then structured validation error is returned
    And error indicates invalid file type

  @REQ-META-011 @error
  Scenario: Special file entry validation fails if compressed
    Given an open NovusPack package
    And a special metadata file entry with compression
    When special file entry is validated
    Then structured validation error is returned
    And error indicates compression not allowed
