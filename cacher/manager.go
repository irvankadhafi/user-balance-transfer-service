package cacher

import (
	"github.com/go-redsync/redsync/v4"
	redigosync "github.com/go-redsync/redsync/v4/redis/redigo"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/jpillora/backoff"
	"time"
)

const (
	// Override these when constructing the cache keeper
	defaultTTL          = 10 * time.Second
	defaultNilTTL       = 5 * time.Minute
	defaultLockDuration = 1 * time.Minute
	defaultLockTries    = 1
	defaultWaitTime     = 15 * time.Second
)

var nilValue = []byte("null")

type (
	CacheManager interface {
		Get(key string) (any, error)
		GetOrLock(key string) (any, *redsync.Mutex, error)
		StoreWithoutBlocking(Item) error
		StoreMultiWithoutBlocking([]Item) error
		DeleteByKeys([]string) error
		StoreNil(cacheKey string) error
		IncreaseCachedValueByOne(key string) error
		Expire(string, time.Duration) error

		GetTTL(string) (int64, error)

		AcquireLock(string) (*redsync.Mutex, error)
		SetDefaultTTL(time.Duration)
		SetNilTTL(time.Duration)
		SetConnectionPool(*redigo.Pool)
		SetLockConnectionPool(*redigo.Pool)
		SetLockDuration(time.Duration)
		SetLockTries(int)
		SetWaitTime(time.Duration)
		SetDisableCaching(bool)
	}
)

type cacheManager struct {
	connPool       *redigo.Pool
	nilTTL         time.Duration
	defaultTTL     time.Duration
	waitTime       time.Duration
	disableCaching bool

	lockConnPool *redigo.Pool
	lockDuration time.Duration
	lockTries    int
}

// NewCacheManager :nodoc:
func NewCacheManager() CacheManager {
	return &cacheManager{
		defaultTTL:     defaultTTL,
		nilTTL:         defaultNilTTL,
		lockDuration:   defaultLockDuration,
		lockTries:      defaultLockTries,
		waitTime:       defaultWaitTime,
		disableCaching: false,
	}
}

func (k *cacheManager) Get(key string) (cachedItem any, err error) {
	if k.disableCaching {
		return
	}

	cachedItem, err = get(k.connPool.Get(), key)
	if err != nil && err != ErrKeyNotExist && err != redigo.ErrNil || cachedItem != nil {
		return
	}

	return nil, nil
}

// GetOrLock :nodoc:
func (k *cacheManager) GetOrLock(key string) (cachedItem any, mutex *redsync.Mutex, err error) {
	if k.disableCaching {
		return
	}

	cachedItem, err = get(k.connPool.Get(), key)
	if err != nil && err != ErrKeyNotExist && err != redigo.ErrNil || cachedItem != nil {
		return
	}

	mutex, err = k.AcquireLock(key)
	if err == nil {
		return
	}

	start := time.Now()
	for {
		b := &backoff.Backoff{
			Min:    20 * time.Millisecond,
			Max:    200 * time.Millisecond,
			Jitter: true,
		}

		if !k.isLocked(key) {
			cachedItem, err = get(k.connPool.Get(), key)
			if err != nil {
				if err == ErrKeyNotExist {
					mutex, err = k.AcquireLock(key)
					if err == nil {
						return nil, mutex, nil
					}

					goto Wait
				}
				return nil, nil, err
			}
			return cachedItem, nil, nil
		}

	Wait:
		elapsed := time.Since(start)
		if elapsed >= k.waitTime {
			break
		}

		time.Sleep(b.Duration())
	}

	return nil, nil, ErrWaitTooLong
}

// IncreaseCachedValueByOne will increments the number stored at key by one.
// If the key does not exist, it is set to 0 before performing the operation
func (k *cacheManager) IncreaseCachedValueByOne(key string) error {
	if k.disableCaching {
		return nil
	}

	client := k.connPool.Get()
	defer func() {
		_ = client.Close()
	}()

	_, err := client.Do("INCR", key)
	return err
}

