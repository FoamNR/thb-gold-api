# 🌟 Thai Gold Price API

A clean and secure Web Scraping API built with Go to fetch real-time gold prices from the Gold Traders Association of Thailand. Designed following Go Standard Layout practices.

## ✨ Features
* **In-Memory Cache:** Runs a background worker to scrape data every 5 minutes, preventing your IP from being banned.
* **Rate Limiter:** Limits requests to a maximum of 20 requests/minute per IP to prevent DoS attacks.
* **API Key Protection:** Secures endpoints using a custom header check (`X-API-Key`).

---

## 🚀 Installation & Setup

### 1. Prerequisites
Ensure you have **Docker** and **WSL** (for Windows users) installed.

### 2. Configure Environment (.env)
Create a `.env` file in the root project folder and define your API key and Url:
```env GOLD_API_KEY=my-super-secret-key-123
URL_GOLD=https://www.ทองคําราคา.com/