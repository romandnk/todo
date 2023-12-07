package utils

import (
	"net/http"
	"strconv"
	"time"
)

type RequestInfo struct {
	ClientIP    string
	Date        string
	Method      string
	Path        string
	HTTPVersion string
	Latency     string
	UserAgent   string
}

func RequestInformation(r *http.Request, duration time.Duration) RequestInfo {
	clientIP := r.RemoteAddr
	date := time.Now().Format("02/Jan/2006:15:04:05 -0700")
	method := r.Method
	pth := r.URL.Path
	HTTPVersion := r.Proto
	latency := strconv.Itoa(int(duration.Milliseconds())) + "ms"
	userAgent := r.Header.Get("User-Agent")
	if userAgent == "" {
		userAgent = "empty"
	}
	return RequestInfo{
		ClientIP:    clientIP,
		Date:        date,
		Method:      method,
		Path:        pth,
		HTTPVersion: HTTPVersion,
		Latency:     latency,
		UserAgent:   userAgent,
	}
}
