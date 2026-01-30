@domain:metadata @m2 @REQ-META-063 @spec(api_metadata.md#51-package-metadata-file-type-65000) @spec(api_metadata.md#515-packagehasmetadatafile-method)
Feature: Package Metadata File Type 65000

  @REQ-META-063 @happy
  Scenario: Package metadata file type 65000 provides metadata file operations
    Given a NovusPack package
    When package metadata file operations are used
    Then AddMetadataFile adds YAML metadata file
    And GetMetadataFile retrieves metadata
    And UpdateMetadataFile updates metadata
    And RemoveMetadataFile removes metadata file
    And HasMetadataFile checks for metadata file

  @REQ-META-063 @happy
  Scenario: AddMetadataFile adds YAML metadata file
    Given a NovusPack package
    And metadata map[string]interface{}
    When AddMetadataFile is called
    Then YAML metadata file is added to package
    And file type is set to 65000
    And file contains structured YAML metadata

  @REQ-META-063 @happy
  Scenario: Package metadata file contains structured information
    Given a NovusPack package
    And a package metadata file
    When metadata file is examined
    Then file contains package description and version information
    And file contains author and license details
    And file contains build and compilation metadata
    And file contains custom package-specific data

  @REQ-META-063 @happy
  Scenario: GetMetadataFile retrieves metadata
    Given a NovusPack package
    And a package with metadata file
    When GetMetadataFile is called
    Then metadata map[string]interface{} is returned
    And metadata contains all package information

  @REQ-META-063 @error
  Scenario: Package metadata file operations handle errors
    Given a NovusPack package
    When invalid metadata or file operations fail
    Then appropriate errors are returned
    And errors follow structured error format
