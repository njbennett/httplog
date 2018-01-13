# httplog
A simple http request logger

[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://github.com/crhntr/httplog)


```go

mux := http.NewServeMux()
mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
  fmt.Fprint(w, "Hello, world!")
})

logMux := httplog.Wrap(mux, httplog.StandardOut{})

w := httptest.NewRecorder()
r := httptest.NewRequest(http.MethodGet, "/", nil)
logMux.ServeHTTP(w, r)

```
