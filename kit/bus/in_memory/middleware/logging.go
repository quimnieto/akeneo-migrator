package middleware

import (
	"context"
	"fmt"
	"time"

	"akeneo-migrator/kit/bus"
	"akeneo-migrator/kit/bus/in_memory"
)

// Logging creates a middleware that logs command execution
func Logging() inmemory.Middleware {
	return func(ctx context.Context, msg bus.Message, next inmemory.NextFunc) (bus.Response, error) {
		start := time.Now()
		
		fmt.Printf("🔄 Executing command: %s\n", msg.Type())
		
		response, err := next(ctx, msg)
		
		duration := time.Since(start)
		
		if err != nil {
			fmt.Printf("❌ Command failed: %s (took %v)\n", msg.Type(), duration)
		} else {
			fmt.Printf("✅ Command completed: %s (took %v)\n", msg.Type(), duration)
		}
		
		return response, err
	}
}
