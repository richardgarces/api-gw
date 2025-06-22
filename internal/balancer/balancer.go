package balancer

import (
	"math/rand"
	"sync"
)

type RoundRobin struct {
	targets []string
	index   int
	mu      sync.Mutex
}

func NewRoundRobin(targets []string) *RoundRobin {
	return &RoundRobin{targets: targets}
}

func (rr *RoundRobin) Next() string {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	if len(rr.targets) == 0 {
		return ""
	}
	target := rr.targets[rr.index%len(rr.targets)]
	rr.index = (rr.index + 1) % len(rr.targets)
	return target
}

func PickTarget(targets []string) string {
	if len(targets) == 0 {
		return ""
	}
	return targets[rand.Intn(len(targets))]
}
