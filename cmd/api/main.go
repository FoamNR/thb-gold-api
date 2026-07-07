package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os" // 🎯 นำเข้า os
	"time"

	"gold-scraper/internal/auth"
	"gold-scraper/internal/cache"
	"gold-scraper/internal/limiter"
	"gold-scraper/internal/scraper"
)

func main() {
	goldCache := cache.NewMemoryCache()
	ipLimiter := limiter.NewIPLimiter()

	go func() {
		for {
			log.Println("[Worker] กำลังไปดึงข้อมูลราคาทองคำจากเว็บภายนอก...")
			data, err := scraper.FetchGoldPrice()
			if err != nil {
				log.Printf("[Worker Error] ดึงข้อมูลไม่สำเร็จ: %v", err)
			} else {
				jsonData, _ := json.Marshal(data)
				goldCache.Set(jsonData)
				log.Println("[Worker] อัปเดตข้อมูลแคชล่าสุดเรียบร้อย")
			}
			time.Sleep(5 * time.Minute)
		}
	}()

	goldHandler := func(w http.ResponseWriter, r *http.Request) {
		cachedData := goldCache.Get()
		if cachedData == nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"success":false,"error":"กำลังดึงข้อมูลระบบครั้งแรก กรุณารอสักครู่"}`))
			return
		}
		w.Write(cachedData)
	}

	secureHandler := limiter.RateLimitMiddleware(ipLimiter, auth.ApiKeyMiddleware(goldHandler))
	http.HandleFunc("/api/gold", secureHandler)

	// 🎯 🎯 🎯 จุดสำคัญ: ตรวจสอบให้มั่นใจว่าโค้ดช่วงนี้ถูกเขียนแบบนี้ เพื่อดึง 'os' มาใช้งานจริง
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("--------------------------------------------------")
	fmt.Printf("🚀 API Server กำลังทำงานบนพอร์ต %s...\n", port)
	fmt.Println("--------------------------------------------------")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}