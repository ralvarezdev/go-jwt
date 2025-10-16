package sync

import (
	"context"
	"time"
)

type (
	// Service is the interface for the SQLite service to store the last sync tokens updated at
	Service interface {
		UpdateLastSyncTokensUpdateAt(
			ctx context.Context,
			updatedAt time.Time,
		) error
		GetLastSyncTokensUpdatedAt(ctx context.Context) (time.Time, error)
	}
)
