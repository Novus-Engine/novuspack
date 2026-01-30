@domain:file_mgmt @REQ-PIPELINE-015 @REQ-PIPELINE-016 @REQ-PIPELINE-017 @spec(api_file_mgmt_transform_pipelines.md#14-disk-space-management)
Feature: Disk Space Management for Pipelines

  @REQ-PIPELINE-015 @REQ-FILEMGMT-343 @happy
  Scenario: Pre-flight disk space check before pipeline
    Given a file requiring multi-stage transformation
    When pipeline is about to start
    Then system estimates required disk space
    And estimation based on file sizes and transformation types
    And available disk space is checked in temp directory
    And pipeline proceeds if sufficient space

  @REQ-PIPELINE-016 @REQ-FILEMGMT-343 @error
  Scenario: Pipeline fails if insufficient disk space
    Given a large file requiring pipeline
    And estimated space requirement is 50GB
    And available disk space is 30GB
    When pipeline is initiated
    Then disk space check fails
    And ErrTypeIO error is returned
    And error message indicates insufficient space
    And error message includes space required and available

  @REQ-PIPELINE-017 @happy
  Scenario: Space estimation for decrypt transformation
    Given a file requiring decryption
    When space is estimated
    Then decrypt estimates approximately same size as encrypted data
    And space calculation accounts for output temporary file

  @REQ-PIPELINE-017 @happy
  Scenario: Space estimation for decompress transformation
    Given a file requiring decompression
    When space is estimated
    Then decompress estimates 2-10x size depending on compression ratio
    And estimation uses conservative multiplier for safety

  @REQ-PIPELINE-017 @happy
  Scenario: Space estimation for compress transformation
    Given a file requiring compression
    When space is estimated
    Then compress estimates 0.1-0.9x size depending on content
    And estimation considers compression algorithm characteristics

  @REQ-PIPELINE-017 @happy
  Scenario: Pipeline requires space for source plus all intermediate stages
    Given a 3-stage pipeline (source, stage1, stage2, final)
    When space is estimated
    Then calculation includes source file space
    And calculation includes all intermediate stage outputs
    And total accounts for simultaneous existence during pipeline
