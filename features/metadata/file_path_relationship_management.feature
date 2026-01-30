@domain:metadata @m2 @REQ-META-102 @spec(api_metadata.md#85-file-path-association)
Feature: File-Path Relationship Management

  @REQ-META-102 @happy
  Scenario: File-path association provides file-path relationships
    Given a NovusPack package
    When file-path association is used
    Then files can be associated with path metadata
    And path associations enable per-path tag inheritance
    And file-path relationships are managed

  @REQ-META-102 @happy
  Scenario: AssociateWithPathMetadata links files to path metadata
    Given a NovusPack package
    And a valid context
    And a FileEntry with paths
    And a PathMetadataEntry
    When AssociateWithPathMetadata is called
    Then file is linked to path metadata entry
    And association enables per-path tag inheritance
    And context supports cancellation

  @REQ-META-102 @happy
  Scenario: Path metadata association enables per-path tag inheritance
    Given a NovusPack package
    And a valid context
    And a FileEntry with multiple paths
    And PathMetadataEntry instances with different inheritance chains
    When GetEffectiveTags is called on each PathMetadataEntry
    Then each path can have different inherited tags
    And inheritance is resolved per-path via PathMetadataEntry.ParentPath
    And context supports cancellation

  @REQ-META-102 @happy
  Scenario: GetPathMetadataForPath retrieves path metadata for specific path
    Given a NovusPack package
    And a valid context
    And a FileEntry with path metadata associations
    When GetPathMetadataForPath is called with a path
    Then PathMetadataEntry for that path is returned
    And returned entry enables inheritance resolution via ParentPath
    And context supports cancellation

  @REQ-META-102 @error
  Scenario: File-path association handles invalid paths
    Given a NovusPack package
    When invalid file or path metadata reference is provided
    Then appropriate error is returned
    And error indicates path problem
    And error follows structured error format
