@domain:validation @REQ-FILEMGMT-366 @REQ-FILEMGMT-367 @REQ-FILEMGMT-368 @REQ-FILEMGMT-369 @REQ-FILEMGMT-370 @REQ-FILEMGMT-371 @REQ-FILEMGMT-372 @REQ-FILEMGMT-373 @REQ-FILEMGMT-374 @spec(api_metadata.md#828-symlink-validation-methods)
Feature: Symlink path validation

  As a package system
  I want to validate symlink paths
  So that symlinks are secure and point to valid targets

  @REQ-FILEMGMT-366 @REQ-FILEMGMT-370 @happy
  Scenario: Validate symlink paths are package-relative
    Given a package
    When I validate symlink paths "/data/symlink.txt" and "/data/target.txt"
    Then validation should succeed
    And both paths should be recognized as package-relative

  @REQ-FILEMGMT-366 @REQ-FILEMGMT-371 @REQ-FILEMGMT-373 @error
  Scenario: Reject symlink paths escaping package root
    Given a package
    When I validate symlink paths "/data/symlink.txt" and "../outside/target.txt"
    Then validation should fail with ErrTypeSecurity
    And the error message should indicate the path escapes package root

  @REQ-FILEMGMT-366 @REQ-FILEMGMT-371 @REQ-FILEMGMT-373 @error
  Scenario: Reject symlink source path escaping package root
    Given a package
    When I validate symlink paths "../outside/symlink.txt" and "/data/target.txt"
    Then validation should fail with ErrTypeSecurity
    And the error message should indicate the source path escapes package root

  @REQ-FILEMGMT-366 @REQ-FILEMGMT-372 @REQ-FILEMGMT-374 @error
  Scenario: Reject symlink with non-existent target
    Given a package
    And no file exists at "/data/target.txt"
    When I validate symlink paths "/data/symlink.txt" and "/data/target.txt"
    Then validation should fail with ErrTypeNotFound
    And the error message should indicate the target does not exist

  @REQ-FILEMGMT-367 @REQ-FILEMGMT-372 @happy
  Scenario: TargetExists returns true for FileEntry
    Given a package
    And a file entry exists at "/data/file.txt"
    When I check if target exists at "/data/file.txt"
    Then TargetExists should return true

  @REQ-FILEMGMT-367 @REQ-FILEMGMT-372 @happy
  Scenario: TargetExists returns true for directory PathMetadataEntry
    Given a package
    And a PathMetadataEntry directory exists at "/data/dir"
    When I check if target exists at "/data/dir"
    Then TargetExists should return true

  @REQ-FILEMGMT-367 @happy
  Scenario: TargetExists returns false for non-existent path
    Given a package
    And no file or directory exists at "/data/missing.txt"
    When I check if target exists at "/data/missing.txt"
    Then TargetExists should return false

  @REQ-FILEMGMT-368 @happy
  Scenario: ValidatePathWithinPackageRoot normalizes valid path
    Given a package
    When I validate path "/data/file.txt"
    Then validation should succeed
    And the normalized path should be returned

  @REQ-FILEMGMT-368 @REQ-FILEMGMT-373 @error
  Scenario: ValidatePathWithinPackageRoot rejects path with ".." escaping root
    Given a package
    When I validate path "/data/../../outside/file.txt"
    Then validation should fail with ErrTypeSecurity
    And the error message should indicate the path escapes package root

  @REQ-FILEMGMT-369 @happy
  Scenario: AddSymlink validates paths before adding
    Given a package
    And a file entry exists at "/data/target.txt"
    When I add a symlink from "/data/symlink.txt" to "/data/target.txt"
    Then validation should be performed
    And the symlink should be added successfully

  @REQ-FILEMGMT-369 @error
  Scenario: AddSymlink rejects symlink with invalid paths
    Given a package
    When I attempt to add a symlink with paths escaping package root
    Then AddSymlink should fail with ErrTypeSecurity
    And the symlink should not be added

  @REQ-FILEMGMT-369 @error
  Scenario: AddSymlink rejects symlink with non-existent target
    Given a package
    And no file exists at "/data/target.txt"
    When I attempt to add a symlink from "/data/symlink.txt" to "/data/target.txt"
    Then AddSymlink should fail with ErrTypeNotFound
    And the symlink should not be added
