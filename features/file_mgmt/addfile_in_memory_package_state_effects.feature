@domain:file_mgmt @REQ-FILEMGMT-436 @spec(api_file_mgmt_addition.md#in-memory-package-state-effects) @spec(api_file_mgmt_addition.md#21-packageaddfile-method)
Feature: AddFile In-Memory Package State Effects

  On success, AddFile updates the package-level list of FileEntry objects
  and PackageInfo. The created or updated FileEntry is visible to
  subsequent in-process operations without a disk write.

  @REQ-FILEMGMT-436 @happy
  Scenario: AddFile updates package list with new or updated FileEntry
    Given an open writable package
    And a filesystem file
    When AddFile completes successfully
    Then package-level list of FileEntry objects is updated
    And created or updated FileEntry is in the list
    And no disk write is required for the update

  @REQ-FILEMGMT-436 @happy
  Scenario: New FileEntry is visible to subsequent operations
    Given an open writable package
    And a filesystem file
    When AddFile completes successfully
    Then ListFiles includes the new file
    And GetFileByPath returns the new FileEntry
    And Find operations can locate the entry
    And visibility is in-process without disk write

  @REQ-FILEMGMT-436 @happy
  Scenario: AddFile updates PackageInfo for new state
    Given an open writable package
    And a filesystem file
    When AddFile completes successfully
    Then PackageInfo reflects new in-memory package state
    And PackageInfo.PackageDataVersion is incremented
    And state is consistent for subsequent operations
