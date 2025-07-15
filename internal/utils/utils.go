package utils

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

func MakeThumbnail(data []byte, thumb *[]byte, filename string) error {

	thumbnail, err := generateWebpThumbnail(data, "static/thumbnails/")
	err = saveThumbnailToStatic(thumbnail, filename)

	if err != nil {
		return err
	}
	*thumb = thumbnail

	return nil
}

func ValidateFilename(filename string) bool {
	match, _ := regexp.MatchString("^.+(pdf|md)$", filename)
	return match
}

func saveThumbnailToStatic(thumbnailData []byte, filename string) error {
	dir := "static/thumbnails"

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	path := filepath.Join(dir, filename)

	return os.WriteFile(path, thumbnailData, 0644)
}

func generateWebpThumbnail(pdfBytes []byte, path string) ([]byte, error) {

	pdfFile, err := os.CreateTemp("", "input-*.pdf")

	if err != nil {
		return nil, err
	}

	defer os.Remove(pdfFile.Name())

	if _, err := pdfFile.Write(pdfBytes); err != nil {
		pdfFile.Close()
		return nil, err
	}

	pdfFile.Close()

	outputPrefix := filepath.Join(os.TempDir(), "outthumb")
	cmd := exec.Command("pdftoppm", "-f", "1", "-l", "1", "-singlefile", "-png", pdfFile.Name(), outputPrefix)

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	pngPath := outputPrefix + ".png"
	defer os.Remove(pngPath)

	webpFile, err := os.CreateTemp("", "thumb-*.webp")

	if err != nil {
		return nil, err
	}

	webpPath := webpFile.Name()
	webpFile.Close()
	defer os.Remove(webpPath)

	cmdWebp := exec.Command("cwebp", pngPath, "-o", webpPath)

	if err := cmdWebp.Run(); err != nil {
		return nil, err
	}

	webpData, err := os.ReadFile(webpPath)

	if err != nil {
		return nil, err
	}

	return webpData, nil
}
