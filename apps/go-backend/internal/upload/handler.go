package upload

import (
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"

	"github.com/emilndersen/acee/apps/go-backend/internal/response"
)

const (
	maxBodySize = 10 << 30 // 10 GB
	thumbWidth  = 400
)

var videoExts = map[string]bool{
	".mp4": true, ".mov": true, ".webm": true, ".avi": true, ".mkv": true,
}

type Handler struct {
	uploadDir string
}

func NewHandler(uploadDir string) *Handler {
	return &Handler{uploadDir: uploadDir}
}

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	if err := r.ParseMultipartForm(64 << 20); err != nil {
		response.WriteError(w, http.StatusBadRequest, "file too large")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "file is required")
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext == "" {
		ext = ".bin"
	}

	contentType := header.Header.Get("Content-Type")
	isImage := strings.HasPrefix(contentType, "image/")
	isVideo := strings.HasPrefix(contentType, "video/") || videoExts[ext]

	if !isImage && !isVideo {
		response.WriteError(w, http.StatusBadRequest, "only image or video files allowed")
		return
	}

	mediaType := "image"
	if isVideo {
		mediaType = "video"
	}

	id := uuid.New().String()
	origName := id + ext

	origDir := filepath.Join(h.uploadDir, "original")
	thumbDir := filepath.Join(h.uploadDir, "thumbs")
	os.MkdirAll(origDir, 0755)
	os.MkdirAll(thumbDir, 0755)

	origPath := filepath.Join(origDir, origName)
	dst, err := os.Create(origPath)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to save file")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to save file")
		return
	}

	thumbURL := ""
	if isImage {
		thumbName := id + "_thumb.jpg"
		go h.generateThumb(origPath, filepath.Join(thumbDir, thumbName))
		thumbURL = fmt.Sprintf("/uploads/thumbs/%s", thumbName)
	}

	response.WriteJSON(w, http.StatusCreated, map[string]any{
		"ok": true,
		"file": map[string]any{
			"id":         id,
			"filename":   header.Filename,
			"image_url":  fmt.Sprintf("/uploads/original/%s", origName),
			"thumb_url":  thumbURL,
			"media_type": mediaType,
			"size":       header.Size,
		},
	})
}

func (h *Handler) generateThumb(srcPath, dstPath string) {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return
	}
	defer srcFile.Close()

	img, _, err := image.Decode(srcFile)
	if err != nil {
		return
	}

	thumb := imaging.Resize(img, thumbWidth, 0, imaging.Lanczos)

	out, err := os.Create(dstPath)
	if err != nil {
		return
	}
	defer out.Close()

	jpeg.Encode(out, thumb, &jpeg.Options{Quality: 85})
}
