package config

const VERSION = "1.1.0"

const THUMB_PATH string = "static/thumbnails/"
const LOCAL string = "127.0.0.1"
const PAGE_SIZE int = 8

const MAX_UPLOAD_SIZE = 60 << 20

type Mode int

const (
	FS Mode = iota + 1
	S3
)

const VIEWER_PATH = "/opt/gochive/lib/pdfjs/web/viewer.html"
