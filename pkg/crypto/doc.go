// Package crypto provides password hashing, JWT token generation/verification,
// and OTP (one-time password) utilities.
//
// Password hashing using bcrypt:
//
//	hash, err := crypto.EncodePassword("s3cret")
//	match := crypto.CheckPassword("s3cret", hash) // true
//
// JWT token generation and verification (HS256, requires 32-byte secret):
//
//	token, claims, err := crypto.GenerateToken(
//	    secretKey,    // 32-byte string
//	    userID,
//	    sessionID,
//	    &orgID,       // optional, can be nil
//	    24*time.Hour, // token duration
//	)
//
//	claims, err := crypto.VerifyToken(secretKey, token)
//	// claims.UserID, claims.SessionID, claims.OrgID, claims.ExpiresAt
//
// OTP generation and constant-time hash comparison (HMAC-SHA256):
//
//	code := crypto.GenNumericCode(6) // "483921"
//
//	hash := crypto.HashOTP(pepper, userID, "email_verification", code)
//	valid := crypto.EqualHashOTP(hash, storedHash) // constant-time comparison
package crypto
