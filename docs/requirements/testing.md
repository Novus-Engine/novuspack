# Testing Requirements

## Testing Infrastructure

- REQ-TEST-001: Presence checks for @spec and @REQ tags in features. [testing.md#0-overview](../tech_specs/testing.md#0-overview)
- REQ-TEST-002: Cross-check tech spec scenarios against feature files by domain; fail if any scenario lacks a mapped feature or any feature references a missing spec anchor. [testing.md#0-overview](../tech_specs/testing.md#0-overview)
- REQ-TEST-003: Testing coverage targets are defined per domain [type: documentation-only]. [testing.md#0-overview](../tech_specs/testing.md#0-overview)

## Dual Encryption Testing

- REQ-TEST-004: Dual encryption testing requirements define encryption testing needs [type: documentation-only]. [testing.md#1-dual-encryption-testing-requirements](../tech_specs/testing.md#1-dual-encryption-testing-requirements)
- REQ-TEST-005: ML-KEM encryption testing validates quantum-safe encryption. [testing.md#11-ml-kem-encryption-testing](../tech_specs/testing.md#11-ml-kem-encryption-testing)
- REQ-TEST-006: AES-256-GCM encryption testing validates traditional encryption. [testing.md#12-aes-256-gcm-encryption-testing](../tech_specs/testing.md#12-aes-256-gcm-encryption-testing)
- REQ-TEST-007: Dual encryption integration testing validates combined encryption. [testing.md#13-dual-encryption-integration-testing](../tech_specs/testing.md#13-dual-encryption-integration-testing)

## File Validation Testing

- REQ-TEST-008: File validation testing requirements define validation testing needs [type: documentation-only]. [testing.md#2-file-validation-testing-requirements](../tech_specs/testing.md#2-file-validation-testing-requirements)
- REQ-TEST-009: Empty file testing validates empty file handling. [testing.md#21-empty-file-testing](../tech_specs/testing.md#21-empty-file-testing)
- REQ-TEST-010: Path normalization testing validates path normalization. [testing.md#22-path-normalization-testing](../tech_specs/testing.md#22-path-normalization-testing)
- REQ-TEST-011: Compression error handling testing validates error handling. [testing.md#23-compression-error-handling-testing](../tech_specs/testing.md#23-compression-error-handling-testing)
- REQ-TEST-012: Hash-based deduplication testing validates deduplication. [testing.md#24-hash-based-deduplication-testing](../tech_specs/testing.md#24-hash-based-deduplication-testing)
