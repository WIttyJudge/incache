package cache

import "time"

const (
	ItemNoTTL      time.Duration = -1
	ItemDefaultTTL time.Duration = 0
)

type Item struct {
	Value     interface{}
	TTL       time.Duration
	ExpiresAt time.Time
}

func newItem(value interface{}, ttl time.Duration) Item {
	expiresAt := time.Now().Add(ttl)

	item := Item{
		Value:     value,
		TTL:       ttl,
		ExpiresAt: expiresAt,
	}

	return item
}

func (i Item) Expired() bool {
	if i.TTL <= 0 {
		return false
	}

	return time.Now().After(i.ExpiresAt)
}
