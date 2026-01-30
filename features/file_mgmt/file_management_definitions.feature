@skip @domain:file_mgmt @m2 @REQ-FILEMGMT-050 @spec(api_file_mgmt_addition.md#21-addfile-package-method)
Feature: File Management Definitions

# This feature captures high-level file management definitions and flow expectations from the file management specs.
# Detailed runnable scenarios live in the dedicated file_mgmt feature files.

  @REQ-FILEMGMT-050 @documentation
  Scenario: AddFile adds new content to the package via a FileSource and options
    Given a package configured for writing
    When the caller invokes AddFile with a FileSource and AddFileOptions
    Then a new FileEntry is created for the added content
    And the FileEntry includes captured metadata as configured by options

  @REQ-FILEMGMT-072 @documentation
  Scenario: AddFile follows the documented implementation flow and processing stages
    Given a file is being added to the package
    When AddFile executes the addition flow
    Then it validates inputs and selects processing options
    And it applies configured transformation and processing stages in the documented order
