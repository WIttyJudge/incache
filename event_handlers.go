package incache

import (
	"sync"
)

func defaultInsertionEvent(key string, value interface{}) {}
func defaultEvictionEvent(key string, value interface{})  {}

type eventHandlers struct {
	wg          *sync.WaitGroup
	onInsertion func(key string, value interface{})
	onEviction  func(key string, value interface{})
}

func newEventHandlers() *eventHandlers {
	return &eventHandlers{
		wg:          &sync.WaitGroup{},
		onInsertion: defaultInsertionEvent,
		onEviction:  defaultEvictionEvent,
	}
}

func (c *eventHandlers) OnInsertion(fn func(key string, value interface{})) {
	c.onInsertion = func(key string, value interface{}) {
		c.wg.Add(1)

		go func() {
			fn(key, value)
			c.wg.Done()
		}()
	}
}

func (c *eventHandlers) OnEviction(fn func(key string, value interface{})) {
	c.onEviction = func(key string, value interface{}) {
		c.wg.Add(1)

		go func() {
			fn(key, value)
			c.wg.Done()
		}()
	}
}

func (c *eventHandlers) Wait() {
	c.wg.Wait()
}
