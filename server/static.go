package server

import (
	"fmt"
	"net/http"
	"os"
)

// StaticServer ist ein einfacher Server für statische Dateien
type StaticServer struct {
	rootDir string
}

// NewStaticServer erstellt einen neuen StaticServer
func NewStaticServer(rootDir string) *StaticServer {
	return &StaticServer{
		rootDir: rootDir,
	}
}

// Start startet den StaticServer auf dem angegebenen Port
func (s *StaticServer) Start(port int) error {
	// Stelle sicher, dass das Verzeichnis existiert
	if err := os.MkdirAll(s.rootDir, 0755); err != nil {
		return err
	}

	// Erstelle einen FileServer für das Verzeichnis
	fileServer := http.FileServer(http.Dir(s.rootDir))

	// Erstelle einen Handler, der alle Anfragen an den FileServer weiterleitet
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Wenn der Pfad "/" ist, zeige index.html an
		if r.URL.Path == "/" {
			r.URL.Path = "/index.html"
		}

		// Leite die Anfrage an den FileServer weiter
		fileServer.ServeHTTP(w, r)
	})

	// Starte den Server
	return http.ListenAndServe(fmt.Sprintf(":%d", port), handler)
}
