package incache

import "time"

// The Item struct represents a cached item.
type Item struct {
	Value     interface{}
	TTL       time.Duration
	ExpiresAt time.Time
}

func newItem(value interface{}, ttl time.Duration) Item {
	item := Item{
		Value: value,
		TTL:   ttl,
	}

	if ttl > 0 {
		expiresAt := time.Now().Add(ttl)
		item.ExpiresAt = expiresAt
	}

	return item
}

// Checks whether the item has expired.
func (i Item) Expired() bool {
	if !i.CanExpire() {
		return false
	}

	return time.Now().After(i.ExpiresAt)
}

// Check whether the item can expire.
func (i Item) CanExpire() bool {
	return !i.ExpiresAt.IsZero()
}
