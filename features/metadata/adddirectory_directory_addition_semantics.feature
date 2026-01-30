@domain:file_mgmt @m2 @REQ-FILEMGMT-313 @REQ-FILEMGMT-314 @REQ-FILEMGMT-315 @spec(api_file_mgmt_addition.md#24-adddirectory)
Feature: Directory addition semantics

  Background:
    Given an open writable NovusPack package
    And a valid context

  @REQ-FILEMGMT-313 @REQ-FILEMGMT-314 @happy
  Scenario: AddDirectory recursively adds files from a filesystem directory
    Given the filesystem directory "/home/user/project/assets" exists
    And the filesystem contains "/home/user/project/assets/a.png"
    And the filesystem contains "/home/user/project/assets/sub/b.png"
    When AddDirectory is called with directory path "/home/user/project/assets" and BasePath "/home/user/project"
    Then returned entries include a FileEntry for "/assets/a.png"
    And returned entries include a FileEntry for "/assets/sub/b.png"
    And directory metadata entries exist for "/assets/" and "/assets/sub/"
    And AddFile was called internally for each file
