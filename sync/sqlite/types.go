package sync

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	godatabases "github.com/ralvarezdev/go-databases"
	godatabasessql "github.com/ralvarezdev/go-databases/sql"
)

type (
	// DefaultService is the default implementation of the Service interface
	DefaultService struct {
		godatabasessql.Service
		logger *slog.Logger
	}
)

// NewDefaultService creates a new DefaultService
//
// Parameters:
//
//   - service: the SQL connection service
//   - logger: the logger (optional, can be nil)
//
// Returns:
//
//   - *DefaultService: the DefaultService instance
//   - error: an error if the data source or driver name is empty
func NewDefaultService(
	service godatabasessql.Service,
	logger *slog.Logger,
) (*DefaultService, error) {
	// Check if the service is nil
	if service == nil {
		return nil, godatabases.ErrNilService
	}

	if logger != nil {
		logger = logger.With(
			slog.String("component", "sync_sqlite_service"),
		)
	}

	return &DefaultService{
		Service: service,
		logger:  logger,
	}, nil
}

// Connect opens the database connection
//
// Parameters:
//
//   - ctx: the context
//
// Returns:
//
//   - error: an error if the connection could not be opened
func (d *DefaultService) Connect(ctx context.Context) error {
	// Check if the service is nil
	if d == nil {
		return godatabases.ErrNilService
	}

	// Connect to the database
	db, err := d.Service.Connect()
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
	if _, err = db.ExecContext(ctx, CreateSyncTokensTableQuery); err != nil {
		return err
	}

	return nil
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

	// Update the last sync tokens updated at timestamp
	if _, err := d.ExecWithCtx(
		ctx,
		&UpdateLastSyncTokensUpdatedAtQuery,
		updatedAt.Unix()-1, // Subtract 1 second to avoid ignoring updates within the same second
	); err != nil {
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

	// Get the last sync tokens updated at timestamp
	var updatedAt sql.NullInt64
	row, err := d.QueryRowWithCtx(ctx, &GetLastSyncTokensUpdatedAtQuery)
	if err != nil {
		if d.logger != nil {
			d.logger.Error(
				"Failed to query last sync tokens updated at",
				slog.String("error", err.Error()),
			)
		}
		return time.Time{}, err
	}
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
