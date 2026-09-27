package utils

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
)

func FileUpload(c *gin.Context, file *multipart.FileHeader, prefix string) (string, error) {
	// Validasi Ukuran File (Maksimal 2MB untuk avatar/gambar standar, atau 5MB jika banner)
	maxSize := int64(2 * 1024 * 1024)
	if prefix == "event" || prefix == "banner" || prefix == "thumbnail" {
		maxSize = 5 * 1024 * 1024
	}

	if file.Size > maxSize {
		if maxSize == 5*1024*1024 {
			return "", errors.New("ukuran file terlalu besar (maksimal 5MB)")
		}
		return "", errors.New("ukuran file terlalu besar (maksimal 2MB)")
	}

	// Validasi Ekstensi Gambar
	ext := filepath.Ext(file.Filename)
	re := regexp.MustCompile(`(?i)\.(png|jpg|jpeg|webp)$`)
	if !re.MatchString(ext) {
		return "", errors.New("format file tidak didukung (hanya PNG, JPG, JPEG, WEBP)")
	}

	// folder public sebagai direktori utama
	uploadDir := "public"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", errors.New("gagal membuat direktori penyimpanan file")
	}

	// Generate nama file unik
	filename := fmt.Sprintf("%s_%d%s", prefix, time.Now().UnixNano(), ext)
	location := filepath.Join(uploadDir, filename)

	// Simpan file menggunakan helper Gin
	if err := c.SaveUploadedFile(file, location); err != nil {
		return "", errors.New("gagal menyimpan file ke server")
	}

	return filename, nil
}
