package proxy

import (
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

var client = &http.Client{Timeout: 10 * time.Second}

func ReverseProxyHandler(w http.ResponseWriter, r *http.Request) {
	backendURL := getBackendTarget(r)

	cacheKey := r.Method + ":" + r.URL.String()
	if body, ok := GetFromCache(cacheKey); ok {
		w.Header().Set("X-Cache", "HIT")
		w.Write(body)
		return
	}

	proxyReq, err := http.NewRequest(r.Method, backendURL+r.RequestURI, r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		Log("Failed to create request: " + err.Error())
		return
	}

	proxyReq.Header = r.Header.Clone()

	proxyReq.Header.Set("X-Forwarded-For", r.RemoteAddr)

	resp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		Log("Request failed: " + err.Error())
		return
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		w.Header()[k] = v
	}

	w.WriteHeader(resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	w.Write(body)

	SaveToCache(cacheKey, body)

}

func getBackendTarget(r *http.Request) string {
	for _, route := range ProxyConfig.Routes {
		log.Printf("Request Path: %s", r.Host)
		log.Printf("Server: %s", route.Server)
		log.Printf("Path Prefix: %s", route.PathPrefix)
		if strings.EqualFold(r.Host, route.Server) && strings.HasPrefix(r.URL.Path, route.PathPrefix) {
			log.Printf("Matched route: %s", route.PathPrefix)
			log.Printf("Backend: %s", route.Backend)
			log.Printf("Server: %s", route.Server)
			log.Printf("Request URI: %s", r.RequestURI)
			log.Printf("Request Path: %s", r.URL.Path)
			return route.Backend
		}
	}
	return ""
}
