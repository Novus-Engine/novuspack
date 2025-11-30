@domain:file_mgmt @m2 @REQ-FILEMGMT-046 @spec(api_file_management.md#123-path-and-directory-management)
Feature: Path and Directory Management

  @REQ-FILEMGMT-046 @happy
  Scenario: Path and directory management methods support path operations and directory associations
    Given a file entry
    When path and directory management methods are used
    Then path operations are supported
    And directory associations are supported
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
  Scenario: Path management methods support directory associations
    Given a file entry
    When SetParentDirectory is called with directory entry
    Then parent directory is set
    And GetParentDirectory returns the parent directory
    And GetParentPath returns parent directory path

  @REQ-FILEMGMT-046 @happy
  Scenario: Path management methods support directory depth and hierarchy
    Given a file entry in nested directory structure
    When GetDirectoryDepth is called
    Then directory depth is returned
    And GetAncestorDirectories returns ancestor directories
    And directory hierarchy is correctly represented

  @REQ-FILEMGMT-046 @happy
  Scenario: Path management methods support inherited tags
    Given a file entry with directory associations
    When SetInheritedTags is called
    Then inherited tags are set
    And GetInheritedTags returns inherited tags
    And UpdateInheritedTags updates inherited tags

  @REQ-FILEMGMT-046 @happy
  Scenario: Path management methods respect context
    Given a file entry
    And a valid context
    When path management methods are called
    Then context supports cancellation
    And context supports timeout handling

  @REQ-FILEMGMT-046 @error
  Scenario: Path management methods handle invalid directory associations
    Given a file entry
    And invalid directory entry
    When SetParentDirectory is called with invalid entry
    Then structured error is returned
    And error indicates invalid association
    And error follows structured error format
