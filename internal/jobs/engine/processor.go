// internal/jobs/engine/processor.go
package engine

import (
	"context"
	"fmt"
	"time"

	"tanmore_backend/internal/db/sqlc"
	"tanmore_backend/internal/jobs/engine/router"
	"tanmore_backend/pkg/sqlnull"
)

type ProcessorEngine struct {
	q *sqlc.Queries
}

func NewProcessorEngine(q *sqlc.Queries) *ProcessorEngine {
	return &ProcessorEngine{q: q}
}

func (e *ProcessorEngine) Start(ctx context.Context) error {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-ticker.C:

			// fmt.Println("⏱️ Background job ticked")

			events, err := e.q.FetchUndispatchedEvents(ctx, 50)
			if err != nil {
				continue
			}

			for _, event := range events {
				fmt.Println("➡️ Routing event:", event.ID, event.EventType)

				err := router.Route(ctx, event, e.q)
				if err != nil {
					fmt.Println("❌ Failed to route event:", err)
					continue
				}

				fmt.Println("✅ Marking event dispatched:", event.ID)

				_ = e.q.MarkEventDispatched(ctx, sqlc.MarkEventDispatchedParams{
					ID:           event.ID,
					DispatchedAt: sqlnull.Time(time.Now().UTC()),
				})
			}
		}
	}
}
