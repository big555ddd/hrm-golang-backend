package notification

import (
	"app/internal/logger"
	"sync"
)

var (
	hubInstance *Hub
	once        sync.Once
)

// GetHub returns the singleton Hub instance
func GetHub() *Hub {
	once.Do(func() {
		hubInstance = NewHub()
		go hubInstance.Run() // Start the hub in a goroutine
		logger.Infof("Singleton Hub created and started with pointer: %p", hubInstance)
	})
	return hubInstance
}
