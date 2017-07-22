package static

var (
	mimeTypes = map[string]string{
		".js":    "application/javascript",
		".json":  "application/json",
		".css":   "text/css",
		".html":  "text/html",
		".jpg":   "image/jpeg",
		".png":   "image/png",
		".svg":   "image/svg+xml",
		".ico":   "image/vnd.microsoft.icon",
		".otf":   "application/font-sfnt",
		".ttf":   "application/font-sfnt",
		".woff":  "application/font-woff",
		".woff2": "font/woff2",
		".ijmap": "text/plain",
		".eot":   "application/vnd.ms-fontobject",
		".map":   "application/json",
	}
)

type File struct {
	Type     string
	Hash     string
	Data     []byte
	GzipData []byte
}
