@domain:metadata @m2 @REQ-META-067 @spec(api_metadata.md#55-special-file-management)
Feature: Special File Management

  @REQ-META-067 @happy
  Scenario: GetSpecialFiles returns all special files
    Given an open NovusPack package
    And special files in package
    When GetSpecialFiles is called
    Then array of SpecialFileInfo is returned
    And all special files are included
    And file information is complete

  @REQ-META-067 @happy
  Scenario: GetSpecialFileByType retrieves special file by type
    Given an open NovusPack package
    And special file with known type
    When GetSpecialFileByType is called with file type
    Then SpecialFileInfo is returned
    And file matches requested type
    And file information is accessible

  @REQ-META-067 @happy
  Scenario: RemoveSpecialFile removes special file by type
    Given an open writable NovusPack package
    And special file with known type
    When RemoveSpecialFile is called with file type
    Then special file is removed
    And GetSpecialFileByType returns error for removed file
    And package header flags are updated

  @REQ-META-067 @happy
  Scenario: ValidateSpecialFiles validates all special files
    Given an open NovusPack package
    And special files in package
    When ValidateSpecialFiles is called
    Then all special files are validated
    And validation result is returned
    And invalid files are identified

  @REQ-META-067 @happy
  Scenario: SpecialFileInfo structure contains file information
    Given an open NovusPack package
    And a special file
    When SpecialFileInfo is examined
    Then Type field contains file type
    And Name field contains special file name
    And Size field contains file size in bytes
    And Offset field contains offset in package
    And Data field contains file content
    And Valid field indicates validity
    And Error field contains error message if invalid

  @REQ-META-067 @happy
  Scenario: SpecialFileInfo supports different special file types
    Given an open NovusPack package
    And special files of various types
    When SpecialFileInfo is examined for each file
    Then Type field identifies file type correctly
    And Name field matches type-specific naming convention
    And file information is type-appropriate

  @REQ-META-011 @error
  Scenario: GetSpecialFileByType fails if file type not found
    Given an open NovusPack package
    And no special file with requested type
    When GetSpecialFileByType is called with non-existent type
    Then structured validation error is returned
    And error indicates file type not found

  @REQ-META-011 @error
  Scenario: RemoveSpecialFile fails if file type not found
    Given an open writable NovusPack package
    And no special file with requested type
    When RemoveSpecialFile is called with non-existent type
    Then structured validation error is returned
    And error indicates file type not found

  @REQ-META-011 @error
  Scenario: ValidateSpecialFiles identifies invalid files
    Given an open NovusPack package
    And special files including invalid file
    When ValidateSpecialFiles is called
    Then validation error is returned for invalid file
    And Valid field is false for invalid file
    And Error field contains error message
