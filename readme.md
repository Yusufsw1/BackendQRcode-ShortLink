# 🚀 QR Shorten Pro - Backend (Go API)

API Service built using Go Fiber to shorten URLs, generate QR Codes to Cloudinary, and track analytics.

## 🛠 Tech Stack

- **Language:** Go 1.21+
- **Framework:** Fiber v2
- **ORM:** GORM
- **Database:** PostgreSQL
- **Auth:** Google OAuth2 & JWT
- **Cloud Storage:** Cloudinary (for QR Code)

## 📦 Features

- URL Shortening with NanoID.
- Click Tracking (Atomically incremented).
- Cloudinary Upload integration.
- Protected Routes using JWT Middleware.
- CORS Configured for Frontend integration.

## 🚀 Setup & Installation

1. **Clone Repository:**

   ```bash
   git clone [https://github.com/username/qr-shorten-go.git](https://github.com/username/qr-shorten-go.git)
   cd qr-shorten-go

   ```

2. **Setup .env: Buat file .env di root folder:**

   ```bash
   DB_URL=postgres://user:pass@localhost:5432/dbname
   GOOGLE_CLIENT_SECRET=xxx
   GOOGLE_CLIENT_ID=xxx
   JWT_SECRET=your_jwt_secret
   CLOUDINARY_CLOUD_NAME=xxx
   CLOUDINARY_API_KEY=xxx
   CLOUDINARY_API_SECRET=xxx
   FRONTEND_URL=http://localhost:3000

   ```

3. **Run Application:**
   ```bash
   go mod tidy
   go run main.go
   ```
