@domain:metadata @m2 @REQ-META-099 @spec(api_metadata.md#841-fileentry-directory-properties)
Feature: FileEntry Directory Properties

  @REQ-META-099 @happy
  Scenario: FileEntry directory properties provide directory association
    Given a NovusPack package
    And a FileEntry
    When FileEntry directory properties are examined
    Then DirectoryEntry property points to directory metadata
    And ParentDirectory property points to parent directory metadata
    And InheritedTags property contains cached inherited tags

  @REQ-META-099 @happy
  Scenario: DirectoryEntry property links file to directory
    Given a NovusPack package
    And a FileEntry
    When DirectoryEntry property is set
    Then property points to directory metadata for file's immediate directory
    And directory association enables tag inheritance
    And directory association enables filesystem property management

  @REQ-META-099 @happy
  Scenario: ParentDirectory property enables inheritance resolution
    Given a NovusPack package
    And a FileEntry
    When ParentDirectory property is set
    Then property points to parent directory metadata
    And property enables inheritance resolution
    And property supports tag inheritance hierarchy

  @REQ-META-099 @happy
  Scenario: InheritedTags property caches inherited tags
    Given a NovusPack package
    And a FileEntry
    When InheritedTags property is populated
    Then property contains cached inherited tags from directory hierarchy
    And cached tags improve performance
    And tags are resolved from directory metadata

  @REQ-META-099 @error
  Scenario: FileEntry directory properties handle invalid associations
    Given a NovusPack package
    When invalid directory associations are provided
    Then association validation detects invalid references
    And appropriate errors are returned
