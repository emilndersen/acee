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
	maxUploadSize = 20 << 20 // 20 MB
	thumbWidth    = 400
	thumbHeight   = 0 // auto
)

type Handler struct {
	uploadDir string
}

func NewHandler(uploadDir string) *Handler {
	return &Handler{uploadDir: uploadDir}
}

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		response.WriteError(w, http.StatusBadRequest, "file too large (max 20MB)")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "file is required")
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		response.WriteError(w, http.StatusBadRequest, "only image files allowed")
		return
	}

	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".jpg"
	}
	ext = strings.ToLower(ext)

	id := uuid.New().String()
	origName := id + ext
	thumbName := id + "_thumb" + ext

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

	go h.generateThumb(origPath, filepath.Join(thumbDir, thumbName))

	response.WriteJSON(w, http.StatusCreated, map[string]any{
		"ok": true,
		"file": map[string]any{
			"id":        id,
			"filename":  header.Filename,
			"image_url": fmt.Sprintf("/uploads/original/%s", origName),
			"thumb_url": fmt.Sprintf("/uploads/thumbs/%s", thumbName),
			"size":      header.Size,
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

	thumb := imaging.Resize(img, thumbWidth, thumbHeight, imaging.Lanczos)

	out, err := os.Create(dstPath)
	if err != nil {
		return
	}
	defer out.Close()

	jpeg.Encode(out, thumb, &jpeg.Options{Quality: 85})
}
