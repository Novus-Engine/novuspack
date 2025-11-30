@domain:metadata @m2 @REQ-META-042 @spec(metadata.md#23-metadata-file-api)
Feature: Metadata File API

  @REQ-META-042 @happy
  Scenario: Metadata file API provides package metadata operations
    Given a NovusPack package
    When metadata file API is used
    Then package metadata operations are available
    And special metadata file management is available
    And API definitions are in Package Metadata API specification

  @REQ-META-042 @happy
  Scenario: Metadata file API operations reference authoritative API
    Given a NovusPack package
    When metadata file API is examined
    Then API definitions reference Package Metadata API specification
    And package metadata operations are defined in api_metadata.md
    And special metadata file management is defined in api_metadata.md

  @REQ-META-042 @happy
  Scenario: Metadata file API supports package metadata operations
    Given a NovusPack package
    When metadata file API is used
    Then AddMetadataFile operation is available
    And GetMetadataFile operation is available
    And UpdateMetadataFile operation is available
    And RemoveMetadataFile operation is available

  @REQ-META-042 @error
  Scenario: Metadata file API handles invalid operations
    Given a NovusPack package
    When invalid metadata file operations are performed
    Then appropriate errors are returned
    And errors follow structured error format
