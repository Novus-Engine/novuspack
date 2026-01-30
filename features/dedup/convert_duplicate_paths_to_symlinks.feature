@domain:dedup @conversion @REQ-FILEMGMT-352 @REQ-FILEMGMT-353 @REQ-FILEMGMT-354 @REQ-FILEMGMT-355 @REQ-FILEMGMT-356 @REQ-FILEMGMT-357 @REQ-FILEMGMT-358 @REQ-FILEMGMT-363 @REQ-FILEMGMT-364 @REQ-FILEMGMT-365 @spec(api_file_mgmt_extraction.md#67-convert-paths-to-symlinks) @spec(api_file_mgmt_updates.md#17-symlinkconvertoptions-struct) @spec(api_file_mgmt_updates.md#171-convertpathstosymlinks-methods) @spec(api_file_mgmt_updates.md#175-convertpathstosymlinks-parameters) @spec(api_file_mgmt_updates.md#176-convertpathstosymlinks-returns) @spec(api_file_mgmt_updates.md#177-convertpathstosymlinks-behavior) @spec(api_file_mgmt_updates.md#178-convertpathstosymlinks-error-conditions) @spec(api_file_mgmt_updates.md#1710-convertallpathstosymlinks-behavior) @spec(api_file_mgmt_updates.md#1712-convertsymlinkstohardlinks-behavior) @spec(api_metadata.md#851-symlink-validation-methods)
Feature: Convert duplicate path entries to symlinks

  As a package creator
  I want to convert duplicate path entries into symlinks
  So that I can optimize storage and maintain proper filesystem semantics

  @REQ-FILEMGMT-352 @happy
  Scenario: Convert FileEntry with multiple paths to symlinks
    Given a package with a file entry with 3 paths
    When I convert the duplicate paths to symlinks
    Then the file entry should have 1 path
    And 2 symlink entries should be created
    And the symlinks should point to the primary path
    And the file entry metadata version should be incremented

  @REQ-FILEMGMT-354 @happy
  Scenario: Convert with explicit primary path
    Given a package with a file entry with paths "/data/file.txt" and "/backup/file.txt"
    When I convert with PrimaryPath set to "/data/file.txt"
    Then the primary path should be "/data/file.txt"
    And the symlink should point from "/backup/file.txt" to "/data/file.txt"

  @REQ-FILEMGMT-355 @happy
  Scenario: Convert with custom primary path selector
    Given a package with a file entry with paths "/a/very/long/path.txt" and "/short.txt"
    When I convert with a shortest-path selector function
    Then the primary path should be "/short.txt"
    And the symlink should point from "/a/very/long/path.txt" to "/short.txt"

  @REQ-FILEMGMT-363 @happy
  Scenario: Convert with new PrimaryPath not in FileEntry
    Given a package with a file entry with paths "/data/temp.txt" and "/backup/old.txt"
    When I convert with PrimaryPath set to "/data/canonical.txt"
    Then the file entry should have path "/data/canonical.txt" as primary
    And symlinks should be created from "/data/temp.txt" to "/data/canonical.txt"
    And symlinks should be created from "/backup/old.txt" to "/data/canonical.txt"
    And the file entry PathCount should be 1

  @REQ-FILEMGMT-370 @REQ-FILEMGMT-371 @REQ-FILEMGMT-373 @error
  Scenario: Reject new PrimaryPath with invalid format
    Given a package with a file entry with paths "/data/file1.txt"
    When I convert with PrimaryPath set to "../outside/package.txt"
    Then the conversion should fail with ErrTypeSecurity
    And the error message should indicate the path escapes package root

  @REQ-FILEMGMT-363 @error
  Scenario: Reject new PrimaryPath that conflicts with existing file
    Given a package with a file entry with path "/data/file1.txt"
    And the package already has a file at "/data/existing.txt"
    When I convert with PrimaryPath set to "/data/existing.txt"
    Then the conversion should fail with ErrTypeConflict
    And the error message should indicate the path already exists

  @REQ-FILEMGMT-364 @happy
  Scenario: Preserve path metadata during conversion
    Given a package with a file entry with 2 paths
    And each path has distinct permissions and timestamps
    When I convert with PreservePathMetadata enabled
    Then the symlink entries should preserve the original permissions
    And the symlink entries should preserve the original timestamps

  @REQ-FILEMGMT-352 @error
  Scenario: Reject conversion on signed package
    Given a signed package with a file entry with multiple paths
    When I attempt to convert the duplicate paths to symlinks
    Then the conversion should fail with ErrTypePackageState
    And the error message should indicate the package is signed

  @REQ-FILEMGMT-370 @REQ-FILEMGMT-371 @REQ-FILEMGMT-373 @error
  Scenario: Reject symlink pointing outside package root
    Given a package with a file entry with 2 paths
    When I attempt to convert with a target path containing ".." that escapes package root
    Then the conversion should fail with ErrTypeSecurity
    And the error message should indicate the path escapes package root

  @REQ-FILEMGMT-372 @REQ-FILEMGMT-374 @error
  Scenario: Reject symlink with non-existent target
    Given a package with a file entry with 2 paths
    And the primary path selector chooses a path that does not exist
    When I attempt to convert the duplicate paths to symlinks
    Then the conversion should fail with ErrTypeNotFound
    And the error message should indicate the target does not exist

  @REQ-FILEMGMT-372 @happy
  Scenario: Validate symlink target exists as FileEntry
    Given a package with a file entry "/data/file.txt" with 2 paths
    When I convert with primary path "/data/file.txt"
    Then the system should verify the FileEntry exists
    And the symlinks should be created successfully

  @REQ-FILEMGMT-372 @happy
  Scenario: Validate symlink target exists as directory PathMetadataEntry
    Given a package with a PathMetadataEntry for directory "/data/dir"
    And a file entry with 2 paths, one pointing to a file within "/data/dir"
    When I convert with primary path pointing to the directory
    Then the system should verify the PathMetadataEntry directory exists
    And the symlinks should be created successfully

  @REQ-FILEMGMT-356 @happy
  Scenario: Batch convert all multi-path entries
    Given a package with 10 file entries, each with multiple paths
    When I call ConvertAllPathsToSymlinks
    Then all file entries should have 1 path each
    And symlinks should be created for all duplicate paths
    And the progress callback should be invoked for each entry

  @REQ-FILEMGMT-357 @happy
  Scenario: Query multi-path entries
    Given a package with 5 file entries with multiple paths
    And 10 file entries with single paths
    When I call GetMultiPathEntries
    Then 5 file entries should be returned
    And each returned entry should have PathCount > 1

  @REQ-FILEMGMT-358 @happy
  Scenario: Get multi-path count
    Given a package with 5 file entries with multiple paths
    And 10 file entries with single paths
    When I call GetMultiPathCount
    Then the count should be 5
