@domain:streaming @m2 @REQ-STREAM-022 @spec(api_streaming.md#13-key-methods)
Feature: Streaming: Streaming Key Methods (Core Operations)

  @REQ-STREAM-022 @happy
  Scenario: Key methods provide core streaming operations
    Given a NovusPack package
    And an open file for streaming
    When key streaming methods are used
    Then NewFileStream creates file stream for large files
    And ReadChunk performs sequential chunk reads
    And Seek changes stream position
    And Close closes stream and releases resources
    And GetStats provides stream statistics

  @REQ-STREAM-022 @happy
  Scenario: NewFileStream creates file stream with configuration
    Given a NovusPack package
    And an io.Reader
    And a StreamConfig
    When NewFileStream is called
    Then FileStream is created with reader
    And stream is configured with StreamConfig
    And stream supports compression if configured
    And stream supports encryption if configured
    And stream integrates with buffer pool if configured

  @REQ-STREAM-022 @happy
  Scenario: ReadChunk performs sequential chunk reads
    Given a NovusPack package
    And an open FileStream
    And a valid context
    When ReadChunk is called
    Then chunk of data is returned
    And chunk size matches configured chunk size
    And stream position advances
    And context cancellation is respected

  @REQ-STREAM-022 @error
  Scenario: Key methods handle errors correctly
    Given a NovusPack package
    When streaming operation encounters error
    Then structured error is returned
    And error indicates specific failure
    And error follows structured error format
