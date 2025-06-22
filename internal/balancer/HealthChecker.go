package balancer

import (
    "net/http"
    "sync"
    "time"
)

type HealthChecker struct {
    targets   []string
    healthy   map[string]bool
    mu        sync.RWMutex
    interval  time.Duration
}

func NewHealthChecker(targets []string, interval time.Duration) *HealthChecker {
    hc := &HealthChecker{
        targets:  targets,
        healthy:  make(map[string]bool),
        interval: interval,
    }
    go hc.run()
    return hc
}

func (hc *HealthChecker) run() {
    for {
        for _, t := range hc.targets {
            go hc.check(t)
        }
        time.Sleep(hc.interval)
    }
}

func (hc *HealthChecker) check(target string) {
    resp, err := http.Get(target + "/health")
    hc.mu.Lock()
    defer hc.mu.Unlock()
    if err == nil && resp.StatusCode == 200 {
        hc.healthy[target] = true
    } else {
        hc.healthy[target] = false
    }
}

func (hc *HealthChecker) GetHealthyTargets() []string {
    hc.mu.RLock()
    defer hc.mu.RUnlock()
    var healthy []string
    for _, t := range hc.targets {
        if hc.healthy[t] {
            healthy = append(healthy, t)
        }
    }
    return healthy
}