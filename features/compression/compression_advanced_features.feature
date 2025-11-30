@domain:compression @m2 @REQ-COMPR-049 @spec(api_package_compression.md#11353-advanced-features)
Feature: Compression Advanced Features

  @REQ-COMPR-049 @happy
  Scenario: UseSolidCompression improves compression ratios
    Given compression operations with multiple files
    When UseSolidCompression is enabled
    Then multiple files are treated as single stream
    And compression ratio improves
    And solid compression enhances efficiency

  @REQ-COMPR-049 @happy
  Scenario: ResumeFromOffset supports resumable compression
    Given a compression operation that needs to resume
    When ResumeFromOffset is set to specific offset
    Then compression resumes from that offset
    And previously processed data is skipped
    And compression continues from specified point

  @REQ-COMPR-049 @happy
  Scenario: BufferPoolSize controls buffer pool allocation
    Given compression operations requiring buffer management
    When BufferPoolSize is set
    Then buffer pool size matches configuration
    And buffer allocation is controlled
    And memory usage is predictable

  @REQ-COMPR-049 @happy
  Scenario: ProgressCallback provides real-time progress updates
    Given compression operations with progress tracking
    When ProgressCallback is configured
    Then real-time progress updates are received
    And bytes processed are reported
    And total bytes are reported
    And progress tracking enhances user experience

  @REQ-COMPR-049 @happy
  Scenario: Advanced features enhance streaming execution
    Given advanced streaming compression execution
    When advanced features are configured
    Then solid compression improves ratios
    And resumable operations support recovery
    And progress tracking provides feedback
    And advanced features optimize execution
