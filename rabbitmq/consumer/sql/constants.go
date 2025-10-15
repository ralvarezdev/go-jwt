package sql

const (
	// CreateRefreshTokensTableQuery is the SQL query to create the refresh_tokens table
	CreateRefreshTokensTableQuery = `CREATE TABLE IF NOT EXISTS refresh_tokens (jti TEXT PRIMARY KEY)`

	// InsertRefreshTokenQuery is the SQL query to insert a new refresh token
	InsertRefreshTokenQuery = `INSERT OR IGNORE INTO refresh_tokens (jti) VALUES (?)`

	// DeleteRefreshTokenQuery is the SQL query to delete a refresh token
	DeleteRefreshTokenQuery = `DELETE FROM refresh_tokens WHERE jti = ?`

	// CheckRefreshTokenQuery is the SQL query to check if a refresh token exists
	CheckRefreshTokenQuery = `SELECT COUNT(1) FROM refresh_tokens WHERE jti = ?`

	// CreateAccessTokensTableQuery is the SQL query to create the access_tokens table
	CreateAccessTokensTableQuery = `CREATE TABLE IF NOT EXISTS access_tokens (jti TEXT PRIMARY KEY)`

	// InsertAccessTokenQuery is the SQL query to insert a new access token
	InsertAccessTokenQuery = `INSERT OR IGNORE INTO access_tokens (jti) VALUES (?)`

	// DeleteAccessTokenQuery is the SQL query to delete an access token
	DeleteAccessTokenQuery = `DELETE FROM access_tokens WHERE jti = ?`

	// CheckAccessTokenQuery is the SQL query to check if an access token exists
	CheckAccessTokenQuery = `SELECT COUNT(1) FROM access_tokens WHERE jti = ?`
)
