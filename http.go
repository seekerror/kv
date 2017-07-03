package kv

import (
	"net/http"
	"strings"
)

// TODO(herohde) 7/2/2017: do we need a RW (or WO) handler?

// NewHandler returns a http.Handler that serves data from the given reader.
// It treats the path as the key. For example:
//
//    HTTP GET "foo.com:1234/foo/bar" maps to r.Read("foo/bar")
//
func NewHandler(r Reader) http.Handler {
	return &handler{r}
}

type handler struct {
	r Reader
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := strings.Trim(r.URL.Path, "/")

	data, err := h.r.Read(key)
	if err != nil {
		code := http.StatusInternalServerError
		if err == KeyNotFoundErr {
			code = http.StatusNotFound
		}
		w.WriteHeader(code)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
