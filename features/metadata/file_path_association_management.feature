@domain:metadata @m2 @REQ-META-011 @REQ-META-014 @REQ-META-098 @REQ-META-101 @REQ-META-135 @REQ-META-136 @REQ-META-156 @spec(api_metadata.md#841-fileentry-path-properties,api_metadata.md#843-association-management) @spec(api_metadata.md#path-association-methods) @spec(api_metadata.md#827-packageassociatefilewithpath-method) @spec(api_metadata.md#828-packageupdatefilepathassociations-method)
Feature: File-Path Association Management

  @REQ-META-098 @happy
  Scenario: FileEntry path properties link to path metadata
    Given an open NovusPack package
    And a file entry with path metadata association
    When FileEntry path properties are examined
    Then PathMetadataEntries map contains path to PathMetadataEntry mappings
    And GetPathMetadataForPath retrieves PathMetadataEntry for specific paths
    And path metadata association enables per-path tag inheritance

  @REQ-META-098 @happy
  Scenario: PathMetadataEntry filesystem properties provide filesystem information
    Given an open NovusPack package
    And a PathMetadataEntry with filesystem properties
    When PathMetadataEntry filesystem properties are examined
    Then Mode property is available for Unix/Linux permissions
    And UID and GID properties are available
    And ACL entries are available
    And WindowsAttrs property is available for Windows
    And ExtendedAttrs map is available
    And Flags property is available

  @REQ-META-101 @happy
  Scenario: AssociateWithPathMetadata links file to path metadata
    Given an open writable NovusPack package
    And a FileEntry with paths
    And a PathMetadataEntry
    And a valid context
    When AssociateWithPathMetadata is called
    Then file is linked to path metadata entry
    And association is stored in FileEntry.PathMetadataEntries map
    And path matching is performed between FileEntry.Paths and PathMetadataEntry.Path
    And association enables per-path tag inheritance

  @REQ-META-135 @happy
  Scenario: AssociateFileWithPath links FileEntry to PathMetadataEntry by stored path
    Given an open writable NovusPack package
    And a file exists at stored path "/docs/readme.txt"
    And a PathMetadataEntry exists for "/docs/readme.txt"
    And a valid context
    When AssociateFileWithPath is called with "/docs/readme.txt"
    Then FileEntry is associated with the matching PathMetadataEntry
    And PathMetadataEntry.ParentPath may be set for hierarchy traversal
    And missing parent paths do not cause this operation to fail

  @REQ-META-101 @happy
  Scenario: Path metadata association enables per-path tag inheritance
    Given an open writable NovusPack package
    And a FileEntry with multiple paths
    And PathMetadataEntry instances with different inheritance chains
    And a valid context
    When GetEffectiveTags is called on each PathMetadataEntry
    Then each path can have different inherited tags
    And inheritance is resolved per-path via PathMetadataEntry.ParentPath
    And FileEntry tags are included in effective tags for each path

  @REQ-META-101 @REQ-META-136 @happy
  Scenario: UpdateFilePathAssociations updates all file-path associations
    Given an open writable NovusPack package
    And multiple files with paths
    And path metadata entries
    And a valid context
    When UpdateFilePathAssociations is called
    Then all file-path associations are updated
    And files are linked to correct path metadata entries
    And associations match file paths
    And parent path associations are established when parents exist
    And missing parents do not cause this operation to fail

  @REQ-META-101 @happy
  Scenario: GetPathMetadataForPath retrieves path metadata for specific path
    Given an open NovusPack package
    And a FileEntry with associated path metadata
    And a valid context
    When GetPathMetadataForPath is called with a path
    Then PathMetadataEntry for that path is returned
    And returned entry enables inheritance resolution via ParentPath
    And all path associations are accessible

  @REQ-META-101 @happy
  Scenario: Association management maintains inheritance relationships
    Given an open writable NovusPack package
    And files with path metadata associations
    And path hierarchy with tags
    When associations are maintained
    Then paths inherit tags from parent paths via PathMetadataEntry.ParentPath
    And inheritance relationships are preserved
    And tag inheritance works correctly per path

  @REQ-META-011 @REQ-META-014 @error
  Scenario: Association management operations respect context cancellation
    Given an open writable NovusPack package
    And a cancelled context
    When association management operation is called
    Then structured context error is returned
    And error type is context cancellation

  @REQ-META-011 @error
  Scenario: AssociateWithPathMetadata fails if file does not exist
    Given an open writable NovusPack package
    And a non-existent FileEntry
    And a PathMetadataEntry
    And a valid context
    When AssociateWithPathMetadata is called
    Then structured validation error is returned
    And error indicates path does not match

  @REQ-META-011 @error
  Scenario: AssociateWithPathMetadata fails if path metadata does not exist
    Given an open writable NovusPack package
    And a FileEntry
    And a non-existent PathMetadataEntry
    And a valid context
    When AssociateWithPathMetadata is called
    Then structured validation error is returned
    And error indicates path does not match
