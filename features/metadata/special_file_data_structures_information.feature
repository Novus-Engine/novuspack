@domain:metadata @m2 @REQ-META-068 @spec(api_metadata.md#551-special-file-data-structures)
Feature: Special File Data Structures

  @REQ-META-068 @happy
  Scenario: Special file data structures define special file information structures
    Given a NovusPack package
    When special file data structures are examined
    Then SpecialFileInfo structure is defined
    And ManifestData structure is defined
    And IndexData structure is defined
    And SignatureData structure is defined

  @REQ-META-068 @happy
  Scenario: SpecialFileInfo structure provides special file metadata
    Given a NovusPack package
    And SpecialFileInfo structure
    When structure is examined
    Then Type field contains FileType
    And Name field contains special file name
    And Size field contains file size in bytes
    And Offset field contains offset in package
    And Data field contains file content
    And Valid field indicates whether file is valid
    And Error field contains error message if invalid

  @REQ-META-068 @happy
  Scenario: ManifestData structure defines package manifest
    Given a NovusPack package
    And ManifestData structure
    When structure is examined
    Then Version field contains manifest version
    And Package field contains PackageInfo
    And Dependencies field contains dependency array
    And Structure field contains file organization
    And Install field contains installation instructions

  @REQ-META-068 @happy
  Scenario: IndexData and SignatureData structures define indexing and signature data
    Given a NovusPack package
    And IndexData structure
    And SignatureData structure
    When structures are examined
    Then IndexData contains Version, Files, Navigation, Search fields
    And SignatureData contains Version, Signatures, TrustChain, Validation fields

  @REQ-META-068 @error
  Scenario: Special file data structures validate structure formats
    Given a NovusPack package
    When invalid structure formats are provided
    Then structure validation detects format violations
    And appropriate errors are returned
