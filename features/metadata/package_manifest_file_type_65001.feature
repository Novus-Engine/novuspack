@domain:metadata @m2 @REQ-META-064 @spec(api_metadata.md#52-package-manifest-file-type-65001) @spec(api_metadata.md#525-packagehasmanifestfile-method)
Feature: Package Manifest File Type 65001

  @REQ-META-064 @happy
  Scenario: Package manifest file type 65001 provides manifest file operations
    Given a NovusPack package
    When package manifest file operations are used
    Then AddManifestFile adds manifest file
    And GetManifestFile retrieves manifest
    And UpdateManifestFile updates manifest
    And RemoveManifestFile removes manifest file
    And HasManifestFile checks for manifest file

  @REQ-META-064 @happy
  Scenario: AddManifestFile adds package manifest file
    Given a NovusPack package
    And ManifestData
    When AddManifestFile is called
    Then manifest file is added to package
    And file type is set to 65001
    And file contains ManifestData

  @REQ-META-064 @happy
  Scenario: Package manifest file defines package structure
    Given a NovusPack package
    And a package manifest file
    When manifest file is examined
    Then file defines file organization and structure
    And file defines dependency requirements
    And file defines installation instructions
    And file defines package relationships

  @REQ-META-064 @happy
  Scenario: GetManifestFile retrieves manifest
    Given a NovusPack package
    And a package with manifest file
    When GetManifestFile is called
    Then ManifestData is returned
    And manifest contains all package structure information

  @REQ-META-064 @error
  Scenario: Package manifest file operations handle errors
    Given a NovusPack package
    When invalid manifest or file operations fail
    Then appropriate errors are returned
    And errors follow structured error format
