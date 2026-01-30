@domain:file_mgmt @m2 @REQ-FILEMGMT-079 @spec(api_file_mgmt_compression.md#71-file-compression-operations) @spec(api_file_mgmt_file_entry.md#fileentry-compression-purpose)
Feature: File Compression Operations
    As a package developer
    I want per-file compression and decompression capabilities
    So that individual files can be compressed independently

    Background:
        Given file compression operations support per-file compression
        And operations include CompressFile, DecompressFile, and GetFileCompressionInfo

    @REQ-FILEMGMT-079 @happy
    Scenario: CompressFile compresses an individual file
        When CompressFile is called on an uncompressed file
        Then the file content MUST be compressed using the specified algorithm
        And the FileEntry MUST be updated with compression metadata

    @REQ-FILEMGMT-079 @happy
    Scenario: CompressFile supports multiple compression algorithms
        When CompressFile is called with a compression type
        Then supported types MUST include Zstd, LZ4, and LZMA
        And the specified algorithm MUST be applied to the file content

    @REQ-FILEMGMT-079 @happy
    Scenario: DecompressFile decompresses an individual file
        When DecompressFile is called on a compressed file
        Then the file content MUST be decompressed
        And the FileEntry MUST be updated to reflect uncompressed state

    @REQ-FILEMGMT-079 @happy
    Scenario: GetFileCompressionInfo returns compression details
        When GetFileCompressionInfo is called on a file
        Then it MUST return compression type, compressed size, and uncompressed size
        And information MUST be accurate for both compressed and uncompressed files

    @REQ-FILEMGMT-079 @happy
    Scenario: Compression operations preserve file identity
        When a file is compressed or decompressed
        Then the FileID MUST remain unchanged
        And all file paths MUST remain associated with the file

    @REQ-FILEMGMT-079 @happy
    Scenario: Compression errors return structured errors
        When compression or decompression fails
        Then a structured PackageError MUST be returned
        And the error type MUST be ErrTypeCompression

    @REQ-FILEMGMT-079 @happy
    Scenario: Compressed files can be read without explicit decompression
        When ReadFile is called on a compressed file
        Then the content MUST be automatically decompressed
        And the caller receives uncompressed content transparently
