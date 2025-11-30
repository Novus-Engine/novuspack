@domain:metadata @m2 @REQ-META-089 @spec(api_metadata.md#82-directory-management-methods)
Feature: Directory Management Methods

  @REQ-META-089 @happy
  Scenario: Directory management methods provide directory operations
    Given a NovusPack package
    And a valid context
    When directory management methods are used
    Then directory metadata management methods are available
    And directory information query methods are available
    And directory validation methods are available
    And special metadata file management methods are available
    And directory association management methods are available
    And context supports cancellation

  @REQ-META-089 @happy
  Scenario: Directory metadata management methods provide CRUD operations
    Given a NovusPack package
    And a valid context
    When directory metadata management is used
    Then GetDirectoryMetadata retrieves directory entries
    And SetDirectoryMetadata sets directory entries
    And AddDirectory adds a directory entry
    And RemoveDirectory removes a directory entry
    And UpdateDirectory updates a directory entry
    And context supports cancellation

  @REQ-META-089 @happy
  Scenario: Directory information query methods provide directory information
    Given a NovusPack package
    And a valid context
    When directory information queries are used
    Then GetDirectoryInfo gets directory information by path
    And ListDirectories lists all directories
    And GetDirectoryHierarchy gets directory hierarchy mapping
    And context supports cancellation

  @REQ-META-089 @happy
  Scenario: Directory validation methods validate directory metadata
    Given a NovusPack package
    And a valid context
    When directory validation is performed
    Then ValidateDirectoryMetadata validates all directory metadata
    And GetDirectoryConflicts gets directory conflicts
    And context supports cancellation

  @REQ-META-089 @error
  Scenario: Directory management methods handle errors
    Given a NovusPack package
    When invalid directory operations are performed
    Then appropriate errors are returned
    And errors follow structured error format
