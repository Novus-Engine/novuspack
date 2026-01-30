@domain:metadata @m2 @REQ-META-089 @REQ-META-133 @REQ-META-134 @REQ-META-155 @REQ-META-157 @REQ-META-158 @REQ-META-159 @spec(api_metadata.md#82-pathmetadata-management-methods) @spec(api_metadata.md#path-information-query-methods) @spec(api_metadata.md#8222-getpathinfo-returns) @spec(api_metadata.md#8231-listpaths-returns) @spec(api_metadata.md#8241-listdirectories-returns) @spec(api_metadata.md#8251-getdirectorycount-returns) @spec(api_metadata.md#8261-getpathhierarchy-returns)
Feature: Path Metadata Management Methods

  @REQ-META-089 @happy
  Scenario: Path metadata management methods provide path operations
    Given a NovusPack package
    And a valid context
    When path metadata management methods are used
    Then path metadata management methods are available
    And path information query methods are available
    And path validation methods are available
    And special metadata file management methods are available
    And path association management methods are available
    And context supports cancellation

  @REQ-META-089 @happy
  Scenario: Path metadata management methods provide CRUD operations
    Given a NovusPack package
    And a valid context
    When path metadata management is used
    Then GetPathMetadata retrieves path entries
    And SetPathMetadata sets path entries
    And AddPathMetadata adds a path entry
    And RemovePathMetadata removes a path entry
    And UpdatePathMetadata updates a path entry
    And context supports cancellation

  @REQ-META-089 @REQ-META-134 @happy
  Scenario: Path information query methods provide path information
    Given a NovusPack package
    And a valid context
    When path information queries are used
    Then GetPathInfo gets path information by path
    And ListPaths lists all paths
    And GetPathHierarchy gets path hierarchy mapping
    And GetPathHierarchy returns a map of parent paths to child path slices
    And context supports cancellation

  @REQ-META-157 @happy
  Scenario: GetPathInfo returns PathInfo and error
    Given a NovusPack package
    And a valid context
    When GetPathInfo is called
    Then a PathInfo is returned on success
    And error is nil on success

  @REQ-META-158 @happy
  Scenario: ListPaths returns slice of PathInfo and error
    Given a NovusPack package
    And a valid context
    When ListPaths is called
    Then a slice of PathInfo is returned on success
    And error is nil on success

  @REQ-META-159 @happy
  Scenario: ListDirectories returns slice of PathInfo and error
    Given a NovusPack package
    And a valid context
    When ListDirectories is called
    Then a slice of directory PathInfo is returned on success
    And error is nil on success

  @REQ-META-133 @happy
  Scenario: GetDirectoryCount returns total count of directories
    Given a NovusPack package
    And a valid context
    When GetDirectoryCount is called
    Then an integer count of directories is returned
    And error is nil on success

  @REQ-META-089 @happy
  Scenario: Path validation methods validate path metadata
    Given a NovusPack package
    And a valid context
    When path validation is performed
    Then ValidatePathMetadata validates all path metadata
    And GetPathConflicts gets path conflicts
    And context supports cancellation

  @REQ-META-089 @error
  Scenario: Path metadata management methods handle errors
    Given a NovusPack package
    When invalid path operations are performed
    Then appropriate errors are returned
    And errors follow structured error format
