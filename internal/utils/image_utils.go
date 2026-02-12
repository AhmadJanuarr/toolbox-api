package utils

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"github.com/chai2010/webp"
)

func LoadImage(inputPath string) (image.Image, string, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, " ", fmt.Errorf("Gagal membuka file input : %v", err)
	}
	defer file.Close()
	img, format, err := image.Decode(file)
	if err != nil {
		return nil, " ", fmt.Errorf("Gagal decode gambar : %v", err)
	}
	return img, format, nil
}

func SaveImage(img image.Image, outputPath string, format string, quality int) error {
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("Gagal membuat direktori: %v", err)
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("Gagal membuat file output: %v", err)
	}
	defer outputFile.Close()

	var encodeErr error
	switch format {
	case "jpeg", "jpg":
		encodeErr = jpeg.Encode(outputFile, img, &jpeg.Options{Quality: quality})
	case "png":
		encoder := png.Encoder{CompressionLevel: png.BestCompression}
		encodeErr = encoder.Encode(outputFile, img)
	case "webp":
		encodeErr = webp.Encode(outputFile, img, &webp.Options{Lossless: false, Quality: float32(quality)})
	default:
		encodeErr = fmt.Errorf("format target tidak didukung: %s (gunakan jpg, png, atau webp)", format)
	}
	if encodeErr != nil {
		os.Remove(outputPath)
		return fmt.Errorf("gagal encode ke format %s: %v", format, encodeErr)
	}
	return nil
}
