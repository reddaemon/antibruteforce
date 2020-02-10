package bucket

import (
	"context"
	"errors"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Bucket struct {
	Capacity  uint
	Remaining uint
	Drop      time.Time
	Rate      time.Duration
	Mutex     sync.Mutex
}

type Storage interface {
	Add(ctx context.Context, key string, capacity uint, rate time.Duration) error
	Drop(ctx context.Context, keys []string) error
	CleanStorage() error
}

type MemRepo struct {
	buckets map[string]*Bucket
	mutex   sync.Mutex
	l       *zap.Logger
}

const cleanTimer = 10

func (r *MemRepo) initCleaner() {
	go func() {
		for {
			time.Sleep(time.Duration(cleanTimer) * time.Minute)
		}
	}()
}

func NewMemRepo(logger *zap.Logger) *MemRepo {
	r := &MemRepo{buckets: make(map[string]*Bucket, 1024), l: logger}
	r.initCleaner()

	return r
}

func (r *MemRepo) Add(ctx context.Context, key string, capacity uint, rate time.Duration) error {
	r.mutex.Lock()
	b, ok := r.buckets[key]

	if !ok {
		b = &Bucket{
			Capacity:  capacity,
			Remaining: capacity - 1,
			Drop:      time.Now().Add(rate),
			Rate:      rate,
		}

		r.buckets[key] = b
		r.mutex.Unlock()

		return nil
	}
	r.mutex.Unlock()

	if time.Now().After(b.Drop) {
		b.Drop = time.Now().Add(b.Rate)
		b.Remaining = b.Capacity
	}

	if b.Remaining == 0 {
		return errors.New("buffer overflow")
	}

	b.Remaining--

	return nil
}

func (r *MemRepo) Drop(ctx context.Context, keys []string) error {
	for _, key := range keys {
		b, ok := r.buckets[key]
		if !ok {
			continue
		}

		b.Remaining = b.Capacity
	}

	return nil
}

func (r *MemRepo) CleanStorage() error {
	r.mutex.Lock()
	for k := range r.buckets {
		delete(r.buckets, k)
	}

	r.mutex.Lock()

	return nil
}
