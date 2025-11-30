@domain:metadata @m2 @REQ-META-025 @spec(metadata.md#128-novuspack-special-files)
Feature: NovusPack Metadata Implementation

  @REQ-META-025 @happy
  Scenario: NovusPack special files provides metadata file reference support
    Given a NovusPack package
    When NovusPack special files tag value type is used
    Then NovusPackMetadata type (0x10) supports special metadata file references
    And value is stored as UTF-8 string
    And reference points to NovusPack special metadata files

  @REQ-META-025 @happy
  Scenario: NovusPackMetadata type stores metadata file references
    Given a NovusPack package
    And a tag with NovusPackMetadata value type
    When NovusPackMetadata tag is set
    Then value is stored as UTF-8 string
    And reference identifies special metadata file
    And reference can point to metadata, manifest, index, or signature files

  @REQ-META-025 @happy
  Scenario: NovusPackMetadata supports special file references
    Given a NovusPack package
    And special metadata files in package
    When NovusPackMetadata tags are used
    Then references can point to package metadata files
    And references can point to package manifest files
    And references can point to package index files
    And references can point to signature files

  @REQ-META-025 @error
  Scenario: NovusPackMetadata validates reference format
    Given a NovusPack package
    When invalid metadata file reference is provided
    Then reference validation detects invalid references
    And appropriate error is returned
