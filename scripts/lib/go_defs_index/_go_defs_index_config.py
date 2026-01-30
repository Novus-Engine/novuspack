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
        ("FileEntry Helper Functions", "medium"),
        ("Metadata Types", "medium"),
        ("Metadata Methods", "strong"),
        ("Metadata Helper Functions", "strong"),
    ],
    "fileentrytag": [
        ("FileEntry Helper Functions", "strong"),
        ("Metadata Methods", "medium"),
        ("Metadata Helper Functions", "medium"),
    ],
    "pathmetadatatag": [
        ("Metadata Methods", "strong"),
        ("Metadata Helper Functions", "strong"),
    ],
    "metadata": [
        ("Metadata Types", "strong"),
        ("Metadata Methods", "strong"),
        ("Metadata Helper Functions", "strong"),
        ("Package Metadata Methods", "medium"),
    ],
    "pathmetadata": [
        ("Metadata Types", "strong"),
        ("Metadata Methods", "strong"),
        ("Metadata Helper Functions", "strong"),
    ],
    "fileentry": [
        ("FileEntry", "strong"),
        ("FileEntry Methods", "strong"),
        ("FileEntry Helper Functions", "strong"),
    ],
    "appid": [
        ("Package Metadata Methods", "medium"),
    ],
    "vendorid": [
        ("Package Metadata Methods", "medium"),
    ],
    "comment": [
        ("Package Metadata Methods", "strong"),
        ("Metadata Types", "medium"),
    ],
    "packagecomment": [
        ("Package Metadata Methods", "strong"),
        ("Metadata Types", "medium"),
    ],
    "validatecomment": [
        ("Package Helper Functions", "strong"),
        ("Package Metadata Methods", "medium"),
        ("Package Metadata Methods", "medium"),
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
        ("Deduplication Types", "strong"),
        ("Deduplication Methods", "strong"),
        ("Deduplication Helper Functions", "strong"),
    ],
    "dedup": [
        ("Deduplication Types", "medium"),
        ("Deduplication Methods", "medium"),
        ("Deduplication Helper Functions", "medium"),
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
