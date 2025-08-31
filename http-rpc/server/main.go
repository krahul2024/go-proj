package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"srv/route"
	"srv/util"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	util.Info("Environment variables loaded successfully!")

	router := http.NewServeMux()

	fileRouter := route.File()
	router.Handle("/file/", http.StripPrefix("/file", fileRouter))

	server := http.Server {
		Addr : fmt.Sprintf(":%v", 4000),
		Handler: ReqLogger(router),
		// add things to this server, otherwise provided by the default server
	}

	util.Info ((fmt.Sprintf("Server running on : %v", 4000)))
	if err := server.ListenAndServe(); err != nil {
		util.Fatal(err, "Error starting the server!")
		return
	}
}

func ReqLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		var requestBody string
		if r.Body != nil {
			buf, _ := io.ReadAll(r.Body)
			requestBody = string(buf)
			r.Body = io.NopCloser(bytes.NewBuffer(buf))
		}

		ip := r.RemoteAddr
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			ip = forwarded
		}

		referrer := r.Referer()
		if referrer == "" {
			referrer = "N/A"
		}

		headers := map[string]string{}
		for name, values := range r.Header {
			headers[name] = strings.Join(values, ", ")
		}

		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		util.Info (
			"Request Information",
			map[string]any{
				"protocol":  r.Proto,
				"status":    rw.statusCode,
				"method":    r.Method,
				"url":       r.URL.Path,
				"ip":        ip,
				"duration":  duration.String(),
			},
			map[string]any{
				"referrer": referrer,
				"headers":  headers,
				"body":     requestBody,
			},
			)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
