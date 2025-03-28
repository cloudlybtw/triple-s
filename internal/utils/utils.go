package utils

import (
	"net/http"
	"regexp"
	"strings"
)

func ResponceXML(w http.ResponseWriter, r *http.Request, code int, data []byte) {
	Responce(w, r, "application/xml", code, data)
}

func Responce(w http.ResponseWriter, r *http.Request, content_type string, code int, data []byte) {
	w.Header().Add("Content-Type", content_type)
	w.WriteHeader(code)
	if data != nil {
		w.Write(data)
	}
}

func CheckName(name string) bool {
	if len(name) < 3 || len(name) > 63 {
		return false
	}

	if !regexp.MustCompile(`^[a-z0-9\.-]+$`).MatchString(name) {
		return false
	}

	if strings.HasPrefix(name, "-") || strings.HasPrefix(name, ".") || strings.HasSuffix(name, "-") || strings.HasSuffix(name, ".") {
		return false
	}

	if strings.Contains(name, "--") || strings.Contains(name, "..") {
		return false
	}

	if regexp.MustCompile(`^\d+\.\d+\.\d+\.\d+$`).MatchString(name) {
		return false
	}

	if name == "buckets.csv" || name == "objects.csv" {
		return false
	}

	return true
}
