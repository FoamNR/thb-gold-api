FROM golang:alpine

WORKDIR /app

# 1. ก๊อปปี้ไฟล์ทั้งหมดในโปรเจกต์เข้าไปใน Container ก่อนเลย
COPY . .

# 2. 🎯 สั่งให้ Go เช็กไฟล์ทั้งหมด (รวมถึงไฟล์ที่สะกดผิดหรือถูก) แล้วดึงทุก Package ที่จำเป็นมาติดตั้ง
RUN go mod tidy

# 3. สั่ง Build จากตำแหน่งโฟลเดอร์โครงสร้างใหม่
RUN go build -o gold-scraper ./cmd/api/

# 4. สั่งรันโปรแกรมเมื่อคอนเทนเนอร์เริ่มทำงาน
CMD ["./gold-scraper"]