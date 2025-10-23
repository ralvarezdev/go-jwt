package sqlite

const (
	// CreateRefreshTokensTableQuery is the SQL query to create the refresh_tokens table
	CreateRefreshTokensTableQuery = `
CREATE TABLE IF NOT EXISTS refresh_tokens (id TEXT PRIMARY KEY, expires_at DATETIME NOT NULL);
`

	// CreateAccessTokensTableQuery is the SQL query to create the access_tokens table
	CreateAccessTokensTableQuery = `
CREATE TABLE IF NOT EXISTS access_tokens (id TEXT PRIMARY KEY, parent_refresh_token_id TEXT, expires_at DATETIME NOT NULL);
`
)

var (
	// InsertRefreshTokenQuery is the SQL query to insert a new refresh token
	InsertRefreshTokenQuery = `
INSERT OR IGNORE INTO refresh_tokens (id, expires_at) VALUES (?, ?);
`

	// DeleteRefreshTokenQuery is the SQL query to delete a refresh token
	DeleteRefreshTokenQuery = `
DELETE FROM refresh_tokens WHERE id = ?;
`

	// CheckRefreshTokenQuery is the SQL query to check if a refresh token exists
	CheckRefreshTokenQuery = `
SELECT COUNT(1) FROM refresh_tokens WHERE id = ? AND expires_at > CURRENT_TIMESTAMP;
`

	// InsertAccessTokenQuery is the SQL query to insert a new access token
	InsertAccessTokenQuery = `
INSERT OR IGNORE INTO access_tokens (id, parent_refresh_token_id, expires_at) VALUES (?, ?, ?);
`

	// DeleteAccessTokenQuery is the SQL query to delete an access token
	DeleteAccessTokenQuery = `
DELETE FROM access_tokens WHERE id = ?;
`

	// DeleteAccessTokenByRefreshTokenQuery deletes access tokens by refresh token JTI
	DeleteAccessTokenByRefreshTokenQuery = `
DELETE FROM access_tokens WHERE parent_refresh_token_id = ?;
`

	// CheckAccessTokenQuery is the SQL query to check if an access token exists
	CheckAccessTokenQuery = `
SELECT COUNT(1) FROM access_tokens WHERE id = ? AND expires_at > CURRENT_TIMESTAMP;
`
)
