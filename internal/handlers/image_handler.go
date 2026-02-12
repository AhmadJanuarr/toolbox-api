package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"toolkits/internal/services"

	"github.com/gin-gonic/gin"
)

const (
	MaxFileSize = 5 << 20 // 5MB
	UploadDir   = "./temp/uploads"
)

// Allowed formats whitelist
var allowedFormats = map[string]bool{
	"jpg":  true,
	"jpeg": true,
	"png":  true,
	"webp": true,
}

func ConvertImage(c *gin.Context) {
	// 1. Validasi Input Form (Target Format)
	targetFormat := strings.ToLower(c.PostForm("targetFormat"))
	if targetFormat == "" || !allowedFormats[targetFormat] {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Format target tidak valid atau tidak didukung (gunakan: jpg, jpeg, png, webp)",
		})
		return
	}

	// 2. Menerima File
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "File tidak ditemukan atau tidak valid",
			"error":   err.Error(),
		})
		return
	}

	// 3. Validasi Ukuran File
	if file.Size > MaxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": fmt.Sprintf("Ukuran file tidak boleh lebih dari %d MB", MaxFileSize>>20),
		})
		return
	}

	// 4. Generate Nama File Aman (Mencegah Path Traversal & Overwrite)
	// Format: timestamp-filename_asli.ext
	safeFilename := fmt.Sprintf("%d-%s", time.Now().Unix(), filepath.Base(file.Filename))
	srcPath := filepath.Join(UploadDir, safeFilename)

	// 5. Simpan File Sementara
	if err := c.SaveUploadedFile(file, srcPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Gagal menyimpan file sementara",
			"error":   err.Error(),
		})
		return
	}

	// 6. Proses Konversi
	resultPath, err := services.ProcessImageConversion(srcPath, targetFormat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Gagal mengonversi gambar",
			"error":   err.Error(),
		})
		defer os.Remove(srcPath)
		return
	}

	// 7. Response Sukses
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Konversi berhasil",
		"data": gin.H{
			"original_file": file.Filename,
			"result_path":   resultPath,
			"format":        targetFormat,
		},
	})
}

func CompressionImage(c *gin.Context) {

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "File tidak ditemukan atau tidak valid",
			"error":   err.Error(),
		})
		return
	}

	qualityStr := c.PostForm("quality")
	if qualityStr == "" {
		qualityStr = "80"
	}

	quality, err := strconv.Atoi(qualityStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Quality harus berupa angka",
			"error":   err.Error(),
		})
		return
	}

	safeFilename := fmt.Sprintf("%d-%s", time.Now().Unix(), filepath.Base(file.Filename))
	srcPath := filepath.Join(UploadDir, safeFilename)

	if err := c.SaveUploadedFile(file, srcPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Gagal menyimpan file sementara",
			"error":   err.Error(),
		})
		return
	}

	resultPath, err := services.ProcessImageCompression(srcPath, quality)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Gagal mengompres gambar",
			"error":   err.Error(),
		})
		defer os.Remove(srcPath)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Kompresi berhasil",
		"data": gin.H{
			"original_file": file.Filename,
			"result_path":   resultPath,
			"quality":       quality,
		},
	})
}

func ResizeImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "File wajib diupload", "error": err.Error()})
		return
	}

	widthStr := c.PostForm("width")
	heightStr := c.PostForm("height")

	width, errW := strconv.Atoi(widthStr)
	height, errH := strconv.Atoi(heightStr)

	if errW != nil || errH != nil || width <= 0 || height <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Width dan Height harus berupa angka positif",
		})
		return
	}

	safeFilename := fmt.Sprintf("%d-%s", time.Now().Unix(), filepath.Base(file.Filename))
	srcPath := filepath.Join(UploadDir, safeFilename)

	os.MkdirAll(UploadDir, 0755)

	if err := c.SaveUploadedFile(file, srcPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Gagal menyimpan file", "error": err.Error()})
		return
	}

	resultPath, err := services.ProcessImageResize(srcPath, width, height)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Gagal resize gambar", "error": err.Error()})
		return
	}
	defer os.Remove(srcPath)

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Resize berhasil",
		"data": gin.H{
			"original_file": file.Filename,
			"result_path":   resultPath,
			"width":         width,
			"height":        height,
		},
	})
}
