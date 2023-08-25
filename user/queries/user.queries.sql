-- name: GetUser :one
SELECT *
FROM users
WHERE email = $1;

-- name: RegisterUser :exec
INSERT INTO users (
        email,
        full_name,
        intro,
        profile,
        role,
        password_hash
    )
VALUES ($1, $2, $3, $4, $5, $6);

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = $2
WHERE email = $1;