package route

import "net/http"

func File() *http.ServeMux {
	mux := http.NewServeMux()

	return mux
}
