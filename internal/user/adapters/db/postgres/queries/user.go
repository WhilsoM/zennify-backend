package queries

const (
	UserCreate = `
INSERT INTO "user" (id, username, password_hash, created_at)
VALUES ($1, $2, $3, $4)
`

	UserGetByUsername = `
SELECT id, username, password_hash, created_at
FROM "user"
WHERE username = $1
`

	UserGetByID = `
SELECT id, username, password_hash, created_at
FROM "user"
WHERE id = $1
`
)
