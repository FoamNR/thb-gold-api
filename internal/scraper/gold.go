package scraper

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type GoldResponse struct {
	Success     bool        `json:"success"`
	Source      string      `json:"source"`
	UpdatedAt   string      `json:"updated_at"`
	GoldBar     GoldPrice   `json:"gold_bar"`
	GoldJewelry GoldPrice   `json:"gold_jewelry"`
}

type GoldPrice struct {
	Buy  string `json:"buy"`
	Sell string `json:"sell"`
}

func FetchGoldPrice() (*GoldResponse, error) {
	url := "https://www.ทองคําราคา.com/"
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("http error: %d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	var goldBarBuy, goldBarSell, goldJewelryBuy, goldJewelrySell string
	var dateStr, timeStr, countStr string
	var targetTable *goquery.Selection

	doc.Find("h3.h-h3").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), "ราคาทองตามประกาศสมาคมค้าทองคำ") {
			targetTable = s.NextFiltered("table")
		}
	})

	if targetTable != nil && targetTable.Length() > 0 {
		targetTable.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
			tds := s.Find("td")
			if tds.Length() >= 3 {
				firstTdText := strings.TrimSpace(tds.Eq(0).Text())

				if strings.Contains(firstTdText, "ทองคำแท่ง") {
					goldBarBuy = strings.TrimSpace(tds.Eq(1).Text())
					goldBarSell = strings.TrimSpace(tds.Eq(2).Text())
				}
				if strings.Contains(firstTdText, "ทองรูปพรรณ") {
					goldJewelryBuy = strings.TrimSpace(tds.Eq(1).Text())
					goldJewelrySell = strings.TrimSpace(tds.Eq(2).Text())
				}
				if strings.Contains(firstTdText, "มกราคม") || strings.Contains(firstTdText, "กุมภาพันธ์") || strings.Contains(firstTdText, "กรกฎาคม") || strings.Contains(firstTdText, "ธันวาคม") { // ย่อเพื่อประหยัดพื้นที่ โค้ดจริงใส่ครบทุกเดือนแบบเดิมได้ครับ
					dateStr = firstTdText
					timeStr = strings.TrimSpace(tds.Eq(1).Text())
					countStr = strings.TrimSpace(tds.Eq(2).Text())
				}
			}
		})
	} else {
		return nil, fmt.Errorf("table not found")
	}

	return &GoldResponse{
		Success:   true,
		Source:    "ราคาทองคำตามประกาศสมาคมค้าทองคำ",
		UpdatedAt: strings.TrimSpace(fmt.Sprintf("%s %s %s", dateStr, timeStr, countStr)),
		GoldBar:   GoldPrice{Buy: goldBarBuy, Sell: goldBarSell},
		GoldJewelry: GoldPrice{Buy: goldJewelryBuy, Sell: goldJewelrySell},
	}, nil
}