"""
Go definitions index configuration.

Holds keyword mappings and priority phrases used during placement scoring.
"""

KEYWORD_TO_SECTION_MAPPING = {
    # Security and Encryption Domain
    "acl": [
        ("Encryption and Security Types", "strong"),
        ("Encryption and Security Methods", "strong"),
        ("Encryption and Security Helper Functions", "strong"),
        ("Security", "strong"),
    ],
    "access control": [
        ("Encryption and Security Types", "strong"),
        ("Encryption and Security Methods", "strong"),
        ("Encryption and Security Helper Functions", "strong"),
    ],
    "access control list": [
        ("Encryption and Security Types", "strong"),
        ("Encryption and Security Methods", "strong"),
        ("Encryption and Security Helper Functions", "strong"),
    ],
    "encryption": [
        ("Encryption and Security Types", "strong"),
        ("Encryption and Security Methods", "strong"),
        ("Encryption and Security Helper Functions", "strong"),
    ],
    "encrypt": [
        ("Encryption and Security Types", "medium"),
        ("Encryption and Security Methods", "medium"),
        ("Encryption and Security Helper Functions", "medium"),
    ],
    "decrypt": [
        ("Encryption and Security Types", "medium"),
        ("Encryption and Security Methods", "medium"),
        ("Encryption and Security Helper Functions", "medium"),
    ],
    "security": [
        ("Encryption and Security Types", "strong"),
        ("Encryption and Security Methods", "strong"),
        ("Encryption and Security Helper Functions", "strong"),
    ],
    "mlkem": [
        ("Encryption and Security Types", "strong"),
        ("Encryption and Security Methods", "strong"),
        ("Encryption and Security Helper Functions", "strong"),
    ],
    "aes": [
        ("Encryption and Security Types", "medium"),
        ("Encryption and Security Methods", "medium"),
        ("Encryption and Security Helper Functions", "medium"),
    ],
    "chacha": [
        ("Encryption and Security Types", "medium"),
        ("Encryption and Security Methods", "medium"),
        ("Encryption and Security Helper Functions", "medium"),
    ],
    "cipher": [
        ("Encryption and Security Types", "medium"),
        ("Encryption and Security Methods", "medium"),
        ("Encryption and Security Helper Functions", "medium"),
    ],
    # Error Handling Domain
    "error": [
        ("Error Types", "strong"),
        ("Error Methods", "strong"),
        ("Error Helper Functions", "strong"),
    ],
    "error context": [
        ("Error Types", "strong"),
        ("Error Methods", "strong"),
        ("Error Helper Functions", "strong"),
    ],
    "err": [
        ("Error Types", "strong"),
        ("Error Methods", "strong"),
        ("Error Helper Functions", "strong"),
    ],
    "packageerror": [
        ("Error Types", "strong"),
        ("Error Methods", "strong"),
        ("Error Helper Functions", "strong"),
    ],
    "validation": [
        ("Error Types", "medium"),
        ("Error Methods", "medium"),
        ("Error Helper Functions", "medium"),
        ("Security Validation", "strong"),
    ],
    "validate": [
        ("Error Types", "medium"),
        ("Error Methods", "medium"),
        ("Error Helper Functions", "medium"),
        ("Security Validation", "strong"),
    ],
    "verify": [
        ("Error Types", "medium"),
        ("Error Methods", "medium"),
        ("Error Helper Functions", "medium"),
    ],
    # Metadata and Tags Domain
    "tag": [
        ("Tag Methods", "strong"),
        ("FileEntry Types", "medium"),
        ("FileEntry Helper Functions", "medium"),
    ],
    "fileentrytag": [
        ("FileEntry Helper Functions", "strong"),
    ],
    "pathmetadatatag": [
        ("Package Path Metadata Methods", "strong"),
        ("Package Metadata Type Methods", "strong"),
        ("Package Metadata Helper Functions", "strong"),
    ],
    "metadata": [
        ("Package Metadata Types", "strong"),
        ("Package Comment Methods", "strong"),
        ("Package Identity Methods", "strong"),
        ("Package Special File Methods", "strong"),
        ("Package Path Metadata Methods", "strong"),
        ("Package Symlink Methods", "strong"),
        ("Package Metadata-Only Methods", "strong"),
        ("Package Info Methods", "strong"),
        ("Package Metadata Validation Methods", "strong"),
        ("Package Metadata Internal Methods", "strong"),
        ("Package Metadata Type Methods", "strong"),
        ("Package Metadata Helper Functions", "strong"),
    ],
    "pathmetadata": [
        ("Package Metadata Types", "strong"),
        ("Package Path Metadata Methods", "strong"),
        ("Package Metadata Type Methods", "strong"),
        ("Package Metadata Helper Functions", "strong"),
    ],
    "fileentry": [
        ("FileEntry Types", "strong"),
        ("FileEntry Query Methods", "strong"),
        ("FileEntry Data Methods", "strong"),
        ("FileEntry Temp File Methods", "strong"),
        ("FileEntry Serialization Methods", "strong"),
        ("FileEntry Path Methods", "strong"),
        ("FileEntry Transformation Methods", "strong"),
        ("FileEntry Helper Functions", "strong"),
    ],
    "appid": [
        ("Package Identity Methods", "medium"),
    ],
    "vendorid": [
        ("Package Identity Methods", "medium"),
    ],
    "comment": [
        ("Package Comment Methods", "strong"),
        ("Package Metadata Type Methods", "strong"),
        ("Package Metadata Types", "medium"),
    ],
    "packagecomment": [
        ("Package Comment Methods", "strong"),
        ("Package Metadata Type Methods", "strong"),
        ("Package Metadata Types", "medium"),
    ],
    "validatecomment": [
        ("Package Helper Functions", "strong"),
        ("Package Metadata Helper Functions", "medium"),
        ("Package Comment Methods", "medium"),
    ],
    "validatepathlength": [
        ("Package Helper Functions", "strong"),
    ],
    # Streaming and Buffer Domain
    "streaming": [
        ("Streaming and Buffer Types", "strong"),
        ("Streaming and Buffer Methods", "strong"),
        ("Streaming and Buffer Helper Functions", "strong"),
    ],
    "stream": [
        ("Streaming and Buffer Types", "medium"),
        ("Streaming and Buffer Methods", "medium"),
        ("Streaming and Buffer Helper Functions", "medium"),
    ],
    "buffer": [
        ("Streaming and Buffer Types", "strong"),
        ("Streaming and Buffer Methods", "strong"),
        ("Streaming and Buffer Helper Functions", "strong"),
    ],
    "bufferpool": [
        ("Streaming and Buffer Types", "strong"),
        ("Streaming and Buffer Methods", "strong"),
        ("Streaming and Buffer Helper Functions", "strong"),
    ],
    "chunk": [
        ("Streaming and Buffer Types", "medium"),
        ("Streaming and Buffer Methods", "medium"),
        ("Streaming and Buffer Helper Functions", "medium"),
    ],
    # Compression Domain
    "compression": [
        ("Compression Types", "strong"),
        ("Compression Methods", "strong"),
        ("Compression Helper Functions", "strong"),
    ],
    "compress": [
        ("Compression Types", "medium"),
        ("Compression Methods", "medium"),
        ("Compression Helper Functions", "medium"),
    ],
    "decompress": [
        ("Compression Types", "medium"),
        ("Compression Methods", "medium"),
        ("Compression Helper Functions", "medium"),
    ],
    # Signature Domain
    "signature": [
        ("Signature Types", "strong"),
        ("Signature Methods", "strong"),
        ("Signature Helper Functions", "strong"),
    ],
    "sign": [
        ("Signature Types", "medium"),
        ("Signature Methods", "medium"),
        ("Signature Helper Functions", "medium"),
    ],
    "signing": [
        ("Signature Types", "medium"),
        ("Signature Methods", "medium"),
        ("Signature Helper Functions", "medium"),
    ],
    # Deduplication Domain
    "deduplication": [
        ("Package File Management Methods", "strong"),
        ("Package Information and Queries Methods", "strong"),
    ],
    "dedup": [
        ("Package File Management Methods", "medium"),
        ("Package Information and Queries Methods", "medium"),
    ],
    # FileType System Domain
    "filetype": [
        ("FileType System Types", "strong"),
        ("FileType System Methods", "strong"),
        ("FileType System Helper Functions", "strong"),
    ],
    "file type": [
        ("FileType System Types", "strong"),
        ("FileType System Methods", "strong"),
        ("FileType System Helper Functions", "strong"),
    ],
    "mimetype": [
        ("FileType System Types", "medium"),
        ("FileType System Methods", "medium"),
        ("FileType System Helper Functions", "medium"),
    ],
    # Generic Types Domain
    "generic": [
        ("Generic Types", "strong"),
        ("Generic Methods", "strong"),
        ("Generic Helper Functions", "strong"),
    ],
    "option": [
        ("Generic Types", "medium"),
        ("Generic Methods", "medium"),
        ("Generic Helper Functions", "medium"),
    ],
    "result": [
        ("Generic Types", "medium"),
        ("Generic Methods", "medium"),
        ("Generic Helper Functions", "medium"),
    ],
}

# Hardcoded priority phrases for phrase matching (strong weight)
PRIORITY_PHRASES = [
    "error context",
    "packageerror",
    "fileentry",
    "path metadata",
    "access control",
    "access control list",
    "fileentry tag",
    "path metadata tag",
]

DOMAIN_FILE_MAP = {
    # Deprecated: replaced by pattern-based detection in
    # lib/go_defs_index/_go_defs_index_scoring_domain.py.
    # Kept temporarily for reference until fully removed.
    "api_generics.md": "generic",
    "api_streaming.md": "streaming",
    "api_package_compression.md": "compression",
    "api_security.md": "encryption",
    "api_metadata.md": "metadata",
    "api_deduplication.md": "deduplication",
    "api_signatures.md": "signature",
    "file_type_system.md": "filetype",
    "api_writing.md": "writing",
}
