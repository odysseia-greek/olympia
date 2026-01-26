package monos

import (
	"context"
	"fmt"
	"time"

	"github.com/odysseia-greek/agora/plato/logging"
)

// WaitForDictionarySettled polls Elasticsearch until the dictionary doc count:
//   - is >= minDocs
//   - has not changed for stableFor duration
//
// It stops waiting after maxWait and returns false.
//
// pollEvery should be something like 10s.
func (m *MelissosHandler) WaitForDictionarySettled(
	ctx context.Context,
	minDocs int64,
	pollEvery time.Duration,
	stableFor time.Duration,
	maxWait time.Duration,
) bool {
	start := time.Now()
	ticker := time.NewTicker(pollEvery)
	defer ticker.Stop()

	var (
		lastCount      int64 = -1
		lastChangeTime       = time.Now()
	)

	for {
		if time.Since(start) >= maxWait {
			logging.Info("Timed out waiting for dictionary to settle")
			return false
		}

		select {
		case <-ctx.Done():
			return false

		case <-ticker.C:
			reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			resp, err := m.Elastic.Query().CountRaw(reqCtx, m.Index, nil)
			cancel()

			if err != nil {
				logging.Debug(fmt.Sprintf("Dictionary count query failed: %v", err))
				continue
			}
			count := resp.Count

			// Sanity floor
			if count < minDocs {
				logging.Debug(fmt.Sprintf("Dictionary count %d below minDocs %d; waiting...", count, minDocs))
				continue
			}

			// Stability tracking
			if count != lastCount {
				logging.Info(fmt.Sprintf("Dictionary count changed: %d -> %d", lastCount, count))
				lastCount = count
				lastChangeTime = time.Now()
				continue
			}

			// unchanged this tick
			if time.Since(lastChangeTime) >= stableFor {
				logging.Info(fmt.Sprintf("Dictionary settled at %d docs for %s (minDocs=%d)", count, stableFor, minDocs))
				return true
			}

			logging.Debug(fmt.Sprintf(
				"Dictionary count stable at %d for %s/%s; waiting...",
				count, time.Since(lastChangeTime).Round(time.Second), stableFor,
			))
		}
	}
}
