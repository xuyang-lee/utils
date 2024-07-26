package reentrantlock

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"sync"
	"time"
)

var defaultExpiration = 10 * time.Second

// ReentrantLock is a reentrant distributed lock implemented by Redis
type ReentrantLock struct {
	client       *redis.Client
	key          string
	holder       string
	expiration   time.Duration
	keepAlive    bool
	mutex        sync.Mutex
	keepAliveTtl time.Duration
}

func NewReentrantLock(client *redis.Client, key string, expiration time.Duration) *ReentrantLock {
	// Generate a unique identifier for the process
	holder := uuid.New().String()
	if expiration == 0 {
		expiration = defaultExpiration
	}
	return &ReentrantLock{
		client:       client,
		key:          key,
		holder:       holder,
		expiration:   expiration,
		keepAlive:    false,
		keepAliveTtl: expiration / 2, // Renew the key before half the expiration time
	}
}

func (r *ReentrantLock) Lock() (bool, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	lockScript := `
    local currentHolder = redis.call("HGET", KEYS[1], "holder")
    if currentHolder == false then
        -- Lock is not held, acquire it
        redis.call("HSET", KEYS[1], "holder", ARGV[1], "count", 1)
        redis.call("PEXPIRE", KEYS[1], ARGV[2])
        return 1
    elseif currentHolder == ARGV[1] then
        -- Lock is held by current process, reentrant
        local count = redis.call("HINCRBY", KEYS[1], "count", 1)
        redis.call("PEXPIRE", KEYS[1], ARGV[2])
        return 1
    else
        return 0
    end
    `

	acquired, err := r.client.Eval(context.Background(), lockScript, []string{r.key}, r.holder, int(r.expiration/time.Millisecond)).Bool()
	if err != nil {
		return false, err
	}

	if acquired && !r.keepAlive {
		// Start keep alive mechanism
		r.startKeepAlive()
	}

	return acquired, nil
}

func (r *ReentrantLock) Unlock() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	unlockScript := `
    local currentHolder = redis.call("HGET", KEYS[1], "holder")
    if currentHolder == ARGV[1] then
        local count = redis.call("HINCRBY", KEYS[1], "count", -1)
        if count <= 0 then
            redis.call("DEL", KEYS[1])
			return 0
        end
        return 1
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

		// Stop keep alive mechanism
		r.stopKeepAlive()
	case 1:
		return nil
	case -1:
		return fmt.Errorf("unlock failed: lock is held by a different holder")
	}

	return nil
}

// startKeepAlive starts the keep-alive mechanism for the lock.
func (r *ReentrantLock) startKeepAlive() {
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
func (r *ReentrantLock) stopKeepAlive() {
	r.keepAlive = false
}
