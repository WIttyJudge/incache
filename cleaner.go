package incache

import "time"

type expiredDeleter interface {
	DeleteExpired()
}

// The structure is supposed to control an automatic cleanup background
// process that calls DeleteExpired() method every time specified
// in cleanupInterval variable.
type cleaner struct {
	cleanupInterval time.Duration

	closeCh chan struct{}
}

func newCleaner(cleanupInterval time.Duration) *cleaner {
	return &cleaner{
		cleanupInterval: cleanupInterval,

		closeCh: make(chan struct{}),
	}
}

func (c *cleaner) start(ed expiredDeleter) {
	go func() {
		ticker := time.NewTicker(c.cleanupInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				ed.DeleteExpired()
			case <-c.closeCh:
				return
			}
		}
	}()
}

func (c *cleaner) close() {
	c.closeCh <- struct{}{}
}
