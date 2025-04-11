package server

import (
	"fmt"
	"net/http"
	"os"
)

type StaticServer struct {
	rootDir string
}

func NewStaticServer(rootDir string) *StaticServer {
	return &StaticServer{
		rootDir: rootDir,
	}
}

func (s *StaticServer) Start(port int) error {
	if err := os.MkdirAll(s.rootDir, 0755); err != nil {
		return err
	}

	fileServer := http.FileServer(http.Dir(s.rootDir))

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			r.URL.Path = "/index.html"
		}
		fileServer.ServeHTTP(w, r)
	})

	return http.ListenAndServe(fmt.Sprintf(":%d", port), handler)
}
