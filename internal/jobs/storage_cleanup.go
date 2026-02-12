package jobs

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

type StorageCleaner struct {
	Dirs   []string
	MaxAge time.Duration
}

func StorageCleanup(dirs []string, maxAge time.Duration) *StorageCleaner {
	return &StorageCleaner{
		Dirs:   dirs,
		MaxAge: maxAge,
	}
}

// Fungsi internal untuk bersih-bersih
func (s *StorageCleaner) clean() {
	for _, dir := range s.Dirs {
		// Gunakan WalkDir (lebih efisien dari Walk di Go modern)
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			// 1. Handle Error akses file
			if err != nil {
				log.Println("⚠️ Error accessing:", path, err)
				return nil // Lanjut ke file berikutnya
			}

			// 2. Skip jika itu adalah Folder (termasuk root folder temp/)
			if info.IsDir() {
				return nil
			}

			// 3. Logic Cek Umur File (Hanya jalan untuk FILE)
			if time.Since(info.ModTime()) > s.MaxAge {
				if err := os.Remove(path); err == nil {
					log.Println("🗑️ Deleted old file:", path)
				} else {
					log.Println("❌ Failed to delete:", path, err)
				}
			}
			return nil
		})
	}
}

func (s *StorageCleaner) Start() {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			s.clean()
		}
	}()
}
