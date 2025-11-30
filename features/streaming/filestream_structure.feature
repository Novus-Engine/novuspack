@domain:streaming @m2 @REQ-STREAM-020 @spec(api_streaming.md#121-filestream-struct)
Feature: FileStream Structure

  @REQ-STREAM-020 @happy
  Scenario: FileStream struct provides file stream structure
    Given a NovusPack package
    When FileStream struct is examined
    Then struct contains reader field for underlying reader interface
    And struct contains file field for direct file access
    And struct contains bufReader field for buffered reading
    And struct contains fileSize field for total file size
    And struct contains position field for current read position
    And struct contains closed field for closure status
    And struct contains chunkSize field for read chunk size
    And struct contains config field for stream configuration
    And struct contains fileEntry field for associated file entry

  @REQ-STREAM-020 @happy
  Scenario: Reader field provides underlying reader interface
    Given a NovusPack package
    And a FileStream
    When stream operations are performed
    Then reader field provides io.Reader interface
    And reader enables reading from various sources
    And reader supports different input types

  @REQ-STREAM-020 @happy
  Scenario: File field enables direct file access
    Given a NovusPack package
    And a FileStream with file handle
    When direct file access is needed
    Then file field provides os.File handle
    And file handle enables direct file operations
    And file access supports efficient reading

  @REQ-STREAM-020 @happy
  Scenario: BufReader field provides buffered reading
    Given a NovusPack package
    And a FileStream
    When buffered reading is performed
    Then bufReader field provides bufio.Reader
    And buffered reader improves read efficiency
    And buffering reduces system call overhead

  @REQ-STREAM-020 @happy
  Scenario: FileSize and position fields track stream state
    Given a NovusPack package
    And a FileStream
    When stream state is examined
    Then fileSize field stores total size of file
    And position field stores current read position
    And state tracking enables stream information methods

  @REQ-STREAM-020 @happy
  Scenario: Closed field tracks stream closure status
    Given a NovusPack package
    And a FileStream
    When closure status is checked
    Then closed field indicates whether stream is closed
    And boolean field enables IsClosed method
    And closure tracking prevents operations on closed streams

  @REQ-STREAM-020 @happy
  Scenario: ChunkSize and config fields configure stream behavior
    Given a NovusPack package
    And a FileStream
    When stream behavior is configured
    Then chunkSize field determines read chunk size
    And config field provides StreamConfig for stream settings
    And configuration enables customizable stream behavior

  @REQ-STREAM-020 @happy
  Scenario: FileEntry field associates stream with package file
    Given a NovusPack package
    And a FileStream
    When file entry information is needed
    Then fileEntry field provides associated FileEntry
    And file entry enables access to file metadata
    And association links stream to package structure

  @REQ-STREAM-020 @error
  Scenario: FileStream struct handles invalid state correctly
    Given a NovusPack package
    And a FileStream with invalid state
    When stream operations are attempted
    Then structured error is returned
    And error indicates invalid state
    And error follows structured error format
