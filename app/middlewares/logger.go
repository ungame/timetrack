package middlewares

import (
	"log"
	"net/http"
	"time"
)

type statusCodeRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusCodeRecorder) WriteHeader(status int) {
	r.ResponseWriter.WriteHeader(status)
	r.status = status
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var (
			start    = time.Now()
			proto    = request.Proto
			method   = request.Method
			uri      = request.RequestURI
			recorder = &statusCodeRecorder{writer, 0}
		)

		next.ServeHTTP(recorder, request)

		log.Printf("%s %s %s %d %s %s\n",
			proto, method, uri, recorder.status,
			http.StatusText(recorder.status),
			time.Since(start).String())
	})
}
