package handlers_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/0xMoonrise/gochive/internal/config"
	"github.com/0xMoonrise/gochive/internal/core"
	"github.com/0xMoonrise/gochive/internal/handlers"
	"github.com/0xMoonrise/gochive/internal/server"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var PDF []byte = []byte("%PDF-1.4\n" +
	"1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n" +
	"2 0 obj\n<< /Type /Pages /Kids [3 0 R] /Count 1 >>\nendobj\n" +
	"3 0 obj\n<< /Type /Page /Parent 2 0 R /MediaBox [0 0 200 200] /Contents 4 0 R /Resources << >> >>\nendobj\n" +
	"4 0 obj\n<< /Length 0 >>\nstream\n\nendstream\nendobj\n" +
	"xref\n0 5\n" +
	"0000000000 65535 f \n" +
	"0000000009 00000 n \n" +
	"0000000058 00000 n \n" +
	"0000000115 00000 n \n" +
	"0000000219 00000 n \n" +
	"trailer\n<< /Size 5 /Root 1 0 R >>\n" +
	"startxref\n268\n" +
	"%%EOF\n")

func projectRoot() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..")
}

func TestMain(m *testing.M) {
	root := projectRoot()

	if err := godotenv.Load(filepath.Join(root, ".env")); err != nil {
		slog.Warn("no .env file found, relying on real env vars")
	}

	code := m.Run()
	os.Exit(code)
}

func newTestConfig(t *testing.T) *config.Config {
	t.Helper()
	return &config.Config{
		Mode: config.FS,
		Data: t.TempDir() + "/",
		FS: config.FSClientConfig{
			Root: t.TempDir() + "/",
		},
	}
}

func setupTestApp(t *testing.T) (*core.App, *gin.Engine) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	app := core.NewApp()
	app.Config = newTestConfig(t)

	err := app.Run(
		core.StageDB,
		core.StageStorage,
	)
	assert.NoError(t, err)

	r := server.NewEngine()
	t.Cleanup(func() { app.Cleanup() })
	return app, r
}

