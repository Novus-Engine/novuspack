@domain:file_format @m1 @REQ-FILEFMT-015 @spec(package_file_format.md#4142-path-entries)
Feature: Validate Path Entries and Per-Path Metadata

  @REQ-FILEFMT-015 @happy
  Scenario: Path entry structure contains all required fields
    Given a file entry with path entries
    When path entries are parsed
    Then each path entry has PathLength (2 bytes)
    And each path entry has UTF-8 Path string
    And each path entry has Mode (4 bytes)
    And each path entry has UserID (4 bytes)
    And each path entry has GroupID (4 bytes)
    And each path entry has ModTime (8 bytes)
    And each path entry has CreateTime (8 bytes)
    And each path entry has AccessTime (8 bytes)

  @REQ-FILEFMT-015 @happy
  Scenario: Primary path is first entry (index 0)
    Given a file entry with multiple path entries
    When path entries are parsed
    Then the first path entry (index 0) is the primary path
    And additional paths are secondary paths
    And all paths point to the same content

  @REQ-FILEFMT-015 @happy
  Scenario: PathLength matches actual UTF-8 path length
    Given a file entry with path entries
    And PathLength is specified for each path
    When path entries are parsed
    Then PathLength matches the actual UTF-8 encoded path length in bytes
    And path data is correctly bounded

  @REQ-FILEFMT-015 @happy
  Scenario: Path entries support per-path metadata
    Given a file entry with multiple paths
    When path entries are parsed
    Then each path can have different Mode
    And each path can have different UserID
    And each path can have different GroupID
    And each path can have different timestamps (ModTime, CreateTime, AccessTime)
    And per-path metadata is preserved

  @REQ-FILEFMT-015 @happy
  Scenario: Path strings are UTF-8 encoded and not null-terminated
    Given a file entry with path entries
    When path entries are parsed
    Then path strings are valid UTF-8 encoded
    And paths are not null-terminated
    And path length is determined by PathLength field

  @REQ-FILEFMT-015 @error
  Scenario: Invalid UTF-8 path bytes are rejected
    Given a path entry with invalid UTF-8 bytes
    When path entries are parsed
    Then a structured invalid path error is returned
    And error indicates invalid UTF-8 encoding
    And error follows structured error format

  @REQ-FILEFMT-015 @error
  Scenario: PathLength mismatch is rejected
    Given a path entry where PathLength does not match actual path length
    When path entries are parsed
    Then a structured invalid path error is returned
    And error indicates length mismatch
    And error follows structured error format

  @happy
  Scenario: WriteTo serializes path entry to binary format
    Given a PathEntry with values
    When path entry WriteTo is called with writer
    Then path entry is written to writer
    And written data matches path entry content

  @happy
  Scenario: ReadFrom deserializes path entry from binary format
    Given a reader with valid path entry data
    When path entry ReadFrom is called with reader
    Then path entry is read from reader
    And path entry fields match reader data
    And path entry is valid

  @happy
  Scenario: Path entry round-trip serialization preserves all fields
    Given a PathEntry with all fields set
    When path entry WriteTo is called with writer
    And ReadFrom is called with written data
    Then all path entry fields are preserved
    And path entry is valid

  @error
  Scenario: ReadFrom fails with incomplete data
    Given a reader with incomplete path entry data
    When path entry ReadFrom is called with reader
    Then structured IO error is returned
    And error indicates read failure
