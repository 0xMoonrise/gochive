package utils

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/0xMoonrise/gochive/internal/config"
	"github.com/chai2010/webp"
	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/single_threaded"
)

var pool pdfium.Pool
var instance pdfium.Pdfium

func renderPage(file []byte, page int, output string) (error, []byte) {
	var rawImage bytes.Buffer
	pool = single_threaded.Init(single_threaded.Config{})
	instance, err := pool.GetInstance(time.Second * 30)

	if err != nil {
		log.Fatal(err)
	}

	doc, err := instance.OpenDocument(&requests.OpenDocument{
		File: &file,
	})

	if err != nil {
		return err, nil
	}

	defer instance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
		Document: doc.Document,
	})

	pageRender, err := instance.RenderPageInDPI(&requests.RenderPageInDPI{
		DPI: 200,
		Page: requests.Page{
			ByIndex: &requests.PageByIndex{
				Document: doc.Document,
				Index:    page,
			},
		},
	})

	if err != nil {
		return err, nil
	}

	err = webp.Encode(&rawImage, pageRender.Result.Image, &webp.Options{Quality: 100})
	if err != nil {
		return err, nil
	}

	f, err := os.Create(filepath.Join(config.THUMB_PATH, output))
	if err != nil {
		return err, nil
	}
	defer f.Close()

	_, err = f.Write(rawImage.Bytes())
	if err != nil {
		return err, nil
	}

	return nil, rawImage.Bytes()
}

func MakeThumbnail(rawFile []byte, filename string) (error, []byte) {

	err, rawImage := renderPage(rawFile, 0, filename)
	if err != nil {
		return err, nil
	}

	return nil, rawImage
}

func ValidateFilename(filename string) bool {

	match, _ := regexp.MatchString("^.+(pdf|md)$", filename)
	return match

}
