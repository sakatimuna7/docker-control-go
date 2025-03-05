package controllers

import (
	"docker-control-go/src/helpers"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// FileEntry mewakili satu file atau folder
type FileEntry struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Type string `json:"type"` // "file" atau "directory"
}

// ListFilesHandler mengembalikan daftar file dalam direktori
func ListFilesHandler(c *fiber.Ctx) error {
	dir := c.Query("path", "/") // Default ke root "/"

	entries, err := os.ReadDir(dir)
	if err != nil {
		return helpers.ErrorResponse(c, 500, "Failed to read directory", err)
	}

	var files []FileEntry
	for _, entry := range entries {
		fileType := "directory"

		if !entry.IsDir() {
			fileType = strings.TrimPrefix(filepath.Ext(entry.Name()), ".") // Ambil ekstensi tanpa titik
			if fileType == "" {
				fileType = "file" // Jika tidak ada ekstensi
			}
		}

		files = append(files, FileEntry{
			Name: entry.Name(),
			Path: filepath.Join(dir, entry.Name()),
			Type: fileType,
		})
	}

	return helpers.SuccessResponse(c, 200, "Files retrieved", files)
}
