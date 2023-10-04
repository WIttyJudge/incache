package incache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewItem(t *testing.T) {
	item := newItem("value", 5*time.Second)

	require.NotNil(t, item)
	assert.Equal(t, "value", item.Value)
	assert.Equal(t, 5*time.Second, item.TTL)
	assert.WithinDuration(t, time.Now().Add(item.TTL), item.ExpiresAt, time.Second)
}

func TestNewItemWithoutTTL(t *testing.T) {
	item := newItem("value", 0)

	require.NotNil(t, item)
	assert.Equal(t, time.Duration(0), item.TTL)
	assert.True(t, item.ExpiresAt.IsZero())
}

func TestItemIsExpired(t *testing.T) {
	item := newItem("value", 1*time.Millisecond)
	time.Sleep(1 * time.Millisecond)

	assert.True(t, item.Expired())
	assert.True(t, item.CanExpire())
	assert.True(t, time.Now().After(item.ExpiresAt))
	assert.False(t, item.ExpiresAt.IsZero())
}

func TestItemIsNotExpired(t *testing.T) {
	item := newItem("value", 5*time.Second)

	assert.False(t, item.Expired())
	assert.True(t, item.CanExpire())
	assert.False(t, time.Now().After(item.ExpiresAt))
	assert.False(t, item.ExpiresAt.IsZero())
}

func TestItemNeverExpires(t *testing.T) {
	item := newItem("value", 0)

	assert.False(t, item.Expired())
	assert.False(t, item.CanExpire())
	assert.True(t, time.Now().After(item.ExpiresAt))
	assert.True(t, item.ExpiresAt.IsZero())
}
