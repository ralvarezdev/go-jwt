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
	// Service is the default implementation of the Service interface
	Service struct {
		godatabasessql.Service
		logger *slog.Logger
	}
)

// NewService creates a new Service
//
// Parameters:
//
//   - service: the SQL connection service
//   - logger: the logger (optional, can be nil)
//
// Returns:
//
//   - *Service: the Service instance
//   - error: an error if the data source or driver name is empty
func NewService(
	service godatabasessql.Service,
	logger *slog.Logger,
) (*Service, error) {
	// Check if the service is nil
	if service == nil {
		return nil, godatabases.ErrNilService
	}

	if logger != nil {
		logger = logger.With(
			slog.String("component", "sync_sqlite_service"),
		)
	}

	return &Service{
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
func (d *Service) Connect(ctx context.Context) error {
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
func (d *Service) UpdateLastSyncTokensUpdateAt(
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
func (d *Service) GetLastSyncTokensUpdatedAt(ctx context.Context) (
	time.Time,
	error,
) {
	// Check if the service is nil
	if d == nil {
		return time.Time{}, godatabases.ErrNilService
	}

	// Get the last sync tokens updated at timestamp
	var updatedAt sql.NullTime
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
	if rowErr := row.Scan(&updatedAt); rowErr != nil {
		if d.logger != nil {
			d.logger.Error(
				"Failed to get last sync tokens updated at",
				slog.String("error", rowErr.Error()),
			)
		}
		return time.Time{}, rowErr
	}

	// Return zero time if updatedAt is null
	if !updatedAt.Valid {
		return time.Time{}, nil
	}
	return updatedAt.Time, nil
}
