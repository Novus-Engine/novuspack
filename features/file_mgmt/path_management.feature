@domain:file_mgmt @m2 @REQ-FILEMGMT-046 @spec(api_file_mgmt_file_entry.md#5-path-management)
Feature: Path Management

  @REQ-FILEMGMT-046 @happy
  Scenario: Path management methods support path operations and path metadata associations
    Given a file entry
    When path management methods are used
    Then path operations are supported
    And path metadata associations are supported
    And methods provide comprehensive path management

  @REQ-FILEMGMT-046 @happy
  Scenario: GetPrimaryPath returns primary path from file entry
    Given a file entry with primary path "documents/file.txt"
    When GetPrimaryPath is called
    Then "documents/file.txt" is returned
    And primary path is correctly retrieved

  @REQ-FILEMGMT-046 @happy
  Scenario: Path management methods support symlink operations
    Given a file entry with symlinks
    When HasSymlinks is called
    Then true is returned if symlinks exist
    And GetSymlinkPaths returns symlink path entries
    And ResolveAllSymlinks returns resolved paths

  @REQ-FILEMGMT-046 @happy
  Scenario: Path management methods support path metadata associations
    Given a file entry with paths
    And a PathMetadataEntry
    When AssociateWithPathMetadata is called
    Then path metadata association is set
    And GetPathMetadataForPath returns the PathMetadataEntry for a path
    And path metadata association enables per-path tag inheritance

  @REQ-FILEMGMT-046 @happy
  Scenario: Path metadata association enables per-path tag inheritance
    Given a file entry with multiple paths
    And PathMetadataEntry instances with different inheritance chains
    When GetEffectiveTags is called on each PathMetadataEntry
    Then each path can have different inherited tags
    And inheritance is resolved per-path via PathMetadataEntry.ParentPath
    And path hierarchy is correctly represented

  @REQ-FILEMGMT-046 @happy
  Scenario: Path metadata association provides per-path tag access
    Given a file entry with path metadata associations
    When GetPathMetadataForPath is called for a specific path
    Then PathMetadataEntry for that path is returned
    And GetInheritedTags can be called on the PathMetadataEntry
    And GetEffectiveTags includes FileEntry tags and inherited path tags

  @REQ-FILEMGMT-046 @error
  Scenario: Path management methods handle invalid path metadata associations
    Given a file entry
    And invalid PathMetadataEntry
    When AssociateWithPathMetadata is called with invalid entry
    Then structured error is returned
    And error indicates invalid association
    And error follows structured error format
