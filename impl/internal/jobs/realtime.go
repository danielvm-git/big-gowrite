package jobs

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RealtimeChannel is the Postgres NOTIFY channel used for realtime fanout
// (replaces Redis pub/sub). Contexts attach the notify_realtime() trigger to
// their tables; subscribers LISTEN on this channel.
const RealtimeChannel = "gowrite_realtime"

// Notification is one realtime event delivered over LISTEN/NOTIFY.
type Notification struct {
	Payload string
}

// Listen subscribes to the realtime channel and delivers notifications to the
// returned channel until ctx is canceled. It acquires a dedicated connection
// from the pool for the lifetime of the subscription.
func Listen(ctx context.Context, pool *pgxpool.Pool) (<-chan Notification, error) {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("jobs: acquire listen conn: %w", err)
	}
	if _, err := conn.Exec(ctx, "LISTEN "+RealtimeChannel); err != nil {
		conn.Release()
		return nil, fmt.Errorf("jobs: LISTEN: %w", err)
	}

	out := make(chan Notification)
	go func() {
		defer conn.Release()
		defer close(out)
		for {
			n, err := conn.Conn().WaitForNotification(ctx)
			if err != nil {
				return // ctx canceled or conn closed
			}
			select {
			case out <- Notification{Payload: n.Payload}:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out, nil
}

// Notify publishes a raw payload on the realtime channel (for contexts that
// emit events explicitly rather than via a table trigger).
func (q *Queue) Notify(ctx context.Context, payload string) error {
	_, err := q.pool.Exec(ctx, "SELECT pg_notify($1, $2)", RealtimeChannel, payload)
	if err != nil {
		return fmt.Errorf("jobs: notify: %w", err)
	}
	return nil
}
