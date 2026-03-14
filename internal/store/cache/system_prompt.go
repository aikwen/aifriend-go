package cache

import (
	"sync"

	"github.com/cloudwego/eino/schema"
)

type SystemPromptCache struct {
	mu   sync.RWMutex
	data map[string][]*schema.Message
}

func NewSystemPromptCache() *SystemPromptCache {
	return &SystemPromptCache{
		data: make(map[string][]*schema.Message),
	}
}

func (c *SystemPromptCache) Get(title string) ([]*schema.Message, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	msgs, ok := c.data[title]
	return msgs, ok
}

func (c *SystemPromptCache) Set(title string, msgs []*schema.Message) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[title] = msgs
}

func (c *SystemPromptCache) Delete(title string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, title)
}

func (c *SystemPromptCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = make(map[string][]*schema.Message)
}