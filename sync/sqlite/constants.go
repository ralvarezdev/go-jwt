package sync

const (
	// CreateSyncTokensTableQuery is the SQL query to create the sync_tokens table
	CreateSyncTokensTableQuery = `
CREATE TABLE IF NOT EXISTS sync_tokens (id INTEGER PRIMARY KEY AUTOINCREMENT, updated_at DATETIME NOT NULL);
`

	// UpdateLastSyncTokensUpdatedAtQuery is the SQL query to insert a new sync tokens record
	UpdateLastSyncTokensUpdatedAtQuery = `
INSERT INTO sync_tokens (updated_at) VALUES (?);
`

	// GetLastSyncTokensUpdatedAtQuery is the SQL query to get the last sync tokens record
	GetLastSyncTokensUpdatedAtQuery = `
SELECT updated_at FROM sync_tokens ORDER BY updated_at DESC LIMIT 1;
`
)
