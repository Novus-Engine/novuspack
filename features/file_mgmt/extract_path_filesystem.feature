@domain:file_mgmt @extraction @filesystem @REQ-FILEMGMT-397 @REQ-FILEMGMT-398 @REQ-FILEMGMT-399 @REQ-FILEMGMT-400 @spec(api_file_mgmt_extraction.md#1-extractpath-package-method)
Feature: ExtractPath filesystem extraction

  @REQ-FILEMGMT-397 @REQ-FILEMGMT-398 @REQ-FILEMGMT-399 @happy
  Scenario: ExtractPath extracts single file to filesystem
    Given an open NovusPack package
    And file "documents/data.txt" exists in the package
    And session base is set to "/tmp/extract"
    When ExtractPath is called with the file path
    Then file is extracted to filesystem
    And file content matches original
    And file is located at "/tmp/extract/documents/data.txt"

  @REQ-FILEMGMT-397 @REQ-FILEMGMT-398 @REQ-FILEMGMT-399 @happy
  Scenario: ExtractPath extracts directory subtree to filesystem
    Given an open NovusPack package
    And directory "documents/" exists with files
    And session base is set to "/tmp/extract"
    When ExtractPath is called with the directory path
    Then all files in directory are extracted
    And directory structure is preserved
    And files are located under "/tmp/extract/documents/"

  @REQ-FILEMGMT-402 @error
  Scenario: ExtractPath requires session base for default-relative extraction
    Given an open NovusPack package
    And file "data.txt" exists in the package
    And session base is not set
    And no call-time destination override is provided
    When ExtractPath is called with the file path
    Then ErrTypeValidation error is returned
    And error indicates session base is required

  @REQ-FILEMGMT-402 @happy
  Scenario: ExtractPath uses session base from package runtime property
    Given an open NovusPack package
    And session base is set to "/tmp/extract"
    And file "data.txt" exists in the package
    When ExtractPath is called with the file path
    Then file is extracted to "/tmp/extract/data.txt"

  @REQ-FILEMGMT-402 @happy
  Scenario: ExtractPath uses session base from call-time options
    Given an open NovusPack package
    And file "data.txt" exists in the package
    When ExtractPath is called with SessionBase option "/tmp/extract"
    Then file is extracted to "/tmp/extract/data.txt"
    And package runtime session base is updated to "/tmp/extract"

  @REQ-FILEMGMT-400 @happy
  Scenario: ExtractPath creates parent directories as needed
    Given an open NovusPack package
    And file "deep/nested/path/file.txt" exists
    And session base is set to "/tmp/extract"
    When ExtractPath is called with the file path
    Then parent directories are created
    And file is extracted to "/tmp/extract/deep/nested/path/file.txt"

  @REQ-FILEMGMT-400 @happy
  Scenario: ExtractPath applies filesystem metadata from PathMetadataEntry
    Given an open NovusPack package
    And file "data.txt" exists with PathMetadataEntry
    And PathMetadataEntry has permissions and timestamps
    And session base is set to "/tmp/extract"
    When ExtractPath is called with the file path
    Then file is extracted with metadata applied
    And permissions match PathMetadataEntry
    And timestamps match PathMetadataEntry (best effort)

  @REQ-FILEMGMT-404 @error
  Scenario: ExtractPath returns error when file not found
    Given an open NovusPack package
    And file does not exist at specified path
    When ExtractPath is called with non-existent path
    Then ErrFileNotFound error is returned
    And error indicates file not found

  @REQ-FILEMGMT-404 @error
  Scenario: ExtractPath returns error when package not open
    Given a package that is not open
    When ExtractPath is called
    Then ErrPackageNotOpen error is returned

  @REQ-FILEMGMT-398 @happy
  Scenario: ExtractPath handles Windows path semantics when isWindows is true
    Given an open NovusPack package
    And file "data.txt" exists in the package
    And session base is set to "C:\\extract"
    When ExtractPath is called with isWindows=true
    Then file is extracted using Windows path semantics
    And file is located at "C:\\extract\\data.txt"

  @REQ-FILEMGMT-398 @happy
  Scenario: ExtractPath handles Unix path semantics when isWindows is false
    Given an open NovusPack package
    And file "data.txt" exists in the package
    And session base is set to "/tmp/extract"
    When ExtractPath is called with isWindows=false
    Then file is extracted using Unix path semantics
    And file is located at "/tmp/extract/data.txt"
