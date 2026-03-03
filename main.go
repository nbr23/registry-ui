package main

import (
	"embed"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//go:embed static
var staticFiles embed.FS

var registryURL *url.URL

func main() {
	raw := os.Getenv("REGISTRY_URL")
	if raw == "" {
		raw = "http://localhost:5000"
	}
	var err error
	registryURL, err = url.Parse(raw)
	if err != nil {
		log.Fatalf("invalid REGISTRY_URL: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/", handleRegistryProxy)
	mux.HandleFunc("/", handleStatic)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func validateCredentials(username, password string) bool {
	req, err := http.NewRequest("GET", registryURL.String()+"/v2/", nil)
	if err != nil {
		return false
	}
	req.SetBasicAuth(username, password)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func handleRegistryProxy(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok || !validateCredentials(username, password) {
		w.Header().Set("WWW-Authenticate", `Basic realm="Docker Registry"`)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	target := registryURL.String() + strings.TrimPrefix(r.URL.Path, "/api")
	if r.URL.RawQuery != "" {
		target += "?" + r.URL.RawQuery
	}

	proxyReq, err := http.NewRequest(r.Method, target, r.Body)
	if err != nil {
		http.Error(w, "proxy error", http.StatusInternalServerError)
		return
	}

	proxyReq.SetBasicAuth(username, password)

	for _, h := range []string{"Content-Type", "Content-Length"} {
		if v := r.Header.Get(h); v != "" {
			proxyReq.Header.Set(h, v)
		}
	}

	if strings.Contains(r.URL.Path, "/manifests/") {
		proxyReq.Header.Set("Accept", strings.Join([]string{
			"application/vnd.docker.distribution.manifest.v2+json",
			"application/vnd.docker.distribution.manifest.list.v2+json",
			"application/vnd.oci.image.manifest.v1+json",
			"application/vnd.oci.image.index.v1+json",
			"application/vnd.docker.distribution.manifest.v1+prettyjws",
		}, ", "))
	}

	client := &http.Client{
		Timeout: 5 * time.Minute,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, "registry unreachable", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for _, h := range []string{
		"Content-Type", "Content-Length", "Docker-Content-Digest",
		"Link", "Docker-Distribution-Api-Version",
		"Location", "Www-Authenticate",
	} {
		if v := resp.Header.Get(h); v != "" {
			w.Header().Set(h, v)
		}
	}

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func handleStatic(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok || !validateCredentials(username, password) {
		w.Header().Set("WWW-Authenticate", `Basic realm="Docker Registry"`)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	path := "static" + r.URL.Path

	data, err := staticFiles.ReadFile(path)
	if err != nil {
		data, _ = staticFiles.ReadFile("static/index.html")
		path = "static/index.html"
	}

	ct := mime.TypeByExtension(filepath.Ext(path))
	if ct == "" {
		ct = http.DetectContentType(data)
	}
	w.Header().Set("Content-Type", ct)
	w.Write(data)
}
