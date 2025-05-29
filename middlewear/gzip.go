package middlewear

import (
	"compress/gzip"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type GzipWriter struct {
	gzipWriter *gzip.Writer
	httpWriter http.ResponseWriter
}

func (g *GzipWriter) Header() http.Header {
	return g.httpWriter.Header()
}

func (g *GzipWriter) WriteHeader(statusCode int) {
	g.httpWriter.WriteHeader(statusCode)
}

func (g *GzipWriter) Write(b []byte) (int, error) {
	if g.gzipWriter != nil {
		return g.gzipWriter.Write(b)
	}
	return g.httpWriter.Write(b)
}

func (g *GzipWriter) Close() {
	if g.gzipWriter != nil {
		g.gzipWriter.Close()
	}
}

func NewGzipWriter(c *gin.Context) *GzipWriter {
	if !strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
		return &GzipWriter{
			httpWriter: c.Writer,
		}
	}

	c.Writer.Header().Set("Content-Encoding", "gzip")
	c.Writer.Header().Set("Vary", "Accept-Encoding")

	gz, _ := gzip.NewWriterLevel(c.Writer, gzip.DefaultCompression)

	return &GzipWriter{
		gzipWriter: gz,
		httpWriter: c.Writer,
	}
}
