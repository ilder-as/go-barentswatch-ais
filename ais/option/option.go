package option

import (
	"net/http"
	"time"
)

type Option func(r *http.Request)

func Since(t time.Time) Option {
	return func(r *http.Request) {
		params := r.URL.Query()
		params.Add("since", t.Format(time.RFC3339))
		r.URL.RawQuery = params.Encode()
	}
}
