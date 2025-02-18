package api

import (
	"bytes"

	"github.com/gin-gonic/gin"
)

type Writer struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func NewWriter(w gin.ResponseWriter) *Writer {
	return &Writer{body: bytes.NewBufferString(""), ResponseWriter: w}
}

func (w Writer) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w Writer) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func (w Writer) String() string {
	return w.body.String()
}
