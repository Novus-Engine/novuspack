@domain:file_mgmt @m2 @REQ-FILEMGMT-312 @spec(api_file_mgmt_file_entry.md#2-fileentry-creation)
Feature: FileEntry Creation
    As a package developer
    I want FileEntry constructor and initialization methods
    So that file entries are created with proper defaults and tag synchronization

    Background:
        Given FileEntry creation methods are available
        And NewFileEntry provides proper initialization

    @REQ-FILEMGMT-312 @happy
    Scenario: NewFileEntry creates FileEntry with required fields
        When NewFileEntry is called with path and content
        Then it MUST return a FileEntry with initialized fields
        And Path, Size, and CRC32 fields MUST be set

    @REQ-FILEMGMT-312 @happy
    Scenario: NewFileEntry initializes tag storage
        When NewFileEntry creates a FileEntry
        Then tag storage MUST be initialized
        And tags MUST be synchronized between TaggedData and Tags fields

    @REQ-FILEMGMT-312 @happy
    Scenario: NewFileEntry sets default compression state
        When NewFileEntry creates a FileEntry without compression options
        Then IsCompressed MUST be false
        And CompressionType MUST be 0 (none)

    @REQ-FILEMGMT-312 @happy
    Scenario: NewFileEntry sets default encryption state
        When NewFileEntry creates a FileEntry without encryption options
        Then IsEncrypted MUST be false
        And EncryptionType MUST be 0 (none)

    @REQ-FILEMGMT-312 @happy
    Scenario: NewFileEntry generates unique FileID
        When NewFileEntry creates multiple FileEntry instances
        Then each FileEntry MUST have a unique FileID
        And FileID generation MUST be deterministic or random as appropriate

    @REQ-FILEMGMT-312 @happy
    Scenario: FileEntry creation validates path format
        When NewFileEntry is called with an invalid path
        Then it MUST return an error
        And the error MUST indicate the path validation failure