// Expire Set expire a key
func (k *cacheManager) Expire(key string, duration time.Duration) (err error) {
	if k.disableCaching {
		return nil
	}

	client := k.connPool.Get()
	defer func() {
		_ = client.Close()
	}()

	_, err = client.Do("EXPIRE", key, int64(duration.Seconds()))
	return
}

func (k *cacheManager) GetTTL(name string) (value int64, err error) {
	client := k.connPool.Get()
	defer func() {
		_ = client.Close()
	}()

	val, err := client.Do("TTL", name)
	if err != nil {
		return
	}

	value = val.(int64)
	return
}

func (k *cacheManager) StoreWithoutBlocking(c Item) error {
	if k.disableCaching {
		return nil
	}

	client := k.connPool.Get()
	defer func() {
		_ = client.Close()
	}()

	_, err := client.Do("SETEX", c.GetKey(), k.decideCacheTTL(c), c.GetValue())
	return err
}

// StoreMultiWithoutBlocking Store multiple items
func (k *cacheManager) StoreMultiWithoutBlocking(items []Item) error {
	if k.disableCaching {
		return nil
	}

	client := k.connPool.Get()
	defer func() {
		_ = client.Close()
	}()

	err := client.Send("MULTI")
	if err != nil {
		return err
	}
	for _, item := range items {
		err = client.Send("SETEX", item.GetKey(), k.decideCacheTTL(item), item.GetValue())
		if err != nil {
			return err
		}
	}

	_, err = client.Do("EXEC")
	return err
}

// DeleteByKeys Delete by multiple keys
func (k *cacheManager) DeleteByKeys(keys []string) error {
	if k.disableCaching {
		return nil
	}

	if len(keys) <= 0 {
		return nil
	}

	client := k.connPool.Get()
	defer func() {
		_ = client.Close()
	}()

	var redisKeys []any
	for _, key := range keys {
		redisKeys = append(redisKeys, key)
	}

	_, err := client.Do("DEL", redisKeys...)
	return err
}

// AcquireLock :nodoc:
func (k *cacheManager) AcquireLock(key string) (*redsync.Mutex, error) {
	p := redigosync.NewPool(k.lockConnPool)
	r := redsync.New(p)
	m := r.NewMutex("lock:"+key,
		redsync.WithExpiry(k.lockDuration),
		redsync.WithTries(k.lockTries))

	return m, m.Lock()
}

// SetDefaultTTL :nodoc:
func (k *cacheManager) SetDefaultTTL(d time.Duration) {
	k.defaultTTL = d
}

func (k *cacheManager) SetNilTTL(d time.Duration) {
	k.nilTTL = d
}

// SetConnectionPool :nodoc:
func (k *cacheManager) SetConnectionPool(c *redigo.Pool) {
	k.connPool = c
}

// SetLockConnectionPool :nodoc:
func (k *cacheManager) SetLockConnectionPool(c *redigo.Pool) {
	k.lockConnPool = c
}

// SetLockDuration :nodoc:
func (k *cacheManager) SetLockDuration(d time.Duration) {
	k.lockDuration = d
}

// SetLockTries :nodoc:
func (k *cacheManager) SetLockTries(t int) {
	k.lockTries = t
}

// SetWaitTime :nodoc:
func (k *cacheManager) SetWaitTime(d time.Duration) {
	k.waitTime = d
}

// SetDisableCaching :nodoc:
func (k *cacheManager) SetDisableCaching(b bool) {
	k.disableCaching = b
}

// StoreNil :nodoc:
func (k *cacheManager) StoreNil(cacheKey string) error {
	item := NewItemWithCustomTTL(cacheKey, nilValue, k.nilTTL)
	err := k.StoreWithoutBlocking(item)
	return err
}

func (k *cacheManager) decideCacheTTL(c Item) (ttl int64) {
	if ttl = c.GetTTLInt64(); ttl > 0 {
		return
	}

	return int64(k.defaultTTL.Seconds())
}

func (k *cacheManager) isLocked(key string) bool {
	client := k.lockConnPool.Get()
	defer func() {
		_ = client.Close()
	}()

	reply, err := client.Do("GET", "lock:"+key)
	if err != nil || reply == nil {
		return false
	}

	return true
}
