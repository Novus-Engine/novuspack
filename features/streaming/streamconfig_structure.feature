@domain:streaming @m2 @REQ-STREAM-021 @spec(api_streaming.md#122-streamconfig-struct)
Feature: StreamConfig Structure

  @REQ-STREAM-021 @happy
  Scenario: StreamConfig struct provides stream configuration structure
    Given a NovusPack package
    When StreamConfig struct is used
    Then struct contains BufferSize field for read buffer size
    And struct contains ChunkSize field for read chunk size
    And struct contains MemoryLimit field for maximum memory
    And struct contains UseBufferPool field for buffer pool usage
    And struct contains IsCompressed field for compression status
    And struct contains IsEncrypted field for encryption status

  @REQ-STREAM-021 @happy
  Scenario: BufferSize field configures read buffer size
    Given a NovusPack package
    And a StreamConfig
    When BufferSize field is set
    Then size of read buffer is configured
    And zero value uses default buffer size
    And buffer size affects read performance

  @REQ-STREAM-021 @happy
  Scenario: ChunkSize field configures read chunk size
    Given a NovusPack package
    And a StreamConfig
    When ChunkSize field is set
    Then size of each read chunk is configured
    And zero value triggers calculated chunk size
    And chunk size affects streaming performance

  @REQ-STREAM-021 @happy
  Scenario: MemoryLimit field configures maximum memory
    Given a NovusPack package
    And a StreamConfig
    When MemoryLimit field is set
    Then maximum memory for buffering is configured
    And zero value means no memory limit
    And memory limit prevents excessive memory usage

  @REQ-STREAM-021 @happy
  Scenario: UseBufferPool field enables buffer pool integration
    Given a NovusPack package
    And a StreamConfig
    When UseBufferPool field is set
    Then global buffer pool integration is enabled or disabled
    And buffer pool reduces memory allocations
    And integration improves memory efficiency

  @REQ-STREAM-021 @happy
  Scenario: IsCompressed and IsEncrypted fields indicate data state
    Given a NovusPack package
    And a StreamConfig
    When IsCompressed and IsEncrypted fields are set
    Then IsCompressed indicates whether file is compressed
    And IsEncrypted indicates whether file is encrypted
    And state information enables proper data handling

  @REQ-STREAM-021 @error
  Scenario: StreamConfig struct validates configuration values
    Given a NovusPack package
    And a StreamConfig with invalid values
    When StreamConfig is used to create FileStream
    Then validation error is returned
    And error indicates invalid configuration field
    And error follows structured error format
