@domain:file_mgmt @extraction @security @REQ-FILEMGMT-379 @REQ-FILEMGMT-380 @REQ-FILEMGMT-381 @REQ-FILEMGMT-382 @REQ-FILEMGMT-383 @REQ-FILEMGMT-384 @REQ-FILEMGMT-385 @REQ-FILEMGMT-386 @REQ-FILEMGMT-387 @REQ-FILEMGMT-388 @spec(api_file_mgmt_extraction.md#523-extractoptions-security-limits)
Feature: Extraction security limits

  As a package system
  I want to enforce security limits during extraction
  So that malicious packages cannot cause denial of service

  @REQ-FILEMGMT-379 @error
  Scenario: Reject symlink chain exceeding depth limit
    Given a package with a symlink chain of depth 50
    And MaxSymlinkDepth is set to 40
    When I attempt to extract the symlink chain
    Then extraction should fail with ErrTypeSecurity
    And the error message should indicate symlink depth limit exceeded

  @REQ-FILEMGMT-380 @error
  Scenario: Reject circular symlinks
    Given a package with circular symlinks
    And RejectCircularSymlinks is enabled
    When I attempt to extract the circular symlinks
    Then extraction should fail with ErrTypeSecurity
    And the error message should indicate circular symlink detected

  @REQ-FILEMGMT-381 @error
  Scenario: Reject file with excessive compression ratio
    Given a package with a file compressed at 2000:1 ratio
    And MaxCompressionRatio is set to 1000
    When I attempt to extract the file
    Then extraction should fail with ErrTypeSecurity
    And the error message should indicate compression ratio limit exceeded

  @REQ-FILEMGMT-384 @error
  Scenario: Reject extraction exceeding file count limit
    Given a package with 2,000,000 files
    And MaxFileCount is set to 1,000,000
    When I attempt to extract all files
    Then extraction should fail with ErrTypeSecurity
    And the error message should indicate file count limit exceeded

  @REQ-FILEMGMT-385 @error
  Scenario: Reject extraction exceeding directory depth limit
    Given a package with directory depth of 300
    And MaxDirectoryDepth is set to 250
    When I attempt to extract the deep directory structure
    Then extraction should fail with ErrTypeSecurity
    And the error message should indicate directory depth limit exceeded

  @REQ-FILEMGMT-386 @error
  Scenario: Reject extraction when insufficient disk space
    Given a package requiring 100GB of space
    And available disk space is 50GB
    When I attempt to extract the package
    Then extraction should fail before beginning
    And the error message should indicate insufficient disk space

  @REQ-FILEMGMT-387 @error
  Scenario: Stop extraction when disk space exhausted during extraction
    Given a package being extracted
    And disk space becomes exhausted during extraction
    When the space check detects insufficient space
    Then extraction should stop immediately
    And the error message should indicate disk space exhausted
    And partial files should be cleaned up

  @REQ-FILEMGMT-388 @happy
  Scenario: Verify sufficient space for each file before extracting
    Given a package with multiple files
    And sufficient disk space for all files
    When I extract the package
    Then space should be verified for each file before extraction
    And all files should be extracted successfully

  @REQ-FILEMGMT-382 @happy
  Scenario: Optional total extracted size limit disabled by default
    Given a package
    And ExtractOptions with default values
    When I extract the package
    Then MaxTotalExtractedSize should be 0 (disabled)
    And disk space checks should be used instead

  @REQ-FILEMGMT-383 @happy
  Scenario: Optional individual file size limit disabled by default
    Given a package
    And ExtractOptions with default values
    When I extract the package
    Then MaxFileSize should be 0 (disabled)
    And disk space checks should be used instead
