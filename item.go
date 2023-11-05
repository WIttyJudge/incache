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

	item.setExpiration()

	return item
}

// Expired checks whether the item has expired.
func (i Item) Expired() bool {
	if !i.CanExpire() {
		return false
	}

	return time.Now().After(i.ExpiresAt)
}

// CanExpire checks whether the item can expire.
func (i Item) CanExpire() bool {
	return !i.ExpiresAt.IsZero()
}

func (i *Item) setExpiration() {
	if i.TTL <= 0 {
		return
	}

	i.ExpiresAt = time.Now().Add(i.TTL)
}
