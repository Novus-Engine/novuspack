@domain:dedup @REQ-DEDUP-017 @REQ-DEDUP-018 @REQ-FILEMGMT-362 @spec(api_deduplication.md#122-pathhandling-integration) @spec(api_basic_operations.md#46-package-configuration)
Feature: Deduplication with automatic symlink creation

  As a package creator
  I want deduplication to automatically create symlinks
  So that duplicate content is stored efficiently without manual conversion

  @REQ-DEDUP-018 @REQ-FILEMGMT-362 @happy
  Scenario: Automatic conversion during deduplication
    Given a package with AutoConvertToSymlinks enabled
    And an existing file entry with path "/data/file.txt"
    When I add a duplicate file with path "/backup/file.txt"
    Then the deduplication should create a symlink instead of adding a path
    And the file entry should have 1 path "/data/file.txt"
    And a symlink should be created from "/backup/file.txt" to "/data/file.txt"

  @REQ-DEDUP-017 @happy
  Scenario: PathHandlingSymlinks creates symlinks during deduplication
    Given a package with DefaultPathHandling set to PathHandlingSymlinks
    And an existing file entry with path "/data/file.txt"
    When I add a duplicate file with PathHandling set to PathHandlingSymlinks
    Then the deduplication should create a symlink
    And the file entry should have 1 path
    And a symlink should point to the primary path

  @REQ-DEDUP-017 @happy
  Scenario: PathHandlingHardLinks adds paths during deduplication
    Given a package with DefaultPathHandling set to PathHandlingHardLinks
    And an existing file entry with path "/data/file.txt"
    When I add a duplicate file with PathHandling set to PathHandlingHardLinks
    Then the deduplication should add the path to the existing FileEntry
    And the file entry should have 2 paths
    And no symlinks should be created

  @REQ-DEDUP-017 @happy
  Scenario: PathHandlingDefault uses package default
    Given a package with DefaultPathHandling set to PathHandlingSymlinks
    And an existing file entry with path "/data/file.txt"
    When I add a duplicate file with PathHandling set to PathHandlingDefault
    Then the deduplication should use the package default
    And a symlink should be created

  @REQ-DEDUP-017 @happy
  Scenario: PrimaryPathSelector used during automatic symlink creation
    Given a package with AutoConvertToSymlinks enabled
    And an existing file entry with paths "/long/path/file.txt" and "/short.txt"
    When I add a duplicate file with PrimaryPathSelector set to shortest path
    Then the primary path should be "/short.txt"
    And the symlink should point to "/short.txt"
