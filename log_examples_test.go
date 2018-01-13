package httplog_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/crhntr/httplog"
)

func ExampleWrap() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Hello, world!")
	})

	logMux := httplog.Wrap(mux, httplog.StandardOut{})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	logMux.ServeHTTP(w, r)
}
