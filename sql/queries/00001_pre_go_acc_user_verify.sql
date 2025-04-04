-- name: GetValidOTP :one
SELECT verify_otp, verify_key_hash, verify_key, verify_id
FROM `pre_go_acc_user_verify`
WHERE verify_key_hash = ? AND is_verified = 0;

-- update lai
-- name: UpdateUserVerificationStatus :exec
UPDATE `pre_go_acc_user_verify`
SET is_verified = 1,
    verify_updated_at = now()
WHERE verify_key_hash = ?;

-- name: InsertOTPVerify :execresult
INSERT INTO `pre_go_acc_user_verify` (
    verify_otp,
    verify_key,
    verify_key_hash, 
    verify_type, 
    is_verified, 
    is_deleted, 
    verify_created_at, 
    verify_updated_at
)
VALUES (?, ?, ?, ?, 0, 0, NOW(), NOW());

-- name: InsertOrUpdateOTPVerify :execresult
INSERT INTO pre_go_acc_user_verify (
    verify_otp,
    verify_key,
    verify_key_hash, 
    verify_type, 
    is_verified, 
    is_deleted, 
    verify_created_at, 
    verify_updated_at
)
VALUES (?, ?, ?, ?, 0, 0, NOW(), NOW())
ON DUPLICATE KEY UPDATE 
    verify_otp = VALUES(verify_otp),
    verify_updated_at = NOW();


-- name: GetInfoOTP :one
SELECT verify_id, verify_otp, verify_key, verify_key_hash, verify_type, is_verified, is_deleted, verify_created_at, verify_updated_at
FROM `pre_go_acc_user_verify`
WHERE verify_key_hash = ?;
