package rate

import (
	"sync"
	"time"
)

// Простой токен-бакет на ключ (обычно host/endpoint).
type bucket struct {
	capacity int
	tokens   float64
	fillRate float64 // tokens per second
	last     time.Time
}

type Limiter struct {
	mu sync.Mutex
	m  map[string]*bucket
}

func New() *Limiter { return &Limiter{m: map[string]*bucket{}} }

// rps - множетель заполнения токенов (запросов) в секунду
// burst - сколько доступно максимально токенов (запросов)
// То есть при rps = 2.0 burst = 2 мы получаем дефолтно 2 запроса сразу и каждые половину секунды (1/2.0) будет востанавливаться запрос
// При 0.5 и 1 мы будем иметь 1 запрос в 2 секунды (1/0.5)
func (l *Limiter) Configure(key string, rps float64, burst int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.m[key] = &bucket{capacity: burst, tokens: float64(burst), fillRate: rps, last: time.Now()}
}

func (l *Limiter) getBucket(key string) *bucket {
	b, ok := l.m[key]
	if ok {
		return b
	}

	l.m[key] = &bucket{capacity: 2, tokens: 2, fillRate: 2, last: time.Now()}
	return l.m[key]
}

func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	b := l.getBucket(key)
	now := time.Now()
	elapsed := now.Sub(b.last).Milliseconds()
	b.tokens += (float64(elapsed) / 1000) * b.fillRate
	if b.tokens > float64(b.capacity) {
		b.tokens = float64(b.capacity)
	}
	b.last = now
	if b.tokens >= 1 {
		b.tokens -= 1
		return true
	}
	return false
}

func (l *Limiter) RateLimited(key string, d time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()

	b := l.getBucket(key)
	b.tokens -= d.Seconds() * b.fillRate
}

func (l *Limiter) Wait(key string) {
	for {
		if l.Allow(key) {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
}
