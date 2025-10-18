package sync

import (
	"context"
	"database/sql"
	"log/slog"
	"sync"
	"time"

	godatabases "github.com/ralvarezdev/go-databases"
	godatabasessql "github.com/ralvarezdev/go-databases/sql"
)

type (
	// DefaultService is the default implementation of the Service interface
	DefaultService struct {
		godatabasessql.Handler
		logger *slog.Logger
		mutex  sync.Mutex
	}
)

// NewDefaultService creates a new DefaultService
//
// Parameters:
//
//   - handler: the SQL connection handler
//   - logger: the logger (optional, can be nil)
//
// Returns:
//
//   - *DefaultService: the DefaultService instance
//   - error: an error if the data source or driver name is empty
func NewDefaultService(
	handler godatabasessql.Handler,
	logger *slog.Logger,
) (*DefaultService, error) {
	// Check if the handler is nil
	if handler == nil {
		return nil, godatabases.ErrNilHandler
	}

	if logger != nil {
		logger = logger.With(
			slog.String("component", "sync_sql_service"),
		)
	}

	return &DefaultService{
		Handler: handler,
		logger:  logger,
	}, nil
}

// Connect opens the database connection
//
// Returns:
//
//   - error: an error if the connection could not be opened
func (d *DefaultService) Connect() error {
	// Check if the service is nil
	if d == nil {
		return godatabases.ErrNilService
	}

	// Lock the mutex to ensure thread safety
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// Connect to the database
	db, err := d.Handler.Connect()
	if err != nil {
		if d.logger != nil {
			d.logger.Error(
				"Failed to connect to database",
				slog.String("error", err.Error()),
			)
		}
		return err
	}

	// Ensure the tables exist
	if _, err = db.Exec(CreateSyncTokensTableQuery); err != nil {
		return err
	}

	return nil
}

// Disconnect closes the database connection
//
// Returns:
//
//   - error: an error if the connection could not be closed
func (d *DefaultService) Disconnect() error {
	// Check if the service is nil
	if d == nil {
		return godatabases.ErrNilService
	}

	// Lock the mutex to ensure thread safety
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// Disconnect from the database
	if err := d.Handler.Disconnect(); err != nil {
		if d.logger != nil {
			d.logger.Error(
				"Failed to disconnect from database",
				slog.String("error", err.Error()),
			)
		}
		return err
	}

	return nil
}

// DB is a helper function to get the database connection
//
// Returns:
//
//   - *sql.DB: the database connection
func (d *DefaultService) DB() (*sql.DB, error) {
	// Lock the mutex to ensure thread safety
	d.mutex.Lock()

	// Get the database connection
	db, err := d.DB()
	if err != nil {
		d.mutex.Unlock()
		if d.logger != nil {
			d.logger.Error(
				"Failed to get database connection",
				slog.String("error", err.Error()),
			)
		}
		return nil, err
	}
	d.mutex.Unlock()

	return db, nil
}

// UpdateLastSyncTokensUpdateAt updates the last sync tokens updated at timestamp
//
// Parameters:
//
//   - ctx: the context
//   - updatedAt: the new timestamp
//
// Returns:
//
//   - error: an error if the timestamp could not be updated
func (d *DefaultService) UpdateLastSyncTokensUpdateAt(
	ctx context.Context,
	updatedAt time.Time,
) error {
	// Check if the service is nil
	if d == nil {
		return godatabases.ErrNilService
	}

	// Get the database connection
	db, err := d.DB()
	if err != nil {
		return err
	}

	// Update the last sync tokens updated at timestamp
	_, err = db.QueryContext(
		ctx,
		UpdateLastSyncTokensUpdatedAtQuery,
		updatedAt.Unix()-1, // Subtract 1 second to avoid ignoring updates within the same second
	)
	if err != nil {
		if d.logger != nil {
			d.logger.Error(
				"Failed to update last sync tokens updated at",
				slog.String("error", err.Error()),
			)
		}
		return err
	}

	return nil
}

// GetLastSyncTokensUpdatedAt gets the last sync tokens updated at timestamp
//
// Parameters:
//
//   - ctx: the context
//
// Returns:
//
//   - time.Time: the last sync tokens updated at timestamp
//   - error: an error if the timestamp could not be retrieved
func (d *DefaultService) GetLastSyncTokensUpdatedAt(ctx context.Context) (
	time.Time,
	error,
) {
	// Check if the service is nil
	if d == nil {
		return time.Time{}, godatabases.ErrNilService
	}

	// Get the database connection
	db, err := d.DB()
	if err != nil {
		return time.Time{}, err
	}

	var updatedAt sql.NullInt64

	// Get the last sync tokens updated at timestamp
	row := db.QueryRowContext(ctx, GetLastSyncTokensUpdatedAtQuery)
	if err = row.Scan(&updatedAt); err != nil {
		if d.logger != nil {
			d.logger.Error(
				"Failed to get last sync tokens updated at",
				slog.String("error", err.Error()),
			)
		}
		return time.Time{}, err
	}

	// If the timestamp is null, it means it has never been set, so return the zero value
	if !updatedAt.Valid {
		return time.Time{}, nil
	}

	return time.Unix(updatedAt.Int64, 0), nil
}
