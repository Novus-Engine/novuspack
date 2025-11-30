@domain:security @m2 @REQ-SEC-105 @spec(api_security.md#534-error-conditions)
Feature: ML-KEM Encryption Operations Error Conditions

  @REQ-SEC-105 @error
  Scenario: ML-KEM encryption operations return ErrEncryptionFailed on encryption failure
    Given an open NovusPack package
    And a valid context
    And ML-KEM key
    And data causing encryption failure
    When ML-KEM Encrypt is called
    Then ErrEncryptionFailed error is returned
    And error indicates encryption failure
    And error follows structured error format

  @REQ-SEC-105 @error
  Scenario: ML-KEM encryption operations return ErrDecryptionFailed on decryption failure
    Given an open NovusPack package
    And a valid context
    And ML-KEM key
    And corrupted or invalid ciphertext
    When ML-KEM Decrypt is called
    Then ErrDecryptionFailed error is returned
    And error indicates decryption failure
    And error follows structured error format

  @REQ-SEC-105 @error
  Scenario: ML-KEM encryption operations return ErrInvalidKey for invalid keys
    Given an open NovusPack package
    And a valid context
    And invalid or corrupted ML-KEM key
    When ML-KEM encryption operation is performed
    Then ErrInvalidKey error is returned
    And error indicates key is invalid or corrupted
    And error follows structured error format

  @REQ-SEC-105 @error
  Scenario: ML-KEM encryption operations return ErrContextCancelled on cancellation
    Given an open NovusPack package
    And a cancelled context
    And ML-KEM key and data
    When ML-KEM encryption operation is performed
    Then ErrContextCancelled error is returned
    And error indicates context was cancelled
    And error follows structured error format

  @REQ-SEC-105 @error
  Scenario: ML-KEM encryption operations return ErrContextTimeout on timeout
    Given an open NovusPack package
    And a context with timeout exceeded
    And ML-KEM key and data
    When ML-KEM encryption operation is performed
    Then ErrContextTimeout error is returned
    And error indicates context timeout exceeded
    And error follows structured error format

  @REQ-SEC-105 @error
  Scenario: ML-KEM encryption operations handle all error conditions consistently
    Given an open NovusPack package
    And a valid context
    And various error conditions
    When ML-KEM encryption operations encounter errors
    Then appropriate error type is returned for each condition
    And all errors follow structured error format
    And error messages provide specific information
