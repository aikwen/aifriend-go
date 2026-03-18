package task

import (
	"log"
	"runtime/debug"
)

func Go(fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("async task panic: %v\n%s", r, debug.Stack())
			}
		}()

		fn()
	}()
}