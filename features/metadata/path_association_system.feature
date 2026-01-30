@domain:metadata @m2 @REQ-META-097 @spec(api_metadata.md#84-path-association-system)
Feature: Path Association System

  @REQ-META-097 @happy
  Scenario: Path association system links FileEntry to PathMetadataEntry
    Given a NovusPack package
    And a valid context
    When path association system is used
    Then FileEntry objects are linked to PathMetadataEntry metadata
    And tag inheritance is enabled
    And filesystem property management is enabled
    And context supports cancellation

  @REQ-META-097 @happy
  Scenario: AssociateWithPathMetadata links files to path metadata
    Given a NovusPack package
    And a valid context
    And a FileEntry with paths
    And a PathMetadataEntry
    When AssociateWithPathMetadata is called
    Then file entry is associated with path metadata entry
    And association is stored in FileEntry.PathMetadataEntries map
    And path matching is performed between FileEntry.Paths and PathMetadataEntry.Path
    And context supports cancellation

  @REQ-META-097 @happy
  Scenario: UpdateFilePathAssociations rebuilds all associations
    Given a NovusPack package
    And a valid context
    When UpdateFilePathAssociations is called
    Then all files are retrieved
    And all path metadata entries are retrieved
    And path map is built
    And each file path is associated with its path metadata entry
    And context supports cancellation

  @REQ-META-097 @error
  Scenario: Path association system handles errors
    Given a NovusPack package
    When invalid file or path metadata references are provided
    Then appropriate errors are returned
    And errors follow structured error format
