# 🚀 QR Shorten Pro - Backend (Go API)

A high-performance API service built with Go Fiber to shorten URLs, generate QR Codes via Cloudinary integration, and track real-time analytics.

## 🛠 Tech Stack

- **Language:** Go 1.21+
- **Framework:** [Fiber v2](https://gofiber.io/)
- **ORM:** [GORM](https://gorm.io/)
- **Database:** PostgreSQL
- **Auth:** Google OAuth2 & JWT (JSON Web Tokens)
- **Cloud Storage:** [Cloudinary](https://cloudinary.com/) (For QR Code image hosting)

## 📦 Key Features

- **Dynamic Shortening:** Generates unique 8-character codes using NanoID.
- **Atomic Click Tracking:** Increments click counts directly in the database to prevent race conditions.
- **QR Code Generation:** Automatically generates QR codes for every shortened link and uploads them to the cloud.
- **Protected Endpoints:** Secure routes using custom JWT Middleware.
- **CORS Enabled:** Fully configured for seamless integration with modern frontend frameworks.

## 🚀 Setup & Installation

1. **Clone the Repository:**

   ```bash
   git clone [https://github.com/yourusername/qr-shorten-go.git](https://github.com/yourusername/qr-shorten-go.git)
   cd qr-shorten-go

   ```

2. **Setup .env: Create files .env in root folder:**

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
