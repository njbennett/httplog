package httplog

import (
	"log"
	"net/http"
	"time"
)

type Logger interface {
	Log(req *http.Request, elapsed time.Duration, status int)
}

type StandardOut struct{}

func (_ StandardOut) Log(r *http.Request, elapsed time.Duration, status int) {
	log.Printf("HTTP\t%-3d\t\t%s\t%s\t%s", status, elapsed, r.Method, r.URL.Path)
}

type logRecord struct {
	http.ResponseWriter
	status int
}

func (r *logRecord) Write(p []byte) (int, error) {
	return r.ResponseWriter.Write(p)
}

// WriteHeader implements ResponseWriter for logRecord
func (r *logRecord) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

// WrapHandler implements ResponseWriter for logRecord
// logger should allow concurent agitccess
func Wrap(f http.Handler, logger Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		record := &logRecord{
			ResponseWriter: w,
		}

		start := time.Now()
		f.ServeHTTP(record, r)

		logger.Log(r, time.Since(start), record.status)
	}
}
