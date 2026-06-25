package main

import (
	"sync"
)

type BotState struct {
	sync.Mutex
	userChannels  map[int64]string
	isInitialized map[int64]bool
	msgQueue      map[int64][]string
}
