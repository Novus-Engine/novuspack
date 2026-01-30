@domain:basic_ops @m2 @REQ-API_BASIC-128 @spec(api_basic_operations.md#321-package-loadspecialmetadatafiles-method)
Feature: Package loadSpecialMetadataFiles method

  @REQ-API_BASIC-128 @happy
  Scenario: loadSpecialMetadataFiles loads all special metadata files
    Given a package opened from disk
    And the package contains special metadata files (file types 65000-65535)
    When loadSpecialMetadataFiles is invoked during package load
    Then all special metadata files are loaded into memory
    And special file types are mapped to their corresponding FileEntry records
    And missing optional special files are handled according to the spec
    And loaded special metadata is available to downstream metadata operations

