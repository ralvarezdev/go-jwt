package sql

const (
	// CreateTableQuery is the SQL query to create the tokens table
	CreateTableQuery = `CREATE TABLE IF NOT EXISTS tokens (jti TEXT PRIMARY KEY)`

	// InsertTokenQuery is the SQL query to insert a new token
	InsertTokenQuery = `INSERT OR IGNORE INTO tokens (jti) VALUES (?)`

	// DeleteTokenQuery is the SQL query to delete a token
	DeleteTokenQuery = `DELETE FROM tokens WHERE jti = ?`

	// CheckTokenQuery is the SQL query to check if a token exists
	CheckTokenQuery = `SELECT COUNT(1) FROM tokens WHERE jti = ?`
)
