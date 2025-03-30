-- name: GetUserByEmailSQLC :one

SELECT  usr_email, usr_id FROM `pre_go_crm_user` WHERE usr_email = ? LIMIT 1;