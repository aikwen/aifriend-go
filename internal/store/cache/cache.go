package cache

type Cache struct {
	SystemPrompt *SystemPromptCache
}

func New() *Cache {
	return &Cache{
		SystemPrompt: NewSystemPromptCache(),
	}
}