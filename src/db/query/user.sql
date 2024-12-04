-- name: CreateUser :one
INSERT INTO users (
  email,
  hashed_password,
  phone,
  first_name,
  birth_date,
  gender,
  location_info,
  bio
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;