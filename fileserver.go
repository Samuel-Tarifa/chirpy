package main

import "net/http"

func makeFileServer(dir string) http.Handler {
    return http.StripPrefix("/app", http.FileServer(http.Dir(dir)))
}