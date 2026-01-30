@domain:file_mgmt @extraction @destination @REQ-FILEMGMT-401 @REQ-FILEMGMT-408 @spec(api_file_mgmt_extraction.md#151-extractpath-destination-resolution) @spec(api_file_mgmt_extraction.md#destpathspec-struct)
Feature: ExtractPath destination overrides

  @REQ-FILEMGMT-401 @REQ-FILEMGMT-408 @happy
  Scenario: Call-time file destination override takes precedence
    Given an open NovusPack package
    And file "data.txt" exists with stored DestPath "/custom/path/data.txt"
    And session base is set to "/tmp/extract"
    When ExtractPath is called with FileDestOverrides["data.txt"]="/override/path/data.txt"
    Then file is extracted to "/override/path/data.txt"
    And stored DestPath is ignored

  @REQ-FILEMGMT-401 @happy
  Scenario: Stored destination takes precedence over default
    Given an open NovusPack package
    And file "data.txt" exists with stored DestPath "/custom/path/data.txt"
    And session base is set to "/tmp/extract"
    When ExtractPath is called without call-time overrides
    Then file is extracted to "/custom/path/data.txt"
    And default path "/tmp/extract/data.txt" is not used

  @REQ-FILEMGMT-401 @happy
  Scenario: Default destination relative to session base when no overrides
    Given an open NovusPack package
    And file "documents/data.txt" exists without stored destination
    And session base is set to "/tmp/extract"
    When ExtractPath is called without overrides
    Then file is extracted to "/tmp/extract/documents/data.txt"

  @REQ-FILEMGMT-408 @happy
  Scenario: Directory destination override applies to all files in directory
    Given an open NovusPack package
    And directory "documents/" exists with files "a.txt" and "b.txt"
    And session base is set to "/tmp/extract"
    When ExtractPath is called with DirDestOverrides["documents/"]="/custom/docs"
    Then "documents/a.txt" is extracted to "/custom/docs/a.txt"
    And "documents/b.txt" is extracted to "/custom/docs/b.txt"

  @REQ-FILEMGMT-408 @happy
  Scenario: Root destination override applies to all paths
    Given an open NovusPack package
    And files "a.txt" and "b/c.txt" exist
    And session base is set to "/tmp/extract"
    When ExtractPath is called with RootDestOverride="/custom/root"
    Then "a.txt" is extracted to "/custom/root/a.txt"
    And "b/c.txt" is extracted to "/custom/root/b/c.txt"

  @REQ-FILEMGMT-409 @happy
  Scenario: IgnoreStoredDestPaths ignores stored destinations
    Given an open NovusPack package
    And file "data.txt" exists with stored DestPath "/custom/path/data.txt"
    And session base is set to "/tmp/extract"
    When ExtractPath is called with IgnoreStoredDestPaths=true
    Then file is extracted to "/tmp/extract/data.txt"
    And stored DestPath is ignored

  @REQ-FILEMGMT-401 @happy
  Scenario: Relative stored destination resolved from default extraction directory
    Given an open NovusPack package
    And file "documents/data.txt" exists with stored DestPath="../custom/data.txt"
    And session base is set to "/tmp/extract"
    When ExtractPath is called
    Then file is extracted to "/tmp/custom/data.txt"
    And relative path is resolved from "/tmp/extract/documents/"

  @REQ-FILEMGMT-401 @happy
  Scenario: Absolute stored destination used as-is
    Given an open NovusPack package
    And file "data.txt" exists with stored DestPath="/absolute/path/data.txt"
    And session base is set to "/tmp/extract"
    When ExtractPath is called
    Then file is extracted to "/absolute/path/data.txt"
    And session base is not used

  @REQ-FILEMGMT-401 @happy
  Scenario: File destination override specifies full file path
    Given an open NovusPack package
    And file "data.txt" exists in package
    And session base is set to "/tmp/extract"
    When ExtractPath is called with FileDestOverrides["data.txt"]="/custom/file.txt"
    Then file is extracted to "/custom/file.txt"
    And file name is "file.txt" not "data.txt"

  @REQ-FILEMGMT-401 @happy
  Scenario: Directory destination override appends file path suffix
    Given an open NovusPack package
    And file "documents/data.txt" exists
    And session base is set to "/tmp/extract"
    When ExtractPath is called with DirDestOverrides["documents/"]="/custom/docs"
    Then file is extracted to "/custom/docs/data.txt"
    And file name "data.txt" is preserved
