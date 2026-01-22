package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"qr-shorten-go/models"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

// Konfigurasi Google OAuth
func GetGoogleConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Endpoint:     google.Endpoint,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	}
}

func GenerateJWT(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // Token expired dalam 3 hari
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// Handler untuk Login
func LoginGoogle(c *fiber.Ctx) error {
	url := GetGoogleConfig().AuthCodeURL("random-state")
	return c.Redirect(url)
}

// Handler untuk Callback (Data balik dari Google)
func CallbackGoogle(c *fiber.Ctx, db *gorm.DB) error {
    code := c.Query("code")
    conf := GetGoogleConfig()

    // 1. Tukar code dengan token Google
    tokenGoogle, err := conf.Exchange(context.Background(), code)
    if err != nil {
        return c.Status(500).SendString("Gagal tukar token")
    }

    // 2. Ambil data user dari Google API
    res, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + tokenGoogle.AccessToken)
    if err != nil {
        return c.Status(500).SendString("Gagal menghubungi Google API")
    }
    defer res.Body.Close()

    var gUser struct {
        ID    string `json:"id"`
        Email string `json:"email"`
        Name  string `json:"name"`
    }

    if err := json.NewDecoder(res.Body).Decode(&gUser); err != nil {
        return c.Status(500).SendString("Gagal membaca data profil user")
    }

    // 3. Simpan ke Postgres
    user := models.User{
        Email:    gUser.Email,
        Name:     gUser.Name,
        GoogleID: gUser.ID,
    }

    if err := db.Where(models.User{Email: user.Email}).FirstOrCreate(&user).Error; err != nil {
        return c.Status(500).SendString("Gagal menyimpan ke database")
    }

    // --- UPDATE DI SINI (BAGIAN JWT) ---
    // 4. Generate JWT Token internal kita sendiri
    tokenJWT, err := GenerateJWT(user.ID) // Pastikan fungsi GenerateJWT sudah kamu buat
    if err != nil {
        return c.Status(500).SendString("Gagal generate akses token")
    }

    // 5. Kirim data ke Next.js (Token inilah yang nanti disimpan di Cookie/LocalStorage)
    frontendURL := "http://localhost:3000/dashboard?token=" + tokenJWT
    
    return c.Redirect(frontendURL)
}

func (h *LinkHandler) GetProfile(c *fiber.Ctx) error {
	// Ambil user_id dari locals (setelah melewati AuthGuard)
	userID := c.Locals("user_id").(uint)

	var user models.User
	// Cari user berdasarkan ID, hilangkan field sensitif jika perlu
	if err := h.DB.First(&user, userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "User tidak ditemukan",
		})
	}

	return c.JSON(fiber.Map{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}