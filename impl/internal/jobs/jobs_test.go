package jobs

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// testPool connects to GOWRITE_TEST_DATABASE_URL, skipping the test if unset.
func testPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	dsn := os.Getenv("GOWRITE_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("GOWRITE_TEST_DATABASE_URL not set; skipping Postgres integration test")
	}
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	t.Cleanup(pool.Close)
	// Clean slate for each test.
	if _, err := pool.Exec(context.Background(), "TRUNCATE jobs RESTART IDENTITY"); err != nil {
		t.Fatalf("truncate: %v", err)
	}
	return pool
}

func TestEnqueueAndClaim(t *testing.T) {
	pool := testPool(t)
	q := New(pool)
	ctx := context.Background()

	id, inserted, err := q.Enqueue(ctx, "webhook", map[string]any{"url": "https://x"}, EnqueueOptions{})
	if err != nil || !inserted || id == 0 {
		t.Fatalf("enqueue: id=%d inserted=%v err=%v", id, inserted, err)
	}

	j, err := q.Claim(ctx)
	if err != nil || j == nil {
		t.Fatalf("claim: %v (job=%v)", err, j)
	}
	if j.ID != id || j.Kind != "webhook" || j.Attempts != 1 {
		t.Fatalf("claimed job = %+v", j)
	}
	if s, _ := q.Status(ctx, id); s != StatusProcessing {
		t.Fatalf("status after claim = %s, want processing", s)
	}

	if err := q.Complete(ctx, id); err != nil {
		t.Fatalf("complete: %v", err)
	}
	if s, _ := q.Status(ctx, id); s != StatusCompleted {
		t.Fatalf("status after complete = %s, want completed", s)
	}
}

func TestDedupIdempotency(t *testing.T) {
	pool := testPool(t)
	q := New(pool)
	ctx := context.Background()

	id1, ins1, err := q.Enqueue(ctx, "build", map[string]any{"cloudEventId": "evt-1"}, EnqueueOptions{DedupKey: "evt-1"})
	if err != nil || !ins1 {
		t.Fatalf("first enqueue: id=%d inserted=%v err=%v", id1, ins1, err)
	}
	id2, ins2, err := q.Enqueue(ctx, "build", map[string]any{"cloudEventId": "evt-1"}, EnqueueOptions{DedupKey: "evt-1"})
	if err != nil {
		t.Fatalf("second enqueue: %v", err)
	}
	if ins2 {
		t.Fatal("duplicate dedup_key should NOT insert a new row")
	}
	if id1 != id2 {
		t.Fatalf("dedup should return the same id: %d vs %d", id1, id2)
	}

	var count int
	_ = pool.QueryRow(ctx, "SELECT count(*) FROM jobs WHERE dedup_key = 'evt-1'").Scan(&count)
	if count != 1 {
		t.Fatalf("expected exactly 1 row for dedup_key, got %d", count)
	}
}

// TestSkipLockedExactlyOnce runs many concurrent claimers against N jobs and
// asserts each job is claimed by exactly one worker (no double-claim).
func TestSkipLocked(t *testing.T) {
	pool := testPool(t)
	q := New(pool)
	ctx := context.Background()

	const n = 20
	for i := 0; i < n; i++ {
		if _, _, err := q.Enqueue(ctx, "work", map[string]any{"i": i}, EnqueueOptions{}); err != nil {
			t.Fatal(err)
		}
	}

	var mu sync.Mutex
	claimed := map[int64]int{}
	var wg sync.WaitGroup
	for w := 0; w < 8; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				j, err := q.Claim(ctx)
				if err != nil || j == nil {
					return
				}
				mu.Lock()
				claimed[j.ID]++
				mu.Unlock()
				_ = q.Complete(ctx, j.ID)
			}
		}()
	}
	wg.Wait()

	if len(claimed) != n {
		t.Fatalf("claimed %d distinct jobs, want %d", len(claimed), n)
	}
	for id, c := range claimed {
		if c != 1 {
			t.Fatalf("job %d claimed %d times (want exactly once)", id, c)
		}
	}
}

func TestNotify(t *testing.T) {
	pool := testPool(t)
	q := New(pool)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ch, err := Listen(ctx, pool)
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	// Give the LISTEN a moment to register, then notify.
	time.Sleep(100 * time.Millisecond)
	if err := q.Notify(ctx, "hello-realtime"); err != nil {
		t.Fatalf("notify: %v", err)
	}

	select {
	case n := <-ch:
		if n.Payload != "hello-realtime" {
			t.Fatalf("payload = %q", n.Payload)
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for notification")
	}
}
