package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"os"

	"gold-scraper/internal/auth"
	"gold-scraper/internal/cache"
	"gold-scraper/internal/limiter"
	"gold-scraper/internal/scraper"
)

func main() {
	// 1. สร้าง Cache และโครงสร้างความปลอดภัย
	goldCache := cache.NewMemoryCache()
	ipLimiter := limiter.NewIPLimiter()

	// 2. 🚀 Start Background Worker (ดึงราคาทองทุกๆ 5 นาทีเก็บเข้าแคช)
	go func() {
		for {
			log.Println("[Worker] กำลังไปดึงข้อมูลราคาทองคำจากเว็บภายนอก...")
			data, err := scraper.FetchGoldPrice()
			if err != nil {
				log.Printf("[Worker Error] ดึงข้อมูลไม่สำเร็จ: %v", err)
			} else {
				jsonData, _ := json.Marshal(data)
				goldCache.Set(jsonData) // เก็บลงแคช
				log.Println("[Worker] อัปเดตข้อมูลแคชล่าสุดเรียบร้อย")
			}
			time.Sleep(5 * time.Minute) // พัก 5 นาทีค่อยไปดึงใหม่
		}
	}()

	// 3. สร้าง Handler สำหรับส่งข้อมูลแคชกลับไปให้ User
	goldHandler := func(w http.ResponseWriter, r *http.Request) {
		cachedData := goldCache.Get()
		if cachedData == nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"success":false,"error":"กำลังดึงข้อมูลระบบครั้งแรก กรุณารอสักครู่"}`))
			return
		}
		w.Write(cachedData)
	}

	// 4. 🛡️ ผูกระบบความปลอดภัย (Rate Limit -> API Key -> ดึงทองจาก Cache)
	secureHandler := limiter.RateLimitMiddleware(ipLimiter, auth.ApiKeyMiddleware(goldHandler))
	http.HandleFunc("/api/gold", secureHandler)

	fmt.Println("--------------------------------------------------")
	fmt.Println("🌟 API ราคาทองคำระบบ Production (Best Practice) เริ่มทำงาน")
	fmt.Println("🔗 URL: http://localhost:8080/api/gold")
	fmt.Println("🔒 ปลอดภัยด้วยระบบ แคชข้อมูล + กันบอทยิงถล่ม + API Key")
	fmt.Println("--------------------------------------------------")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}