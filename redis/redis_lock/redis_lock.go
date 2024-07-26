package redislock

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"sync"
	"time"
)

var defaultExpiration = 10 * time.Second

// Lock is a non-reentrant distributed lock implemented by Redis
type Lock struct {
	m            sync.Mutex
	client       *redis.Client
	key          string
	holder       string
	expiration   time.Duration
	keepAlive    bool
	keepAliveTtl time.Duration
}

func NewLock(client *redis.Client, key string, expiration time.Duration) *Lock {
	// Generate a unique identifier for the process
	holder := uuid.New().String()
	if expiration == 0 {
		expiration = defaultExpiration
	}
	return &Lock{
		client:       client,
		key:          key,
		holder:       holder,
		expiration:   expiration,
		keepAlive:    false,
		keepAliveTtl: expiration / 2, // Renew the key before half the expiration time
	}
}

func (r *Lock) Lock() (bool, error) {
	r.m.Lock()
	defer r.m.Unlock()

	success, err := r.client.SetNX(context.Background(), r.key, r.holder, r.expiration).Result()
	if err != nil {
		return false, err
	}

	if success && !r.keepAlive {
		r.startKeepAlive()
	}

	return success, nil
}

func (r *Lock) Unlock() error {
	r.m.Lock()
	defer r.m.Unlock()

	unlockScript := `
    local currentHolder = redis.call("GET", KEYS[1])
    if currentHolder == ARGV[1] then
		return redis.call("DEL", KEYS[1])
    else
        return -1
    end
    `

	unlockedStatus, err := r.client.Eval(context.Background(), unlockScript, []string{r.key}, r.holder).Int()
	if err != nil {
		return fmt.Errorf("unlock failed: %w", err)
	}

	switch unlockedStatus {
	case 0:
		return fmt.Errorf("unlock failed: lock does not exist")
	case 1:
		// Stop keep alive mechanism
		r.stopKeepAlive()
	case -1:
		return fmt.Errorf("unlock failed: lock is held by a different holder")
	}

	return nil
}

// startKeepAlive starts the keep-alive mechanism for the lock.
func (r *Lock) startKeepAlive() {
	r.keepAlive = true
	go func() {
		ticker := time.NewTicker(r.keepAliveTtl)
		defer ticker.Stop()

		for r.keepAlive {
			select {
			case <-ticker.C:
				r.client.PExpire(context.Background(), r.key, r.expiration)
			}
		}
	}()
}

// stopKeepAlive stops the keep-alive mechanism for the lock.
func (r *Lock) stopKeepAlive() {
	r.keepAlive = false
}
