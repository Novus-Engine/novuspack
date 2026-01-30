# Digital Signature API Requirements

## Status

Signature management and signature validation are deferred to v2.
V1 only enforces signed package immutability and supports clear-signatures workflows via the writing API.
Unless noted otherwise, requirements in this document are v2 future work.

## Signature Management

- REQ-SIG-001: Multiple signature types supported including quantum-safe. [api_signatures.md#2-signature-types](../tech_specs/api_signatures.md#2-signature-types)
- REQ-SIG-004: AddSignature adds new digital signature to package. [api_signatures.md#1-signature-management](../tech_specs/api_signatures.md#1-signature-management)
- REQ-SIG-005: RemoveSignature removes signature by index and all later signatures. [api_signatures.md#1-signature-management](../tech_specs/api_signatures.md#1-signature-management)
- REQ-SIG-006: GetSignatureCount returns total number of signatures. [api_signatures.md#1-signature-management](../tech_specs/api_signatures.md#1-signature-management)
- REQ-SIG-007: GetSignature retrieves signature information by index. [api_signatures.md#1-signature-management](../tech_specs/api_signatures.md#1-signature-management)
- REQ-SIG-008: GetAllSignatures returns all signatures in package. [api_signatures.md#1-signature-management](../tech_specs/api_signatures.md#1-signature-management)
- REQ-SIG-009: ClearAllSignatures removes all signatures. [api_signatures.md#1-signature-management](../tech_specs/api_signatures.md#1-signature-management)
- REQ-SIG-020: Implementation requirements define signature implementation needs [type: architectural]. [api_signatures.md#117-implementation-requirements](../tech_specs/api_signatures.md#117-implementation-requirements)
- REQ-SIG-021: Incremental signing implementation provides sequential signature support. [api_signatures.md#12-incremental-signing-implementation](../tech_specs/api_signatures.md#12-incremental-signing-implementation)
- REQ-SIG-022: Adding subsequent signatures supports incremental signing. [api_signatures.md#121-adding-subsequent-signatures](../tech_specs/api_signatures.md#121-adding-subsequent-signatures)
- REQ-SIG-023: Signature validation process defines validation workflow. [api_signatures.md#122-signature-validation-process](../tech_specs/api_signatures.md#122-signature-validation-process)
- REQ-SIG-024: Key implementation points define critical implementation details [type: architectural]. [api_signatures.md#123-key-implementation-points](../tech_specs/api_signatures.md#123-key-implementation-points)
- REQ-SIG-060: 1.1 Multiple Signature Management (Incremental Signing) is specified and implemented. [api_signatures.md#11-multiple-signature-management-incremental-signing](../tech_specs/api_signatures.md#11-multiple-signature-management-incremental-signing)

## Signature Types

- REQ-SIG-025: Signature type constants define supported signature types. [api_signatures.md#21-signature-type-constants](../tech_specs/api_signatures.md#21-signature-type-constants)
- REQ-SIG-029: ML-DSA implementation provides quantum-safe signature support. [api_signatures.md#23-ml-dsa-crystals-dilithium-implementation](../tech_specs/api_signatures.md#23-ml-dsa-crystals-dilithium-implementation)
- REQ-SIG-030: SLH-DSA implementation provides quantum-safe hash-based signatures. [api_signatures.md#24-slh-dsa-sphincs-implementation](../tech_specs/api_signatures.md#24-slh-dsa-sphincs-implementation)
- REQ-SIG-031: PGP implementation provides OpenPGP signature support. [api_signatures.md#25-pgp-openpgp-implementation](../tech_specs/api_signatures.md#25-pgp-openpgp-implementation)
- REQ-SIG-032: X.509/PKCS#7 implementation provides certificate-based signatures. [api_signatures.md#26-x509pkcs7-implementation](../tech_specs/api_signatures.md#26-x509pkcs7-implementation)

## Signature Structures

- REQ-SIG-026: Signature information structure defines signature metadata format. [api_signatures.md#22-signature-information-structure](../tech_specs/api_signatures.md#22-signature-information-structure)
- REQ-SIG-027: SignatureInfo struct provides signature information structure. [api_signatures.md#221-signatureinfo-struct](../tech_specs/api_signatures.md#221-signatureinfo-struct)
- REQ-SIG-028: SignatureValidationResult struct provides validation result structure. [api_signatures.md#222-signaturevalidationresult-struct](../tech_specs/api_signatures.md#222-signaturevalidationresult-struct)

## Signature Validation

- REQ-SIG-002: Validation returns detailed status per signature. [api_signatures.md#27-signature-validation](../tech_specs/api_signatures.md#27-signature-validation)
- REQ-SIG-010: ValidateAllSignatures validates all signatures in order. [api_signatures.md#27-signature-validation](../tech_specs/api_signatures.md#27-signature-validation)
- REQ-SIG-011: ValidateSignatureType validates signatures of specific type. [api_signatures.md#27-signature-validation](../tech_specs/api_signatures.md#27-signature-validation)
- REQ-SIG-012: ValidateSignatureIndex validates signature by index. [api_signatures.md#27-signature-validation](../tech_specs/api_signatures.md#27-signature-validation)
- REQ-SIG-013: ValidateSignatureWithKey validates signature with specific public key. [api_signatures.md#27-signature-validation](../tech_specs/api_signatures.md#27-signature-validation)
- REQ-SIG-014: ValidateSignatureChain validates signature chain integrity. [api_signatures.md#27-signature-validation](../tech_specs/api_signatures.md#27-signature-validation)

## Immutability

- REQ-SIG-003: Post-sign write operations are blocked [type: constraint] (v1). [api_signatures.md#13-immutability-check](../tech_specs/api_signatures.md#13-immutability-check)

## Existing Package Signing

- REQ-SIG-033: Existing package signing supports signing existing packages. [api_signatures.md#28-existing-package-signing](../tech_specs/api_signatures.md#28-existing-package-signing)
- REQ-SIG-034: Implementation requirements for existing packages define signing needs [type: architectural]. [api_signatures.md#285-implementation-requirements-signpackage](../tech_specs/api_signatures.md#285-implementation-requirements-signpackage)
- REQ-SIG-035: Function usage guide provides signature function guidance [type: documentation-only]. [api_signatures.md#286-function-usage-guide](../tech_specs/api_signatures.md#286-function-usage-guide)
- REQ-SIG-036: When to use AddSignature defines low-level usage [type: documentation-only]. [api_signatures.md#2861-when-to-use-addsignature-low-level](../tech_specs/api_signatures.md#2861-when-to-use-addsignature-low-level)
- REQ-SIG-037: When to use SignPackage functions defines high-level usage [type: documentation-only]. [api_signatures.md#2862-when-to-use-signpackage-functions-high-level](../tech_specs/api_signatures.md#2862-when-to-use-signpackage-functions-high-level)
- REQ-SIG-038: Implementation pattern provides signature implementation guidance [type: documentation-only]. [api_signatures.md#2863-implementation-pattern](../tech_specs/api_signatures.md#2863-implementation-pattern)
- REQ-SIG-068: Secure signing operations using key material MUST execute within `runtime/secret.Do` for signature generation and key-loading workflows. [api_signatures.md#2851-secure-signing-operations-with-runtimesecret](../tech_specs/api_signatures.md#2851-secure-signing-operations-with-runtimesecret)

## Key Management

- REQ-SIG-039: Signing key management provides key handling operations. [api_signatures.md#413-signingkey-structure](../tech_specs/api_signatures.md#413-signingkey-structure)

## Generic Signature Patterns

- REQ-SIG-046: Generic signature patterns provide type-safe signature support. [api_signatures.md#4-generic-signature-patterns](../tech_specs/api_signatures.md#4-generic-signature-patterns)
- REQ-SIG-047: Generic signature strategy interface defines signature strategy contract. [api_signatures.md#41-signaturestrategy-interface](../tech_specs/api_signatures.md#41-signaturestrategy-interface)
- REQ-SIG-048: Generic signature configuration provides type-safe configuration. [api_signatures.md#43-generic-signature-configuration](../tech_specs/api_signatures.md#43-generic-signature-configuration)
- REQ-SIG-049: Generic signature validation provides type-safe validation. [api_signatures.md#44-generic-signature-validation](../tech_specs/api_signatures.md#44-generic-signature-validation)
- REQ-SIG-061: Signature typed representation uses Option[T] as source of truth [type: architectural]. [api_signatures.md#221-signatureinfo-struct](../tech_specs/api_signatures.md#221-signatureinfo-struct)
- REQ-SIG-062: Signature GetData returns (T, error) not (T, bool) [type: constraint]. [api_signatures.md#221-signatureinfo-struct](../tech_specs/api_signatures.md#221-signatureinfo-struct)
- REQ-SIG-063: Signature.IsValid semantics are intentionally incomplete for v1 [type: documentation-only]. [api_signatures.md#42-isvalid-and-isexpired-semantics](../tech_specs/api_signatures.md#42-isvalid-and-isexpired-semantics)
- REQ-SIG-064: All SigningKey keys must be treated as private key material [type: constraint]. [api_signatures.md#416-signaturestrategy-secure-signingkey-operations-with-runtimesecret](../tech_specs/api_signatures.md#416-signaturestrategy-secure-signingkey-operations-with-runtimesecret)
- REQ-SIG-065: SigningKey GetKey must execute within runtime/secret.Do for all keys [type: constraint]. [api_signatures.md#416-signaturestrategy-secure-signingkey-operations-with-runtimesecret](../tech_specs/api_signatures.md#416-signaturestrategy-secure-signingkey-operations-with-runtimesecret)
- REQ-SIG-066: SigningKey SetKey must execute within runtime/secret.Do for all keys [type: constraint]. [api_signatures.md#416-signaturestrategy-secure-signingkey-operations-with-runtimesecret](../tech_specs/api_signatures.md#416-signaturestrategy-secure-signingkey-operations-with-runtimesecret)
- REQ-SIG-067: Signing operations must execute within runtime/secret.Do when accessing keys [type: constraint]. [api_signatures.md#416-signaturestrategy-secure-signingkey-operations-with-runtimesecret](../tech_specs/api_signatures.md#416-signaturestrategy-secure-signingkey-operations-with-runtimesecret)

- REQ-SIG-069: SigningKey GetKey returns a copy (deep copy for slices), returns error when missing/invalid/expired, and should be used within runtime/secret.Do. [api_signatures.md#4151-getkey-behavior](../tech_specs/api_signatures.md#4151-getkey-behavior)
- REQ-SIG-070: SigningKey SetKey overwrites existing key with a copy, does not update timestamps, and validates key type and expiration. [api_signatures.md#4152-setkey-behavior](../tech_specs/api_signatures.md#4152-setkey-behavior)
- REQ-SIG-071: SigningKey GetKey and SetKey error conditions return `*PackageError` with `ErrTypeSignature` and descriptive messages for key-not-set, invalid, expired, and type mismatch. [api_signatures.md#4153-error-conditions](../tech_specs/api_signatures.md#4153-error-conditions)
- REQ-SIG-072: Operations using SigningKey re-validate key state before use and return `*PackageError` with `ErrTypeSignature` and context on failure. [api_signatures.md#4154-operation-requirements](../tech_specs/api_signatures.md#4154-operation-requirements)

- REQ-SIG-073: SigningKey IsValid returns true only when key is set, KeyID non-empty, CreatedAt non-zero, KeyType is valid, and ExpiresAt (if set) is after CreatedAt. [api_signatures.md#421-isvalid-requirements](../tech_specs/api_signatures.md#421-isvalid-requirements)
- REQ-SIG-074: SigningKey IsExpired returns true only when ExpiresAt is set and current time is at or after ExpiresAt. [api_signatures.md#422-isexpired-requirements](../tech_specs/api_signatures.md#422-isexpired-requirements)
- REQ-SIG-075: SigningKey expiration semantics define ExpiresAt=nil means never expires, ExpiresAt <= CreatedAt is invalid, and CreatedAt zero makes key invalid. [api_signatures.md#423-expiration-semantics](../tech_specs/api_signatures.md#423-expiration-semantics)
- REQ-SIG-076: Validation order requires checking IsValid first, then IsExpired, and both must pass for the key to be usable. [api_signatures.md#424-validation-order](../tech_specs/api_signatures.md#424-validation-order)

## Error Handling

- REQ-SIG-050: Error handling provides signature error management. [api_signatures.md#5-error-handling](../tech_specs/api_signatures.md#5-error-handling)
- REQ-SIG-051: Structured error system defines signature error types [type: architectural]. [api_signatures.md#51-structured-error-system](../tech_specs/api_signatures.md#51-structured-error-system)
- REQ-SIG-052: Common signature error types define standard error classifications. [api_signatures.md#52-signature-error-messages](../tech_specs/api_signatures.md#52-signature-error-messages)
- REQ-SIG-053: Specific signature error types define signature-specific errors. [api_signatures.md#53-signature-specific-error-context-types](../tech_specs/api_signatures.md#53-signature-specific-error-context-types)
- REQ-SIG-054: Error type mapping maps legacy to structured errors. [api_signatures.md#52-signature-error-messages](../tech_specs/api_signatures.md#52-signature-error-messages)
- REQ-SIG-055: Structured error examples demonstrate error handling patterns. [api_signatures.md#53-signature-specific-error-context-types](../tech_specs/api_signatures.md#53-signature-specific-error-context-types)
- REQ-SIG-077: ValidationErrorContext structure defines signature validation error context fields. [api_signatures.md#534-validationerrorcontext-structure](../tech_specs/api_signatures.md#534-validationerrorcontext-structure)
- REQ-SIG-056: Creating signature errors supports structured error creation. [api_signatures.md#53-signature-specific-error-context-types](../tech_specs/api_signatures.md#53-signature-specific-error-context-types)
- REQ-SIG-057: Error handling patterns define recommended error handling [type: documentation-only]. [api_signatures.md#52-signature-error-messages](../tech_specs/api_signatures.md#52-signature-error-messages)
- REQ-SIG-058: Function signatures define error handling functions. [api_signatures.md#51-structured-error-system](../tech_specs/api_signatures.md#51-structured-error-system)
- REQ-SIG-059: Core error handling functions provide error management utilities. [api_signatures.md#51-structured-error-system](../tech_specs/api_signatures.md#51-structured-error-system)

## Industry Standards and Comparison

- REQ-SIG-040: Comparison with other implementations compares signature systems [type: documentation-only]. [api_signatures.md#3-comparison-with-other-signed-file-implementations](../tech_specs/api_signatures.md#3-comparison-with-other-signed-file-implementations)
- REQ-SIG-041: Industry standard comparison compares with industry standards [type: documentation-only]. [api_signatures.md#31-industry-standard-comparison](../tech_specs/api_signatures.md#31-industry-standard-comparison)
- REQ-SIG-042: NovusPack advantages document signature advantages [type: documentation-only]. [api_signatures.md#32-novuspack-advantages](../tech_specs/api_signatures.md#32-novuspack-advantages)
- REQ-SIG-043: Industry standard compliance ensures standards alignment [type: architectural]. [api_signatures.md#33-industry-standard-compliance](../tech_specs/api_signatures.md#33-industry-standard-compliance)
- REQ-SIG-044: Signature size comparison compares signature sizes [type: non-functional]. [api_signatures.md#34-signature-size-comparison](../tech_specs/api_signatures.md#34-signature-size-comparison)
- REQ-SIG-045: Verification performance defines validation performance characteristics [type: non-functional]. [api_signatures.md#35-verification-performance](../tech_specs/api_signatures.md#35-verification-performance)

## Context Integration

- REQ-SIG-015: All signature methods accept context.Context and respect cancellation/timeout [type: constraint]. [api_core.md#02-context-integration](../tech_specs/api_core.md#02-context-integration)
- REQ-SIG-019: Context cancellation during signature operations stops operation and returns context error. [api_core.md#02-context-integration](../tech_specs/api_core.md#02-context-integration)

## Validation

- REQ-SIG-016: Signature type parameters validated against supported signature types [type: constraint]. [api_signatures.md#2-signature-types](../tech_specs/api_signatures.md#2-signature-types)
- REQ-SIG-017: Signature index parameters validated (non-negative, within signature count) [type: constraint]. [api_signatures.md#1-signature-management](../tech_specs/api_signatures.md#1-signature-management)
- REQ-SIG-018: Public key parameters validated (non-nil, valid format for signature type) [type: constraint]. [api_signatures.md#27-signature-validation](../tech_specs/api_signatures.md#27-signature-validation)
