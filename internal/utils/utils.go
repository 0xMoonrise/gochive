package utils

import (
	"bytes"
	"io"
	"regexp"
	"time"

	"github.com/chai2010/webp"
	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/single_threaded"
)

var pool pdfium.Pool

// var instance pdfium.Pdfium
func MakeThumbnail(reader io.ReadSeeker, size int64, page int, imageBuffer *bytes.Buffer) (err error) {

	pool = single_threaded.Init(single_threaded.Config{})
	instance, err := pool.GetInstance(time.Second * 30)

	if err != nil {
		return
	}

	doc, err := instance.OpenDocument(&requests.OpenDocument{
		FileReader:     reader,
		FileReaderSize: size,
	})

	if err != nil {
		return
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
		return
	}
	err = webp.Encode(imageBuffer, pageRender.Result.Image, &webp.Options{Quality: 100})
	if err != nil {
		return
	}
	return
}

func ValidateFilename(filename string) bool {

	match, _ := regexp.MatchString("^.+(pdf|md)$", filename)
	return match

}
