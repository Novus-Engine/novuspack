# Security and Encryption Requirements

## Encryption Algorithms

- REQ-CRYPTO-001: Supported ciphers and key sizes are enforced for file encryption [type: constraint]. [security.md#32-encryption-algorithms](../tech_specs/security.md#32-encryption-algorithms)

## Quantum-Safe Encryption

- REQ-CRYPTO-002: Quantum-safe signature and key management options are available.
  [security.md#31-quantum-safe-encryption](../tech_specs/security.md#31-quantum-safe-encryption),
  [api_security.md#523-ml-kem-parameters](../tech_specs/api_security.md#523-ml-kem-parameters)
- REQ-CRYPTO-003: Quantum-safe keys are generated correctly with appropriate security levels.
  [security.md#341-ml-kem-key-management](../tech_specs/security.md#341-ml-kem-key-management),
  [api_security.md#525-ml-kem-error-conditions](../tech_specs/api_security.md#525-ml-kem-error-conditions),
  [api_security.md#4156-key-material-lifetime](../tech_specs/api_security.md#4156-key-material-lifetime)

## Runtime Secret Protection

- REQ-CRYPTO-004: All EncryptionKey keys must be treated as private key material [type: constraint].
  [api_security.md#4144-operation-requirements](../tech_specs/api_security.md#4144-operation-requirements),
  [api_security.md#4151-private-key-material-policy](../tech_specs/api_security.md#4151-private-key-material-policy)
- REQ-CRYPTO-005: GetKey must execute within runtime/secret.Do for all keys [type: constraint].
  [api_security.md#415-secure-encryptionkey-operations-with-runtimesecret](../tech_specs/api_security.md#415-secure-encryptionkey-operations-with-runtimesecret),
  [api_security.md#4141-getkey-behavior](../tech_specs/api_security.md#4141-getkey-behavior)
- REQ-CRYPTO-006: SetKey must execute within runtime/secret.Do for all keys [type: constraint].
  [api_security.md#415-secure-encryptionkey-operations-with-runtimesecret](../tech_specs/api_security.md#415-secure-encryptionkey-operations-with-runtimesecret),
  [api_security.md#4142-setkey-behavior](../tech_specs/api_security.md#4142-setkey-behavior)
- REQ-CRYPTO-007: Encrypt and Decrypt methods must execute within runtime/secret.Do when accessing keys [type: constraint].
  [api_security.md#415-secure-encryptionkey-operations-with-runtimesecret](../tech_specs/api_security.md#415-secure-encryptionkey-operations-with-runtimesecret),
  [api_security.md#4143-encryptionkey-error-conditions](../tech_specs/api_security.md#4143-encryptionkey-error-conditions)
- REQ-CRYPTO-008: File encryption operations must execute within runtime/secret.Do [type: constraint]. [api_security.md#465-secure-file-encryption-operations-with-runtimesecret](../tech_specs/api_security.md#465-secure-file-encryption-operations-with-runtimesecret), [api_security.md#464-packagegetfileencryptioninfo-method](../tech_specs/api_security.md#464-packagegetfileencryptioninfo-method)
- REQ-CRYPTO-009: Key generation operations must execute within runtime/secret.Do [type: constraint].
  [api_security.md#527-secure-encryption-operations-with-runtimesecret](../tech_specs/api_security.md#527-secure-encryption-operations-with-runtimesecret),
  [api_security.md#527-secure-encryption-operations-with-runtimesecret](../tech_specs/api_security.md#527-secure-encryption-operations-with-runtimesecret)
- REQ-CRYPTO-010: Key cleanup operations must execute within runtime/secret.Do [type: constraint]. [api_security.md#539-secure-key-clearing-with-runtimesecret](../tech_specs/api_security.md#539-secure-key-clearing-with-runtimesecret), [api_security.md#538-clear-behavior](../tech_specs/api_security.md#538-clear-behavior), [api_security.md#539-secure-key-clearing-with-runtimesecret](../tech_specs/api_security.md#539-secure-key-clearing-with-runtimesecret)

## Multi-Stage Transformation Pipeline Integration

- REQ-CRYPTO-011: Encryption operations integrate with multi-stage transformation pipelines as individual stages (encrypt stage writes to temporary file, decrypt stage reads from temporary file). [api_file_management.md#211-multi-stage-transformation-pipelines](../tech_specs/api_file_mgmt_transform_pipelines.md)
- REQ-CRYPTO-012: Encryption acts as pipeline stage in transformation sequence (e.g., compress => encrypt for addition, decrypt => decompress for extraction). [api_file_management.md#2113-processingstate-transitions](../tech_specs/api_file_mgmt_transform_pipelines.md)
- REQ-CRYPTO-013: Temporary files for encrypted content use context-aware security (encrypted on disk when possible, with exception for decrypt operations where user intends to decrypt). [api_file_management.md#2117-temporary-file-security](../tech_specs/api_file_mgmt_transform_pipelines.md)
- REQ-CRYPTO-014: Secure cleanup with overwrites for sensitive temporary files containing encrypted or decrypted content. [api_file_management.md#2117-temporary-file-security](../tech_specs/api_file_mgmt_transform_pipelines.md)
- REQ-CRYPTO-015: Context-aware decryption temp file handling allows unencrypted temp files when user is decrypting content to write to disk (user intent is to decrypt). [api_file_management.md#2117-temporary-file-security](../tech_specs/api_file_mgmt_transform_pipelines.md)
