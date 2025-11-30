@domain:streaming @m2 @REQ-STREAM-042 @spec(api_streaming.md#221-bufferpool-struct)
Feature: BufferPool Structure

  @REQ-STREAM-042 @happy
  Scenario: BufferPool struct provides buffer pool structure
    Given a NovusPack package
    When BufferPool struct is examined
    Then struct contains buffers map for buffer ID to buffer data mapping
    And struct contains lastUsed map for last access time tracking
    And struct contains accessCount map for access count tracking
    And struct provides complete buffer pool state management

  @REQ-STREAM-042 @happy
  Scenario: Buffers map stores buffer data by ID
    Given a NovusPack package
    And a BufferPool
    When buffers are allocated and stored
    Then buffers map associates buffer IDs with buffer data
    And map enables efficient buffer retrieval by ID
    And buffer data is stored as byte slices

  @REQ-STREAM-042 @happy
  Scenario: LastUsed map tracks buffer access times
    Given a NovusPack package
    And a BufferPool
    When buffers are accessed
    Then lastUsed map records last access time for each buffer
    And time tracking enables LRU eviction policy
    And access time is stored as time.Time

  @REQ-STREAM-042 @happy
  Scenario: AccessCount map tracks buffer usage statistics
    Given a NovusPack package
    And a BufferPool
    When buffers are accessed multiple times
    Then accessCount map maintains access count for each buffer
    And count enables access pattern analysis
    And statistics support buffer pool optimization

  @REQ-STREAM-042 @happy
  Scenario: BufferPool struct manages buffer lifecycle
    Given a NovusPack package
    And a BufferPool with configuration
    When buffers are allocated and released
    Then struct tracks all buffer state through maps
    And lifecycle management enables efficient buffer reuse
    And state tracking supports eviction policies

  @REQ-STREAM-042 @error
  Scenario: BufferPool struct handles concurrent access safely
    Given a NovusPack package
    And a BufferPool with concurrent access
    When multiple goroutines access struct fields
    Then all map operations are properly synchronized
    And thread safety prevents data races
    And concurrent access is handled correctly
