package main

import (
	"time"
)
func (rl *rateLimiter) allow(ip string, limit int, window time.Duration) bool {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    // cleanup() runs in background gorountines and allow() runs on every request, the mutex make them take turns before data race occurs

    now := time.Now()
    v, exists := rl.visitors[ip]
    
    if !exists || now.After(v.resetAt) { // i.e if never seen this window before or time has passed
        rl.visitors[ip] = &visitor{
            count:   1,
            resetAt: now.Add(window),
        }
        return true
    }
    
    // within window
    if v.count >= limit {
        return false
    }
    
    v.count++
    return true
}


func (rl *rateLimiter) cleanup() {
    ticker := time.NewTicker(1 * time.Minute) // Creates an alarm that goes off every 1 minute, Each time it rings, it sends a signal on a channel called ticker.C.

    for range ticker.C { // waits for the alarm to ring then runs the code below then waits again and then runs again
        rl.mu.Lock()
        now := time.Now()
        for ip, v := range rl.visitors { // locks during cleanup and graps the ip and value (visitors record) during loop
            if now.After(v.resetAt) { // is the current time later than when this visitors window was supposed to reset
                delete(rl.visitors, ip)
            }
        }
        rl.mu.Unlock()
    }
}