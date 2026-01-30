@domain:metadata @m2 @REQ-META-123 @spec(metadata.md#2-package-metadata-file-specification)
Feature: Package Metadata File Specification
    As a package developer
    I want a package-level metadata storage format
    So that package-wide metadata is stored separately from file-specific metadata

    Background:
        Given package metadata is stored in a special metadata file
        And the file has a specific binary format with YAML structure

    @REQ-META-123 @happy
    Scenario: Package metadata file stores package-level properties
        When package metadata is set
        Then properties MUST be stored in the package metadata file
        And properties include name, version, author, description, etc.

    @REQ-META-123 @happy
    Scenario: Package metadata file uses YAML format
        When package metadata is serialized
        Then the content MUST use YAML format
        And the YAML MUST be valid and parseable

    @REQ-META-123 @happy
    Scenario: Package metadata file is discoverable via metadata index
        When a package contains package metadata
        Then the package metadata file MUST be listed in the metadata index
        And the file type identifier MUST indicate package metadata

    @REQ-META-123 @happy
    Scenario: Package metadata file is optional
        When a package is created without metadata
        Then the package metadata file MAY be omitted
        And the absence MUST NOT prevent package operations

    @REQ-META-123 @happy
    Scenario: Package metadata supports custom fields
        When custom metadata fields are added
        Then the package metadata file MUST support arbitrary key-value pairs
        And custom fields MUST coexist with standard fields

    @REQ-META-123 @happy
    Scenario: Package metadata is validated on read
        When package metadata is read
        Then the YAML structure MUST be validated
        And invalid or corrupted metadata MUST be detected and reported
