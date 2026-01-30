@domain:file_mgmt @extraction @platform @REQ-FILEMGMT-375 @REQ-FILEMGMT-376 @REQ-FILEMGMT-377 @REQ-FILEMGMT-378 @spec(api_file_mgmt_extraction.md#522-extractoptions-symlink-behavior)
Feature: Extract symlinks with platform-specific behavior

  As a package user
  I want symlinks to be extracted appropriately for my platform
  So that files are accessible and security is maintained

  @REQ-FILEMGMT-375 @happy
  Scenario: Default extraction behavior on Unix-like systems
    Given a package with symlinks
    And the platform is Unix-like
    When I extract with default options
    Then symlinks should be extracted as symlinks
    And no errors should occur

  @REQ-FILEMGMT-376 @happy
  Scenario: Default extraction behavior on Windows
    Given a package with symlinks
    And the platform is Windows
    When I extract with default options
    Then symlinks should be extracted as regular file copies
    And no privilege errors should occur

  @REQ-FILEMGMT-375 @happy
  Scenario: Unix user explicitly requests symlinks
    Given a package with symlinks
    And the platform is Unix-like
    When I extract with PreserveSymlinks set to Some(true)
    Then symlinks should be extracted as symlinks
    And no errors should occur

  @REQ-FILEMGMT-376 @happy
  Scenario: Windows user explicitly requests file copies
    Given a package with symlinks
    And the platform is Windows
    When I extract with PreserveSymlinks set to Some(false)
    Then symlinks should be extracted as regular file copies
    And no errors should occur

  @REQ-FILEMGMT-377 @REQ-FILEMGMT-378 @error
  Scenario: Windows user explicitly requests symlinks without privileges
    Given a package with symlinks
    And the platform is Windows
    And the user has insufficient privileges
    When I extract with PreserveSymlinks set to Some(true)
    Then the extraction should fail with ErrTypeIO
    And the error message should indicate privilege requirements
    And symlinks should not be silently converted to file copies

  @REQ-FILEMGMT-375 @happy
  Scenario: Unix user explicitly requests file copies
    Given a package with symlinks
    And the platform is Unix-like
    When I extract with PreserveSymlinks set to Some(false)
    Then symlinks should be extracted as regular file copies
    And no errors should occur
