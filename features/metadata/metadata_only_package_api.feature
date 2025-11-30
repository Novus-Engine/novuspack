@domain:metadata @m2 @REQ-META-081 @spec(api_metadata.md#64-metadata-only-package-api)
Feature: Metadata-Only Package API

  @REQ-META-081 @happy
  Scenario: Metadata-only package API provides metadata-only operations
    Given a NovusPack package
    When metadata-only package API is used
    Then IsMetadataOnlyPackage checks if package contains only metadata files
    And ValidateMetadataOnlyPackage validates a metadata-only package
    And CreateMetadataOnlyPackage creates a new metadata-only package
    And AddMetadataOnlyFile adds special metadata file
    And GetMetadataOnlyFiles returns all metadata files
    And ValidateMetadataOnlyIntegrity validates package integrity

  @REQ-META-081 @happy
  Scenario: IsMetadataOnlyPackage checks package type
    Given a NovusPack package
    And a metadata-only package
    When IsMetadataOnlyPackage is called
    Then function returns true if package is metadata-only
    And function returns false if package contains regular files

  @REQ-META-081 @happy
  Scenario: CreateMetadataOnlyPackage creates new metadata-only package
    Given a NovusPack package
    When CreateMetadataOnlyPackage is called
    Then new metadata-only package is created
    And package has FileCount of 0
    And package has no regular content files
    And package is ready for special metadata files

  @REQ-META-081 @happy
  Scenario: AddMetadataOnlyFile adds special metadata file
    Given a NovusPack package
    And a metadata-only package
    And a file type
    And file data
    When AddMetadataOnlyFile is called
    Then special metadata file is added to package
    And file type must be special file type
    And file data is stored

  @REQ-META-081 @error
  Scenario: Metadata-only package API handles errors
    Given a NovusPack package
    When invalid metadata-only package operations are performed
    Then appropriate errors are returned
    And errors follow structured error format
