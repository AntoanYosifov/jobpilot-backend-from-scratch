package handler

import (
	"jobpilot/model"
	"sync"
)

var (
	jobs   = []model.Job{}
	nextID = 1
	mu     sync.Mutex
)
