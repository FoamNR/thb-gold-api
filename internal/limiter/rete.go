package limiter

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type IPlimiter struct {
	ips map[string]*rate.Limiter
	mu  sync.Mutex
}

func NewIPLimiter() *IPlimiter {
	return &IPlimiter{ips: make(map[string]*rate.Limiter)}
}

func (i *IPlimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	l, exists := i.ips[ip]
	if !exists {
		// ให้ยิงได้เฉลี่ย 20 ครั้งต่อนาที (และยิงรัวสุดพร้อมกันได้ไม่เกิน 5 ครั้ง)
		l = rate.NewLimiter(rate.Every(time.Minute/20), 5)
		i.ips[ip] = l
	}
	return l
}

// RateLimitMiddleware ดักจับตาม IP Address
func RateLimitMiddleware(limiter *IPlimiter, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}

		l := limiter.GetLimiter(ip)
		if !l.Allow() {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"success":false,"error":"ยิงถี่เกินไป กรุณารอเว้นระยะ (Too Many Requests)"}`))
			return
		}
		next(w, r)
	}
}