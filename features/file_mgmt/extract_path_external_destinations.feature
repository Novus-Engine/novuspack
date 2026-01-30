@domain:file_mgmt @extraction @security @REQ-FILEMGMT-403 @REQ-FILEMGMT-410 @REQ-FILEMGMT-411 @spec(api_file_mgmt_extraction.md#153-extractpath-external-destination-handling)
Feature: ExtractPath external destination security

  @REQ-FILEMGMT-403 @error
  Scenario: ExtractPath errors on external destination by default
    Given an open NovusPack package
    And file "data.txt" exists with stored DestPath="/external/path/data.txt"
    And session base is set to "/tmp/extract"
    And "/external/path" is outside session base and not a package subdirectory
    When ExtractPath is called without AllowExternalDestinations
    Then ErrTypeSecurity error is returned
    And error indicates external destination not allowed

  @REQ-FILEMGMT-410 @happy
  Scenario: AllowExternalDestinations permits extraction outside session base
    Given an open NovusPack package
    And file "data.txt" exists with stored DestPath="/external/path/data.txt"
    And session base is set to "/tmp/extract"
    And "/external/path" is outside session base
    When ExtractPath is called with AllowExternalDestinations=true
    Then file is extracted to "/external/path/data.txt"
    And extraction completes successfully

  @REQ-FILEMGMT-411 @happy
  Scenario: SkipDisallowedExternalDestinations skips external files
    Given an open NovusPack package
    And file "data.txt" exists with stored DestPath="/external/path/data.txt"
    And file "allowed.txt" exists with stored DestPath="/tmp/extract/allowed.txt"
    And session base is set to "/tmp/extract"
    When ExtractPath is called with SkipDisallowedExternalDestinations=true
    Then "allowed.txt" is extracted
    And "data.txt" is skipped
    And no error is returned

  @REQ-FILEMGMT-403 @happy
  Scenario: Extraction to session base is allowed
    Given an open NovusPack package
    And file "data.txt" exists
    And session base is set to "/tmp/extract"
    When ExtractPath is called
    Then file is extracted to "/tmp/extract/data.txt"
    And extraction completes successfully

  @REQ-FILEMGMT-403 @happy
  Scenario: Extraction to package subdirectory is allowed
    Given an open NovusPack package
    And file "data.txt" exists
    And session base is set to "/tmp/extract"
    And "/tmp/extract/subdir" is a package subdirectory
    When ExtractPath is called with destination "/tmp/extract/subdir/data.txt"
    Then file is extracted successfully
    And extraction completes without error

  @REQ-FILEMGMT-403 @error
  Scenario: Extraction to non-package subdirectory requires allow flag
    Given an open NovusPack package
    And file "data.txt" exists
    And session base is set to "/tmp/extract"
    And "/tmp/extract/other" is not a package subdirectory
    When ExtractPath is called with destination "/tmp/extract/other/data.txt"
    Then ErrTypeSecurity error is returned
    And error indicates external destination

  @REQ-FILEMGMT-410 @REQ-FILEMGMT-411 @happy
  Scenario: SkipDisallowedExternalDestinations takes precedence over AllowExternalDestinations
    Given an open NovusPack package
    And file "external.txt" exists with external destination
    And session base is set to "/tmp/extract"
    When ExtractPath is called with AllowExternalDestinations=true and SkipDisallowedExternalDestinations=true
    Then external files are skipped
    And no error is returned
