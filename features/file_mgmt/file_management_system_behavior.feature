@skip @domain:file_mgmt @m2 @spec(api_file_mgmt_addition.md#215-addfile-behavior) @spec(api_file_mgmt_removal.md#1-file-removal-semantics-and-multi-path-files) @spec(api_file_mgmt_removal.md#21-removefile-purpose)
Feature: File Management System Behavior

# This feature captures high-level file management behaviors from the file management index spec.
# Detailed runnable scenarios live in the dedicated file_mgmt feature files.

  @REQ-FILEMGMT-072 @behavior
  Scenario: AddFile updates in-memory package state without writing to disk
    Given a package configured for writing
    When the caller adds a file via AddFile
    Then a new FileEntry is created in the in-memory package state
    And file metadata is captured according to the configured options
    And the package is not written to disk until a write operation is invoked

  @REQ-FILEMGMT-136 @behavior
  Scenario: RemoveFile updates in-memory state and associated metadata
    Given a package that contains a file entry at a path
    When the caller removes that file via RemoveFile
    Then the file entry is removed from the in-memory package state
    And path metadata is cleaned up according to the removal rules
