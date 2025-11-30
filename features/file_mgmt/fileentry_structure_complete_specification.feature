@domain:file_mgmt @m2 @REQ-FILEMGMT-005 @spec(api_file_management.md#11-fileentry-structure)
Feature: FileEntry structure complete specification

  @happy
  Scenario: FileEntry static fields are correctly defined
    Given a FileEntry instance
    When static fields are examined
    Then FileID is 8 bytes (uint64)
    And OriginalSize and StoredSize are 8 bytes (uint64)
    And RawChecksum and StoredChecksum are 4 bytes (uint32)
    And FileVersion and MetadataVersion are 4 bytes (uint32)
    And PathCount is 2 bytes (uint16)
    And Type is 2 bytes (FileType)
    And CompressionType, CompressionLevel, EncryptionType are 1 byte (uint8)
    And HashCount is 1 byte (uint8)
    And HashDataOffset and OptionalDataOffset are 4 bytes (uint32)
    And HashDataLen and OptionalDataLen are 2 bytes (uint16)
    And Reserved is 4 bytes (uint32)
    And total static size is 64 bytes

  @happy
  Scenario: FileEntry variable-length data structures
    Given a FileEntry instance
    When variable-length data is examined
    Then Paths is a slice of PathEntry
    And Hashes is a slice of HashEntry
    And OptionalData is OptionalData structure
    And Tags is a pointer to tags slice

  @happy
  Scenario: FileEntry data management fields
    Given a FileEntry instance
    When data management fields are examined
    Then Data is byte slice for in-memory content
    And SourceFile is file handle for streaming
    And SourceOffset and SourceSize indicate source location
    And TempFilePath is path to temp file
    And IsDataLoaded indicates memory state
    And IsTempFile indicates temp file state
    And ProcessingState tracks current state

  @happy
  Scenario: FileEntry directory association fields
    Given a FileEntry instance
    When directory association is examined
    Then ParentDirectory points to parent directory
    And InheritedTags points to cached inherited tags
    And directory hierarchy is maintained

  @happy
  Scenario: PathEntry structure contains all fields
    Given a PathEntry instance
    When structure is examined
    Then Path is UTF-8 string
    And Mode is uint32 for permissions
    And UserID and GroupID are uint32
    And ModTime, CreateTime, AccessTime are time.Time
    And IsSymlink indicates symbolic link
    And LinkTarget contains symlink target

  @happy
  Scenario: HashEntry structure contains all fields
    Given a HashEntry instance
    When structure is examined
    Then Type is HashType
    And Purpose is HashPurpose
    And Data is byte slice containing hash

  @happy
  Scenario: OptionalData structure contains all optional fields
    Given an OptionalData instance
    When structure is examined
    Then Tags is slice of Tag
    And PathEncoding, PathFlags are pointers
    And CompressionDictID, SolidGroupID are pointers
    And FileSystemFlags, WindowsAttributes are pointers
    And ExtendedAttributes is map of strings
    And ACLData is byte slice
    And CustomData is map for reserved types
