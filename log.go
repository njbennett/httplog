package httplog

import (
	"fmt"
	"net/http"
	"time"
)

func JSON(req *http.Request, elapsed time.Duration, status int) {
	fmt.Printf(`{"type": "HTTP_REQUEST", "method": %q, "path": %q, "duration": %q, "status": %d}`+"\n", req.Method, req.URL.Path, elapsed, status)
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

func Wrap(f http.Handler, logFns ...Func) http.HandlerFunc {
	var fn Func
	if len(logFns) == 0 {
		fn = JSON
	} else if len(logFns) == 1 {
		fn = logFns[0]
	} else {
		fn = func(req *http.Request, elapsed time.Duration, status int) {
			for _, lg := range logFns {
				lg(req, elapsed, status)
			}
		}
	}
	return func(w http.ResponseWriter, r *http.Request) {
		record := &logRecord{
			ResponseWriter: w,
		}

		start := time.Now()
		f.ServeHTTP(record, r)

		fn(r, time.Since(start), record.status)
	}
}
