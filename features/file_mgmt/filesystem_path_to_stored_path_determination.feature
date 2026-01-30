@domain:file_mgmt @m2 @REQ-FILEMGMT-305 @REQ-FILEMGMT-306 @REQ-FILEMGMT-307 @REQ-FILEMGMT-308 @spec(api_file_mgmt_addition.md#213-filesystem-input-path-and-stored-path-derivation)
Feature: Filesystem path to stored path determination

  Background:
    Given an open writable NovusPack package
    And a valid context

  @REQ-FILEMGMT-305 @happy
  Scenario: First absolute AddFile establishes session base and preserves one parent directory
    When AddFile is called with filesystem path "/home/user/project/assets/texture.png" and default options
    Then the stored path is "/assets/texture.png"
    When AddFile is called with filesystem path "/home/user/project/src/main.go" and default options
    Then the stored path is "/src/main.go"

  @REQ-FILEMGMT-305 @error
  Scenario: AddFile fails when a subsequent absolute filesystem path is outside the established base
    Given AddFile was previously called with filesystem path "/home/user/project/assets/texture.png" and default options
    When AddFile is called with filesystem path "/home/other/file.txt" and default options
    Then a structured validation error is returned
    And the error indicates the file is not under the established base path
    And the error suggests setting BasePath or StoredPath

  @REQ-FILEMGMT-308 @happy
  Scenario: StoredPath overrides filesystem path mapping
    When AddFile is called with filesystem path "/any/filesystem/path/file.txt" and StoredPath "/custom/location/file.txt"
    Then the stored path is "/custom/location/file.txt"
    And the filesystem path is used only as a filesystem input

  @REQ-FILEMGMT-307 @happy
  Scenario: AddFilePattern derives stored paths from pattern base
    Given the filesystem contains "/home/user/project/assets/button.png"
    And the filesystem contains "/home/user/project/ui/icons/save.png"
    When AddFilePattern is called with pattern "/home/user/project/**/*.png"
    Then stored paths include "/project/assets/button.png"
    And stored paths include "/project/ui/icons/save.png"
    And directory structure is preserved for pattern operations

  @REQ-FILEMGMT-308 @error
  Scenario: Mutually exclusive path determination options are rejected
    When AddFile is called with filesystem path "/home/user/project/assets/texture.png" and both StoredPath and BasePath set
    Then a structured validation error is returned
    And the error indicates mutually exclusive path determination options
