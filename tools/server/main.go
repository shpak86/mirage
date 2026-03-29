package main

import (
	"net/http"
	"net/http/httputil"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello!"))
	})

	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{\"message\":\"hello\"}"))
	})

	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		req, _ := httputil.DumpRequest(r, true)
		w.Write([]byte(req))
	})

	http.HandleFunc("/response400", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"message\":\"400\"}"))
	})
	http.ListenAndServe(":8080", nil)
}
