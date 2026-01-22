package handlers

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"qr-shorten-go/models"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"

	"github.com/gofiber/fiber/v2"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/skip2/go-qrcode"
	"gorm.io/gorm"
)

type LinkHandler struct {
	DB *gorm.DB
}



// Fungsi membuat Short Link
func (h *LinkHandler) CreateShorten(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uint)
    
    var req struct {
        LongURL string `json:"long_url"`
    }
    c.BodyParser(&req)

    id, _ := gonanoid.New(8)
    shortUrl := c.BaseURL() + "/" + id
	fullShortURL := c.BaseURL() + "/" + id

    // 1. Generate QR Code ke dalam memori ([]byte)
    qrBytes, _ := qrcode.Encode(shortUrl, qrcode.Medium, 256)

    // 2. Upload ke Cloudinary
    cloudinaryURL, err := UploadToCloudinary(qrBytes, "qr_"+id)
	if err != nil {
	    // Tampilkan error aslinya di terminal untuk debug
	    log.Println("Cloudinary Error:", err) 
	    return c.Status(500).JSON(fiber.Map{
	        "error": "Gagal upload QR",
	        "details": err.Error(),
	    })
	}

    // 3. Simpan link + QRURL ke database
    newLink := models.Link{
        ShortCode: id,
        LongURL:   req.LongURL,
		ShortURL:  fullShortURL,
        UserID:    userID,
        QRURL:     cloudinaryURL, // Simpan URL dari Cloudinary
    }
    h.DB.Create(&newLink)

    return c.JSON(fiber.Map{
		"short_url":  id,
		"full_short_url": fullShortURL,
        "qr_url":     cloudinaryURL,
    })
}

func UploadToCloudinary(imageBytes []byte, fileName string) (string, error) {
    // 1. Ambil ENV (Pastikan tidak kosong)
    cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
    apiKey := os.Getenv("CLOUDINARY_API_KEY")
    apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

    if cloudName == "" || apiKey == "" || apiSecret == "" {
        return "", fmt.Errorf("Cloudinary credentials are missing in .env")
    }

    cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
    if err != nil {
        return "", err
    }

    // 2. Upload
    resp, err := cld.Upload.Upload(context.Background(), bytes.NewReader(imageBytes), uploader.UploadParams{
        PublicID: fileName,
        Folder:   "qr_codes",
    })
    
    if err != nil {
        return "", err
    }

    // 3. Return URL yang aman
    return resp.SecureURL, nil
}

// Fungsi Generate QR Code
func (h *LinkHandler) GetQRCode(c *fiber.Ctx) error {
	url := c.Query("url")
	if url == "" {
		return c.Status(400).SendString("URL query is required")
	}

	// Generate QR Code dalam bentuk PNG
	var png []byte
	png, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		return c.Status(500).SendString("Failed to generate QR")
	}

	c.Set("Content-Type", "image/png")
	return c.Send(png)
}

// Fungsi Redirect (Jika orang klik link pendeknya)
// Fungsi Redirect (Jika orang klik link pendeknya)
func (h *LinkHandler) Resolve(c *fiber.Ctx) error {
    code := c.Params("code")
    var link models.Link

    // 1. Cari di database berdasarkan ShortCode
    // Gunakan .First() untuk mencari satu data
    if err := h.DB.Where("short_code = ?", code).First(&link).Error; err != nil {
        return c.Status(404).SendString("Link tidak ditemukan")
    }

    // 2. Tambah jumlah klik secara atomik (Aman dari race condition)
    // Ini akan menjalankan: UPDATE links SET click_count = click_count + 1 WHERE id = ...
    h.DB.Model(&link).UpdateColumn("click_count", gorm.Expr("click_count + ?", 1))

    // 3. Redirect ke URL asli (301 untuk permanen, 302 untuk sementara)
    // Disarankan 302 agar browser tidak men-cache redirect, 
    // supaya setiap klik selalu masuk ke backend kita dan terhitung.
    return c.Redirect(link.LongURL, 302) 
}

func (h *LinkHandler) GetUserLinks(c *fiber.Ctx) error {
    // Ambil userID dari token JWT yang sudah di-parse middleware
    userID := c.Locals("user_id").(uint)

    var links []models.Link
    // Cari semua link berdasarkan UserID
    if err := h.DB.Where("user_id = ?", userID).Find(&links).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data"})
    }

    return c.JSON(fiber.Map{
        "status": "success",
        "data":   links, // Di sini nanti ada field ShortCode, LongURL, dan QRURL (Cloudinary)
    })
}

func (h *LinkHandler) DeleteLink(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uint)
    linkID := c.Params("id")

    // Hapus link yang ID-nya sesuai DAN dimiliki oleh user yang sedang login
    // Ini penting agar user A tidak bisa menghapus link milik user B
    result := h.DB.Where("id = ? AND user_id = ?", linkID, userID).Delete(&models.Link{})

    if result.Error != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Gagal menghapus link"})
    }

    if result.RowsAffected == 0 {
        return c.Status(404).JSON(fiber.Map{"error": "Link tidak ditemukan atau bukan milik Anda"})
    }

    return c.JSON(fiber.Map{"message": "Link berhasil dihapus"})
}

func (h *LinkHandler) GetStats(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uint)

    // Definisikan struct di dalam fungsi agar rapi
    // Pastikan tag json menggunakan snake_case (huruf kecil)
    var stats struct {
        TotalLinks  int64         `json:"total_links"`
        TotalClicks int64         `json:"total_clicks"`
        Data        []models.Link `json:"data"` 
    }

    // 1. Hitung total link (Count)
    h.DB.Model(&models.Link{}).Where("user_id = ?", userID).Count(&stats.TotalLinks)

    // 2. Hitung total klik (Sum)
    // Gunakan COALESCE agar jika belum ada klik, hasilnya 0 (bukan null)
    h.DB.Model(&models.Link{}).Where("user_id = ?", userID).
        Select("COALESCE(SUM(click_count), 0)").
        Scan(&stats.TotalClicks)

    // 3. ISI ARRAY DATA (Penyebab Chart Kosong)
    // Kita harus memanggil database untuk mengambil recordnya
    if err := h.DB.Where("user_id = ?", userID).Find(&stats.Data).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Database error"})
    }

    // Kirim hasilnya
    return c.JSON(stats)
}