package httplog

import (
	"fmt"
	"net/http"
	"time"
)

func JSON(req *http.Request, elapsed time.Duration, status int) {
	fmt.Printf(`{"type": "HTTP_REQUEST", "method": %q, "path": %q, "duration": %q, "status": %q}`+"\n", req.Method, req.URL.Path, elapsed, status)
}

type Func func(req *http.Request, elapsed time.Duration, status int)

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

func Wrap(f http.Handler, logFn Func) http.HandlerFunc {
	if logFn == nil {
		logFn = JSON
	}
	return func(w http.ResponseWriter, r *http.Request) {
		record := &logRecord{
			ResponseWriter: w,
		}

		start := time.Now()
		f.ServeHTTP(record, r)

		logFn(r, time.Since(start), record.status)
	}
}
