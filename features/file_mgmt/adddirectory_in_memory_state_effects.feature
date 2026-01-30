@domain:file_mgmt @REQ-FILEMGMT-445 @REQ-FILEMGMT-313 @spec(api_file_mgmt_addition.md#in-memory-package-state-effects-adddirectory) @spec(api_file_mgmt_addition.md#25-adddirectory-package-method)
Feature: AddDirectory In-Memory Package State Effects

  AddDirectory updates the package-level list of FileEntry objects for each
  successfully added file. Each added FileEntry is visible to subsequent
  in-process operations without disk write. PackageInfo and PackageDataVersion
  are updated per file.

  @REQ-FILEMGMT-445 @happy
  Scenario: AddDirectory updates package list per successfully added file
    Given an open writable package
    And a filesystem directory with N files
    When AddDirectory completes successfully
    Then package-level list of FileEntry objects is updated N times
    And each successfully added FileEntry is in the list
    And no disk write is required for visibility

  @REQ-FILEMGMT-445 @happy
  Scenario: Added FileEntries visible to subsequent operations
    Given an open writable package
    And AddDirectory has added files from a directory
    When ListFiles or GetFileByPath or Find is called
    Then all added FileEntries are visible
    And visibility is in-process without disk write

  @REQ-FILEMGMT-445 @happy
  Scenario: PackageInfo and PackageDataVersion updated per file
    Given an open writable package
    And initial PackageInfo.PackageDataVersion value
    When AddDirectory successfully adds M files
    Then PackageInfo reflects new in-memory package state
    And PackageDataVersion is incremented for each successfully added or updated file