func TestSetup(t *testing.T) {
	_, r := setupTestApp(t)

	r.GET("/", handlers.Root)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUploadFile(t *testing.T) {
	app, r := setupTestApp(t)

	r.POST("/upload", handlers.UploadFile(app))

	tests := []struct {
		name       string
		fieldName  string
		fileName   string
		content    []byte
		wantStatus int
	}{
		{
			name:       "valid pdf file",
			fieldName:  "file",
			fileName:   uuid.New().String() + ".pdf",
			content:    PDF,
			wantStatus: http.StatusOK,
		},
		{
			name:       "valid markdown file",
			fieldName:  "file",
			fileName:   uuid.New().String() + ".md",
			content:    []byte("# Hello world"),
			wantStatus: http.StatusOK,
		},
		{
			name:       "empty pdf file",
			fieldName:  "file",
			fileName:   "empty.pdf",
			content:    []byte{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty md file",
			fieldName:  "file",
			fileName:   "empty.md",
			content:    []byte{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "without field name",
			fieldName:  "",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "name with path traversal",
			fieldName:  "file",
			fileName:   "../../etc/passwd",
			content:    []byte(PDF),
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			if tt.fieldName != "" {
				part, err := writer.CreateFormFile(tt.fieldName, tt.fileName)
				assert.NoError(t, err)
				_, err = part.Write(tt.content)
				assert.NoError(t, err)
			}
			assert.NoError(t, writer.Close())

			req := httptest.NewRequest(http.MethodPost, "/upload", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code, w.Body.String())
		})
	}
}

type UploadResponse struct {
	Success bool         `json:"success"`
	File    FileResponse `json:"file"`
}

type FileResponse struct {
	ID        int64  `json:"id"`
	Filename  string `json:"filename"`
	Editorial string `json:"editorial"`
	Favorite  bool   `json:"favorite"`
}

func sameFile(t *testing.T, fileA []byte, fileB []byte) {
	original := sha256.Sum256(fileA)
	uploaded := sha256.Sum256(fileB)
	assert.Equal(t, original, uploaded)
}

func uploadTestFile(t *testing.T, r *gin.Engine, filename string, content []byte) UploadResponse {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	assert.NoError(t, err)
	_, err = part.Write(content)
	assert.NoError(t, err)
	assert.NoError(t, writer.Close())

	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, w.Body.String())

	var resp UploadResponse
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	return resp
}

func TestIntegrity(t *testing.T) {
	app, r := setupTestApp(t)

	r.POST("/upload", handlers.UploadFile(app))

	t.Log("--- testing pdf file upload")
	res := uploadTestFile(t, r, uuid.New().String()+".pdf", PDF)
	key := path.Join("files", strconv.FormatInt(res.File.ID, 10))

	pdfFile, err := app.Storage.GetItem(t.Context(), key)
	assert.NoError(t, err)
	defer pdfFile.Reader.Close()
	data, err := io.ReadAll(pdfFile.Reader)
	assert.NoError(t, err)

	sameFile(t, PDF, data)

	t.Log("--- testing MD file upload")
	res = uploadTestFile(t, r, uuid.New().String()+".md", []byte("# Hello world"))
	key = path.Join("files", strconv.FormatInt(res.File.ID, 10))

	mdFile, err := app.Storage.GetItem(t.Context(), key)
	assert.NoError(t, err)
	defer mdFile.Reader.Close()
	data, err = io.ReadAll(mdFile.Reader)
	assert.NoError(t, err)

	sameFile(t, []byte("# Hello world"), data)
}

func TestImageGeneration(t *testing.T) {
	app, r := setupTestApp(t)

	r.GET("/images/:id", handlers.GetImage(app))
	r.POST("/upload", handlers.UploadFile(app))

	t.Log("--- generating image from a pdf file")
	res := uploadTestFile(t, r, uuid.New().String()+".pdf", PDF)
	imageId := strconv.FormatInt(res.File.ID, 10)

	key := path.Join("images", imageId)
	image, err := app.Storage.GetItem(t.Context(), key)
	assert.NoError(t, err)

	data, err := io.ReadAll(image.Reader)
	assert.NoError(t, err)

	sniffLen := min(len(data), 512)
	assert.Equal(t, "image/webp", http.DetectContentType(data[:sniffLen]))

	req := httptest.NewRequest(http.MethodGet, "/images/"+imageId, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, w.Body.String())

	sameFile(t, data, w.Body.Bytes())

	t.Log("--- generating image from a md file")
	res = uploadTestFile(t, r, uuid.New().String()+".md", []byte("# Hello world"))
	imageId = strconv.FormatInt(res.File.ID, 10)

	key = path.Join("images", imageId)
	image, err = app.Storage.GetItem(t.Context(), key)
	assert.Error(t, err)

	req = httptest.NewRequest(http.MethodGet, "/images/"+imageId, nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code, w.Body.String())
}

func editFile(t *testing.T, r *gin.Engine, id int64, filename, editorial string) *httptest.ResponseRecorder {
	t.Helper()

	form := url.Values{}
	form.Set("filename", filename)
	form.Set("editorial", editorial)

	req := httptest.NewRequest(http.MethodPatch, "/edit/"+strconv.FormatInt(id, 10), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// I'm going to assume that, whether it's Markdown or PDF,
// it'll behave the same way in this endpoint /edit/:id
func TestEditFile(t *testing.T) {
	app, r := setupTestApp(t)

	r.POST("/upload", handlers.UploadFile(app))
	r.PATCH("/edit/:id", handlers.SetEditFile(app))

	res := uploadTestFile(t, r, uuid.New().String()+".pdf", PDF)
	assert.Equal(t, "Default", res.File.Editorial)

	newFilename := uuid.New().String() + ".pdf"
	w := editFile(t, r, res.File.ID, newFilename, "New Editorial")
	assert.Equal(t, http.StatusOK, w.Code, w.Body.String())

	updated, err := app.DB.Queries.GetArchive(t.Context(), int(res.File.ID))
	assert.NoError(t, err)
	assert.Equal(t, newFilename, updated.Filename)
	assert.Equal(t, "New Editorial", updated.Editorial)
}

func TestEditFile_InvalidExtension(t *testing.T) {
	app, r := setupTestApp(t)

	r.POST("/upload", handlers.UploadFile(app))
	r.PATCH("/edit/:id", handlers.SetEditFile(app))

	res := uploadTestFile(t, r, uuid.New().String()+".pdf", PDF)

	w := editFile(t, r, res.File.ID, "malicious.exe", "New Editorial")
	assert.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())

	res = uploadTestFile(t, r, uuid.New().String()+".md", []byte("# Hello world"))

	w = editFile(t, r, res.File.ID, "malicious.exe", "New Editorial")
	assert.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())

}

func TestSetFavorite(t *testing.T) {
	app, r := setupTestApp(t)

	r.POST("/upload", handlers.UploadFile(app))
	r.POST("/set_favorite/:id", handlers.SetFavorite(app))

	res := uploadTestFile(t, r, uuid.New().String()+".pdf", PDF)
	form := url.Values{}
	form.Set("favorite", "true")

	req := httptest.NewRequest(http.MethodPost,
		"/set_favorite/"+strconv.FormatInt(res.File.ID, 10),
		strings.NewReader(form.Encode()),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, w.Body.String())

	archive, err := app.DB.Queries.GetArchive(context.Background(), int(res.File.ID))
	assert.NoError(t, err)
	assert.Equal(t, true, archive.Favorite)
}
