@domain:metadata @m2 @REQ-META-099 @spec(api_metadata.md#841-fileentry-path-properties)
Feature: FileEntry Path Properties

  @REQ-META-099 @happy
  Scenario: FileEntry path properties provide path metadata association
    Given a NovusPack package
    And a FileEntry with paths
    When FileEntry path properties are examined
    Then PathMetadataEntries map contains path to PathMetadataEntry mappings
    And each path in FileEntry.Paths can be associated with a PathMetadataEntry
    And path metadata association enables tag inheritance per path

  @REQ-META-099 @happy
  Scenario: PathMetadataEntries map links file paths to path metadata
    Given a NovusPack package
    And a FileEntry with multiple paths
    And PathMetadataEntry instances for those paths
    When AssociateWithPathMetadata is called for each path
    Then PathMetadataEntries map contains entries for each associated path
    And path association enables tag inheritance per path
    And path association enables filesystem property management per path

  @REQ-META-099 @happy
  Scenario: GetPathMetadataForPath retrieves path metadata for specific path
    Given a NovusPack package
    And a FileEntry with associated path metadata
    When GetPathMetadataForPath is called with a path
    Then PathMetadataEntry for that path is returned
    And returned entry enables inheritance resolution via ParentPath
    And returned entry supports tag inheritance hierarchy

  @REQ-META-099 @happy
  Scenario: Path metadata association enables per-path tag inheritance
    Given a NovusPack package
    And a FileEntry with multiple paths
    And PathMetadataEntry instances with different inheritance chains
    When GetEffectiveTags is called on each PathMetadataEntry
    Then each path can have different inherited tags
    And inheritance is resolved per-path via PathMetadataEntry.ParentPath
    And FileEntry tags are included in effective tags for each path

  @REQ-META-099 @error
  Scenario: FileEntry path properties handle invalid associations
    Given a NovusPack package
    When invalid path metadata associations are provided
    Then association validation detects invalid references
    And appropriate errors are returned
