@domain:file_format @m2 @REQ-FILEFMT-081 @spec(package_file_format.md#62-file-index-compression)
Feature: File Index Compression
    As a package format implementer
    I want file index compression capability
    So that index size is reduced and package efficiency is improved

    Background:
        Given the file index can be compressed
        And compression reduces index storage overhead

    @REQ-FILEFMT-081 @happy
    Scenario: File index compression uses LZ4 algorithm
        When the file index is compressed
        Then LZ4 compression MUST be used for the index
        And LZ4 provides fast decompression for quick package opening

    @REQ-FILEFMT-081 @happy
    Scenario: File index compression is applied as single block
        When the file index is compressed
        Then the entire index MUST be compressed as a single block
        And the compressed block MUST be decompressible independently

    @REQ-FILEFMT-081 @happy
    Scenario: Compressed index includes uncompressed size
        When the file index is compressed
        Then the compressed data MUST include the original uncompressed size
        And readers MUST verify decompressed size matches expected size

    @REQ-FILEFMT-081 @happy
    Scenario: File index compression reduces package size
        When packages contain large file indexes
        Then compression MUST reduce the index storage overhead
        And compression ratio SHOULD be significant for text-based index data

    @REQ-FILEFMT-081 @happy
    Scenario: File index decompression occurs during package open
        When a package with compressed index is opened
        Then the index MUST be decompressed during open operation
        And decompression errors MUST prevent package from opening

    @REQ-FILEFMT-081 @happy
    Scenario: Uncompressed indexes remain supported
        When a package has an uncompressed file index
        Then readers MUST support reading uncompressed indexes
        And no decompression MUST be attempted for uncompressed indexes
