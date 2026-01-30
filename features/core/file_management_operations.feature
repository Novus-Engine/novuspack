@domain:core @m2 @REQ-CORE-128 @spec(api_core.md#122-file-management-operations)
Feature: File management operations define add and remove file operations

  @REQ-CORE-128 @happy
  Scenario: PackageWriter exposes add and remove file operations
    Given a package opened for writing
    When file management operations are used
    Then add file and remove file operations are available
    And the operations follow the file management specification
    And in-memory and disk effects are as specified
